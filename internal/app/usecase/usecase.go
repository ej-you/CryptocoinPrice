// Package usecase contains usecases for app and its implementations.
package usecase

import (
	"errors"

	"CryptocoinPrice/internal/app/entity"
)

var (
	ErrNotFound     = errors.New("not found")     // not found error
	ErrValidateData = errors.New("validate data") // validat data error
)

// CoinManageUsecase used to manage observed coins and its prices.
type CoinManageUsecase interface {
	// ObserveCoin creates new observed coin or sets observed on true for existing coin.
	// Before creating new coin it gets price
	// for coin to check that coin exists in the world.
	ObserveCoin(symbol string) (*entity.Coin, error)
	// DisableObserveCoin sets observed on false for coin.
	DisableObserveCoin(symbol string) (*entity.Coin, error)
	// GetNearestPrice returns first price with coin symbol
	// and nearest timestamp for given timestamp.
	GetNearestPrice(symbol string, timestamp int64) (*entity.Price, error)
}

// PriceCollectorUsecase used to get new coin prices.
type PriceCollectorUsecase interface {
	// GetNewObservedCoinPrices gets new prices for observed coins.
	GetNewObservedCoinPrices() (entity.PriceList, error)
	// SaveCoinPrices saves coin prices.
	SaveCoinPrices(priceList entity.PriceList) (entity.PriceList, error)
}
