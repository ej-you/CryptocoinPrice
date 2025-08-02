package usecase

import (
	"CryptocoinPrice/internal/app/entity"
	"CryptocoinPrice/internal/app/repo"
	"errors"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

var _ CoinManageUsecase = (*CoinManageUC)(nil)

type CoinManageUC struct {
	coinRepoDB   repo.CoinRepoDB
	priceRepoDB  repo.PriceRepoDB
	priceRepoAPI repo.PriceRepoAPI
}

// NewCoinManageUC returns new coin manage usecase.
func NewCoinManageUC(coinRepoDB repo.CoinRepoDB,
	priceRepoDB repo.PriceRepoDB, priceRepoAPI repo.PriceRepoAPI) *CoinManageUC {

	return &CoinManageUC{
		coinRepoDB:   coinRepoDB,
		priceRepoDB:  priceRepoDB,
		priceRepoAPI: priceRepoAPI,
	}
}

// ObserveCoin creates new observed coin or sets observed on true for existing coin.
func (u *CoinManageUC) ObserveCoin(symbol string) (*entity.Coin, error) {
	// get coin from DB by symbol
	coin, err := u.coinRepoDB.GetBySymbol(symbol)
	// true if coin is not found
	coinNotFound := errors.Is(err, repo.ErrNotFound)
	// if unknown error
	if err != nil && !coinNotFound {
		return nil, fmt.Errorf("get by symbol: %w", err)
	}

	// if coin is found - set observed on true for coin
	if !coinNotFound {
		coin.Observed = true
		coinUpdates := &entity.CoinPartial{Observed: &coin.Observed}
		if err := u.coinRepoDB.Update(coin.ID, coinUpdates); err != nil {
			return nil, fmt.Errorf("update: %w", err)
		}
		return coin, nil
	}

	// if coin is not found
	// get coin price from API to check that coin exists in the world.
	coinPrice, err := u.priceRepoAPI.OneCoinPrice(symbol)
	// if coin symbol is invalid
	if errors.Is(err, repo.ErrValidateData) {
		return nil, fmt.Errorf("%w: invalid symbol: unexisting coin", ErrValidateData)
	}
	if err != nil {
		return nil, fmt.Errorf("check coin: %w", err)
	}

	// create coin
	coin, err = u.coinRepoDB.Create(symbol)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	// save coin price into DB
	_, err = u.priceRepoDB.Create(coin, coinPrice.Price, time.Now().UTC().Unix())
	if err != nil {
		logrus.Errorf("Save coin price into DB: %v", err)
	}
	return coin, nil
}

// DisableObserveCoin sets observed on false for coin.
func (u *CoinManageUC) DisableObserveCoin(symbol string) (*entity.Coin, error) {
	// get coin from DB by symbol
	coin, err := u.coinRepoDB.GetBySymbol(symbol)
	// if coin is not found
	if errors.Is(err, repo.ErrNotFound) {
		return nil, fmt.Errorf("get coin: %w", ErrNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("get by symbol: %w", err)
	}

	coin.Observed = false
	coinUpdates := &entity.CoinPartial{Observed: &coin.Observed}
	if err := u.coinRepoDB.Update(coin.ID, coinUpdates); err != nil {
		return nil, fmt.Errorf("update: %w", err)
	}
	return coin, nil
}

// GetNearestPrice returns first price with coin symbol
// and nearest timestamp for given timestamp.
func (u *CoinManageUC) GetNearestPrice(symbol string, timestamp int64) (*entity.Price, error) {
	// get coin from DB by symbol
	coin, err := u.coinRepoDB.GetBySymbol(symbol)
	// if coin is not found
	if errors.Is(err, repo.ErrNotFound) {
		return nil, fmt.Errorf("get coin: %w", ErrNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("get coin by symbol: %w", err)
	}

	price, err := u.priceRepoDB.GetNearestTimestamp(coin, timestamp)
	// if price is not found
	if errors.Is(err, repo.ErrNotFound) {
		return nil, fmt.Errorf("get nearest price: %w", ErrNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("get nearest price: %w", err)
	}
	return price, nil
}
