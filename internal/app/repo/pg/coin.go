// Package pg contains PostgreSQL DB repos implementations for entities.
package pg

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"CryptocoinPrice/internal/app/entity"
	"CryptocoinPrice/internal/app/repo"
	// "CryptocoinPrice/internal/app/errors"
)

var _ repo.CoinRepo = (*CoinRepoPG)(nil)

var ErrNotFound = gorm.ErrRecordNotFound // record not found error

type CoinRepoPG struct {
	dbStorage *gorm.DB
}

// NewCoinRepoPG returns new PostgreSQL repo DB instance for coin entity.
func NewCoinRepoPG(dbStorage *gorm.DB) *CoinRepoPG {
	return &CoinRepoPG{
		dbStorage: dbStorage,
	}
}

// Create creates new coin.
func (r *CoinRepoPG) Create(symbol string) (*entity.Coin, error) {
	// init coin for creating
	coin := &entity.Coin{
		ID:       uuid.NewString(),
		Symbol:   symbol,
		Observed: true,
	}
	// create record
	if err := r.dbStorage.Create(coin).Error; err != nil {
		return nil, fmt.Errorf("create coin: %w", err)
	}
	return coin, nil
}

// GetBySymbol returns coin by given symbol.
// If symbol is not found it returns not found error
func (r *CoinRepoPG) GetBySymbol(symbol string) (*entity.Coin, error) {
	coin := &entity.Coin{Symbol: symbol}

	err := r.dbStorage.Where(coin).First(coin).Error
	// if record not found
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return coin, nil
}

// Update updates coin.
// It selects coin by given ID and replace
// all old values (from DB) to new (given).
func (r *CoinRepoPG) Update(coinID string, coinUpdates *entity.CoinPartial) error {
	// update coin
	err := r.dbStorage.Model(&entity.Coin{}).
		Where("id = ?", coinID).
		Updates(coinUpdates).Error
	// error
	if err != nil {
		return fmt.Errorf("update: %w", err)
	}
	return nil
}

// GetList gets all subscriptions and returns it.
func (r *CoinRepoPG) GetObserved() (entity.CoinList, error) {
	coinList := entity.CoinList{}
	err := r.dbStorage.Where("observed = ?", true).Find(&coinList).Error
	if err != nil {
		return nil, err
	}
	return coinList, nil
}
