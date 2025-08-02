package coingecko

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"CryptocoinPrice/config"
)

var (
	_testPriceRepo *PriceRepoCoingecko

	_testCoinSymbols = []string{"btc", "eth", "ton"}
)

func TestMain(m *testing.M) {
	// load config
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("parse config: %v", err)
	}
	// create repo
	_testPriceRepo = NewPriceRepoCoingecko(cfg.App.CoingeckoAPIKey)

	// run tests
	os.Exit(m.Run())
}

func TestCoinRepoCoingecko_OneCoinPrice(t *testing.T) {
	t.Log("Get coin price from API")

	coinPrice, err := _testPriceRepo.OneCoinPrice(_testCoinSymbols[0])
	require.NoError(t, err)

	t.Logf("Coin price: %+v", coinPrice)
}

func TestCoinRepoCoingecko_ManyCoinPrices(t *testing.T) {
	t.Log("Get coins' prices from API")

	coinPricesList, err := _testPriceRepo.ManyCoinPrices(_testCoinSymbols)
	require.NoError(t, err)

	t.Logf("Coins' prices: %+v", coinPricesList)
}
