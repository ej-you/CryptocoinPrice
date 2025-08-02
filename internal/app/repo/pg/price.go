package pg

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"CryptocoinPrice/internal/app/entity"
	"CryptocoinPrice/internal/app/repo"
)

var _ repo.PriceRepoDB = (*PriceRepoPG)(nil)

type PriceRepoPG struct {
	dbStorage *gorm.DB
}

// NewPriceRepoPG returns new PostgreSQL repo DB instance for price entity.
func NewPriceRepoPG(dbStorage *gorm.DB) *PriceRepoPG {
	return &PriceRepoPG{
		dbStorage: dbStorage,
	}
}

// Create creates new price.
// Coin ID must be presented in the given coin instance.
// Also this coin instance passes into the price instance.
func (r *PriceRepoPG) Create(coin *entity.Coin,
	price float64, timestamp int64) (*entity.Price, error) {

	priceObj := &entity.Price{
		ID:        uuid.NewString(),
		CoinID:    coin.ID,
		Price:     fmt.Sprint(price),
		Timestamp: timestamp,
		Coin:      coin,
	}
	if err := r.dbStorage.Create(priceObj).Error; err != nil {
		return nil, err
	}
	return priceObj, nil
}

// GetNearestTimestamp returns price for given coin at the
// given timestamp or the nearest timestamp from the given timestamp.
// Coin ID must be presented in the given coin instance.
// Also this coin instance passes into the price instance.
func (r *PriceRepoPG) GetNearestTimestamp(coin *entity.Coin,
	timestamp int64) (*entity.Price, error) {

	price := &entity.Price{Coin: coin}
	err := r.dbStorage.Raw(`
		SELECT * FROM prices WHERE coin_id = ? ORDER BY ABS(timestamp - ?) LIMIT 1`,
		coin.ID, timestamp).
		Scan(price).Error

	// if record is not found
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, repo.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return price, nil
}
