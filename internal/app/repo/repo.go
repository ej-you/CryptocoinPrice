// Package repo contains repo interfaces for entities and
// its implementations in subdirs.
package repo

import "CryptocoinPrice/internal/app/entity"

type CoinRepo interface {
	Create(symbol string) (*entity.Coin, error)
	GetBySymbol(symbol string) (*entity.Coin, error)
	Update(coinID string, coinUpdates *entity.CoinPartial) error
	GetObserved() (entity.CoinList, error)
}

type PriceRepo interface {
	Create(coin *entity.Coin, price float64, timestamp int) (*entity.Price, error)
	GetNearestTimestamp(coin *entity.Coin, timestamp int) (*entity.Price, error)
}
