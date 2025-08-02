// Package coingecko contains Coingecko API repos implementations for entities.
package coingecko

import (
	"fmt"
	"strings"
	"time"

	resty "github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"

	"CryptocoinPrice/internal/app/entity"
	"CryptocoinPrice/internal/app/repo"
)

var _ repo.PriceRepoAPI = (*PriceRepoCoingecko)(nil)

const (
	_requestTimeout = 2 * time.Second        // timeout for requests to API
	_retryCount     = 3                      // amount of retries attempts in error cases
	_retryInitTime  = 500 * time.Millisecond // time between first request and first retry
	_retryMaxTime   = 2 * time.Second        // max time between request and retry

	_coinDataPriceKey      = "usd"             // key for coin price in coin data map
	_coinDataLastUpdateKey = "last_updated_at" // key for coin last update time in coin data map
)

// rawCoinsData is a raw response from API. It's a map (with keys - coin names)
// of maps with string-any key-values (coin data).
type rawCoinsData map[string]map[string]any

type PriceRepoCoingecko struct {
	apiKey string
	client *resty.Client
}

// NewPriceRepoCoingecko returns new Coingecko API repo instance for price entity.
func NewPriceRepoCoingecko(apiKey string) *PriceRepoCoingecko {
	// init HTTP-client with retry params
	restyClient := resty.New().
		SetTimeout(_requestTimeout).
		SetRetryCount(_retryCount).
		SetRetryWaitTime(_retryInitTime).
		SetRetryMaxWaitTime(_retryMaxTime)

	return &PriceRepoCoingecko{
		apiKey: apiKey,
		client: restyClient,
	}
}

// OneCoinPrice sends request to API for
// one coin (with given symbol) price and returns it.
// The price is in USD.
// API response looks like:
//
//	{
//	  "btc": {
//	    "usd": 115380,
//	    "last_updated_at": 1754050754
//	  }
//	}
func (r *PriceRepoCoingecko) OneCoinPrice(symbol string) (*entity.CoinPriceAPI, error) {
	// map of coins each of which are map with coin data
	var rawData rawCoinsData

	// do request to REST API and parse JSON-response into result
	_, err := r.client.R().
		SetResult(&rawData).
		SetHeader("Accept", "application/json").
		SetHeader("x-cg-demo-api-key", r.apiKey).
		SetQueryParam("vs_currencies", "usd").
		SetQueryParam("include_last_updated_at", "true").
		SetQueryParam("symbols", symbol).
		Get("https://api.coingecko.com/api/v3/simple/price")
	if err != nil {
		return nil, fmt.Errorf("request to api: %w", err)
	}

	// parse coin data into struct
	coinData, err := parseCoinData(rawData, symbol)
	if err != nil {
		return nil, fmt.Errorf("parse coin data: %w", err)
	}
	return coinData, nil
}

// ManyCoinPrices sends request to API for
// many coins' (with given symbols) prices and returns them.
// The prices is in USD.
// API response looks like:
//
//	{
//	  "btc": {
//	    "usd": 115380,
//	    "last_updated_at": 1754050754
//	  },
//	  "eth": {
//	    "usd": 3647.54,
//	    "last_updated_at": 1754050755
//	  }
//	}
func (r *PriceRepoCoingecko) ManyCoinPrices(symbols []string) (entity.CoinPriceAPIList, error) {
	var (
		// map of coins each of which are map with coin data
		rawData rawCoinsData
		// coin parsing errors
		errList = make([]error, 0)
	)

	// do request to REST API and parse JSON-response into result
	_, err := r.client.R().
		SetResult(&rawData).
		SetHeader("Accept", "application/json").
		SetHeader("x-cg-demo-api-key", r.apiKey).
		SetQueryParam("vs_currencies", "usd").
		SetQueryParam("include_last_updated_at", "true").
		SetQueryParam("symbols", strings.Join(symbols, ",")).
		Get("https://api.coingecko.com/api/v3/simple/price")
	if err != nil {
		return nil, fmt.Errorf("request to api: %w", err)
	}

	// init coin prices slice
	symbolsAmount := len(symbols)
	coinPricesList := make(entity.CoinPriceAPIList, 0, symbolsAmount)

	var coinData *entity.CoinPriceAPI
	// parse each coin
	for _, symbol := range symbols {
		coinData, err = parseCoinData(rawData, symbol)
		if err != nil {
			errList = append(errList, err)
		}
		coinPricesList = append(coinPricesList, *coinData)
	}

	errsAmount := len(errList)
	logrus.Infof("Get coin prices: %d/%d", symbolsAmount-errsAmount, symbolsAmount)
	// if no one parsing error is occurred
	if errsAmount == 0 {
		return coinPricesList, nil
	}

	// collect err slice into one error
	var errBuilder strings.Builder
	errBuilder.WriteString(errList[0].Error())
	for _, parseErr := range errList[1:] {
		errBuilder.WriteString(" && ")
		errBuilder.WriteString(parseErr.Error())
	}
	return coinPricesList, fmt.Errorf("parse coins data: %s", errBuilder.String())
}

// parseCoinData parses specific coin data from raw map with coins from API.
// API response looks like:
//
//	{
//	  "btc": {
//	    "usd": 115380,
//	    "last_updated_at": 1754050754
//	  },
//	  "eth": {
//	    "usd": 3647.54,
//	    "last_updated_at": 1754050755
//	  }
//	}
//
// So, given rawCoinsData must be a map (with keys - coin names)
// of maps with string-any key-values (coin data).
// Given symbol value is the name of needed coin to parse.
func parseCoinData(rawCoinsData map[string]map[string]any,
	symbol string) (*entity.CoinPriceAPI, error) {

	coinData, found := rawCoinsData[symbol]
	// if coin data is not found in result
	if !found {
		return nil, fmt.Errorf("%w: coin %s is not found", repo.ErrValidateData, symbol)
	}

	var ok bool // nolint:varnamelen // generally accepted variable name
	coinPriceObj := &entity.CoinPriceAPI{Symbol: symbol}
	// parse coin price
	coinPriceObj.Price, ok = coinData[_coinDataPriceKey].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid coin price: coin data - %v", coinData)
	}
	// parse coin last update time (by default, numbers deserialized into float64)
	floatCoinLastUpdate, ok := coinData[_coinDataLastUpdateKey].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid coin last update time: coin data - %v", coinData)
	}
	// cast to integer
	coinPriceObj.LastUpdate = int64(floatCoinLastUpdate)

	return coinPriceObj, nil
}
