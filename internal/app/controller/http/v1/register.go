// Package http/v1 is a first version of HTTP-controllers.
// It provides registers for HTTP-routes.
// Controllers with handlers for them is in subpackages.
package v1

import (
	fiber "github.com/gofiber/fiber/v2"
)

type CoinManageController interface {
	AddObserve(ctx *fiber.Ctx) error
	RemoveObserve(ctx *fiber.Ctx) error
	GetPrice(ctx *fiber.Ctx) error
}

// RegisterCoinManageEndpoints registers all endpoints for coin manage controller.
func RegisterCoinManageEndpoints(router fiber.Router, controller CoinManageController) {
	currencyPrefix := router.Group("/currency")

	currencyPrefix.Post("/add", controller.AddObserve)
	currencyPrefix.Delete("/remove", controller.RemoveObserve)
	currencyPrefix.Get("/price", controller.GetPrice)
}
