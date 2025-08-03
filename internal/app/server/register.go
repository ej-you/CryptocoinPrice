package server

import (
	"gorm.io/gorm"

	"CryptocoinPrice/config"
	httpv1 "CryptocoinPrice/internal/app/controller/http/v1"
	"CryptocoinPrice/internal/app/controller/http/v1/coinmanage"
	repocoingecko "CryptocoinPrice/internal/app/repo/coingecko"
	repopg "CryptocoinPrice/internal/app/repo/pg"
	"CryptocoinPrice/internal/app/usecase"
	"CryptocoinPrice/internal/pkg/validator"
)

// registerEndpointsV1 register all endpoints for 1st version of API.
func (s *Server) registerEndpointsV1(cfg *config.Config, db *gorm.DB, valid validator.Validator) {
	// create repos
	coinRepoPG := repopg.NewCoinRepoPG(db)
	priceRepoDB := repopg.NewPriceRepoPG(db)
	priceRepoCoingecko := repocoingecko.NewPriceRepoCoingecko(cfg.App.CoingeckoAPIKey)
	// create usecases
	coinManageUC := usecase.NewCoinManageUC(coinRepoPG, priceRepoDB, priceRepoCoingecko)
	// create controllers
	coinManageController := coinmanage.NewController(coinManageUC, valid)
	// register endpoints
	apiV1 := s.fiberApp.Group("/api/v1")
	httpv1.RegisterCoinManageEndpoints(apiV1, coinManageController)
}
