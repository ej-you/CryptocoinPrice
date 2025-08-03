package pricecollector

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"CryptocoinPrice/config"
	repocoingecko "CryptocoinPrice/internal/app/repo/coingecko"
	repopg "CryptocoinPrice/internal/app/repo/pg"
	"CryptocoinPrice/internal/app/usecase"
)

// Price collector.
type PriceCollector struct {
	priceCollectorUC usecase.PriceCollectorUsecase
	tickerInterval   time.Duration
}

// New returns new price collector instance.
func New(cfg *config.Config, db *gorm.DB) *PriceCollector {
	// create repos
	coinRepoPG := repopg.NewCoinRepoPG(db)
	priceRepoDB := repopg.NewPriceRepoPG(db)
	priceRepoCoingecko := repocoingecko.NewPriceRepoCoingecko(cfg.App.CoingeckoAPIKey)
	// create usecases
	priceCollectorUC := usecase.NewPriceCollectorUC(coinRepoPG, priceRepoDB, priceRepoCoingecko)

	return &PriceCollector{
		priceCollectorUC: priceCollectorUC,
		tickerInterval:   cfg.App.PriceCollectInterval,
	}
}

// StartWithShutdown starts price collector and waits for
// context is done for gracefully shutdown collector.
// This method is blocking.
func (p *PriceCollector) StartWithShutdown(ctx context.Context) error {
	logrus.Info("Start price collector")
	defer logrus.Info("Price collector is shutdown")

	ticker := time.NewTicker(p.tickerInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			p.collect()
		case <-ctx.Done():
			return nil
		}
	}
}

// collect collects new observed coin prices and saves them.
func (p *PriceCollector) collect() {
	// get new prices
	newPrices, err := p.priceCollectorUC.GetNewObservedCoinPrices()
	if err != nil {
		logrus.Errorf("Background collect prices: %v", err)
	}
	// save new prices
	if _, err := p.priceCollectorUC.SaveCoinPrices(newPrices); err != nil {
		logrus.Errorf("Background save collected prices: %v", err)
	}
}
