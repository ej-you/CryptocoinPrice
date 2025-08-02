package usecase

import (
	"CryptocoinPrice/internal/app/entity"
	"CryptocoinPrice/internal/app/repo"
	"fmt"
	"slices"
	"time"
)

var _ PriceCollectorUsecase = (*PriceCollectorUC)(nil)

type PriceCollectorUC struct {
	coinRepoDB   repo.CoinRepoDB
	priceRepoDB  repo.PriceRepoDB
	priceRepoAPI repo.PriceRepoAPI
}

// NewPriceCollectorUC returns new price collector usecase.
func NewPriceCollectorUC(coinRepoDB repo.CoinRepoDB,
	priceRepoDB repo.PriceRepoDB, priceRepoAPI repo.PriceRepoAPI) *PriceCollectorUC {

	return &PriceCollectorUC{
		coinRepoDB:   coinRepoDB,
		priceRepoDB:  priceRepoDB,
		priceRepoAPI: priceRepoAPI,
	}
}

// GetNewObservedCoinPrices gets new prices for observed coins.
func (u *PriceCollectorUC) GetNewObservedCoinPrices() (entity.PriceList, error) {
	// get observed coins
	observedCoins, err := u.coinRepoDB.GetObserved()
	if err != nil {
		return nil, fmt.Errorf("get observed coins: %w", err)
	}
	// takes only coins' symbols
	coinSymbols := make([]string, 0, len(observedCoins))
	for _, coin := range observedCoins {
		coinSymbols = append(coinSymbols, coin.Symbol)
	}

	// get coin prices
	coinPrices, err := u.priceRepoAPI.ManyCoinPrices(coinSymbols)

	updateTime := time.Now().UTC().Unix()
	var coin entity.Coin
	// fill price list
	priceList := make(entity.PriceList, len(coinPrices))
	for i, coinPrice := range coinPrices {
		coin = observedCoins[slices.IndexFunc(observedCoins, func(coin entity.Coin) bool {
			return coin.Symbol == coinPrice.Symbol
		})]
		priceList[i] = entity.Price{
			Coin:      &coin,
			Price:     fmt.Sprint(coinPrice.Price),
			Timestamp: updateTime,
			CoinID:    coin.ID,
		}
	}
	return priceList, err
}

// SaveCoinPrices saves coin prices.
func (u *PriceCollectorUC) SaveCoinPrices(priceList entity.PriceList) (entity.PriceList, error) {
	priceList, err := u.priceRepoDB.CreateMany(priceList)
	if err != nil {
		return nil, fmt.Errorf("create many: %w", err)
	}
	return priceList, nil
}
