// Package repo contains repo interfaces for entities and
// its implementations in subdirs.
package repo

import (
	"errors"

	"CryptocoinPrice/internal/app/entity"
)

var (
	ErrNotFound     = errors.New("record not found") // record not found error
	ErrValidateData = errors.New("validate data")    // validat data error
)

type CoinRepoDB interface {
	Create(symbol string) (*entity.Coin, error)
	GetBySymbol(symbol string) (*entity.Coin, error)
	Update(coinID string, coinUpdates *entity.CoinPartial) error
	GetObserved() (entity.CoinList, error)
}

type PriceRepoDB interface {
	Create(coin *entity.Coin, price float64, timestamp int64) (*entity.Price, error)
	CreateMany(priceList entity.PriceList) (entity.PriceList, error)
	GetNearestTimestamp(coin *entity.Coin, timestamp int64) (*entity.Price, error)
}

type PriceRepoAPI interface {
	OneCoinPrice(symbol string) (*entity.CoinPriceAPI, error)
	ManyCoinPrices(symbols []string) (entity.CoinPriceAPIList, error)
}
