package server

import (
	httpv1 "CryptocoinPrice/internal/app/controller/http/v1"
	"CryptocoinPrice/internal/app/controller/http/v1/coinmanage"
	repocoingecko "CryptocoinPrice/internal/app/repo/coingecko"
	repopg "CryptocoinPrice/internal/app/repo/pg"
	"CryptocoinPrice/internal/app/usecase"
)

// registerEndpointsV1 register all endpoints for 1st version of API.
func (s *Server) registerEndpointsV1() {
	// create repos
	coinRepoPG := repopg.NewCoinRepoPG(s.db)
	priceRepoDB := repopg.NewPriceRepoPG(s.db)
	priceRepoCoingecko := repocoingecko.NewPriceRepoCoingecko(s.cfg.App.CoingeckoAPIKey)
	// create usecases
	coinManageUC := usecase.NewCoinManageUC(coinRepoPG, priceRepoDB, priceRepoCoingecko)
	// create controllers
	coinManageController := coinmanage.NewController(coinManageUC, s.valid)
	// register endpoints
	apiV1 := s.fiberApp.Group("/api/v1")
	httpv1.RegisterCoinManageEndpoints(apiV1, coinManageController)
}
