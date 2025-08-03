package coinmanage

import (
	"errors"
	"fmt"

	fiber "github.com/gofiber/fiber/v2"

	httpv1 "CryptocoinPrice/internal/app/controller/http/v1"
	"CryptocoinPrice/internal/app/usecase"
	"CryptocoinPrice/internal/pkg/validator"
)

var _ httpv1.CoinManageController = (*Controller)(nil)

// Controller is a HTTP-controller for coin manage usecase.
type Controller struct {
	uc    usecase.CoinManageUsecase
	valid validator.Validator
}

// NewController returns new coin manage controller.
func NewController(uc usecase.CoinManageUsecase,
	valid validator.Validator) *Controller {

	return &Controller{
		uc:    uc,
		valid: valid,
	}
}

// AddObserve appends coin to observed list.
//
//	@summary		Добавление криптовалюты в список наблюдения
//	@description	Добавление криптовалюты в список наблюдения.
//	@router			/currency/add [post]
//	@id				observe-coin
//	@tags			currency
//	@param			Coin	body	coinObservedInput	true	"Название криптовалюты"
//	@success		204		"Успешное добавление в список наблюдения"
//	@failure		400		"Невалидное тело запроса"
//	@failure		404		"Криптовалюта с таким названием не существует"
func (c *Controller) AddObserve(ctx *fiber.Ctx) error {
	bodyData := &coinObservedInput{}
	// parse body
	if err := ctx.BodyParser(bodyData); err != nil {
		return fmt.Errorf("parse body: %w", err)
	}
	// validate parsed data
	if err := c.valid.Validate(bodyData); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, validateErrStr(err))
	}

	// observe coin
	if _, err := c.uc.ObserveCoin(bodyData.Symbol); err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return err
	}
	return ctx.Status(fiber.StatusNoContent).Send(nil)
}

// RemoveObserve removes coin from observed list.
//
//	@summary		Удаление криптовалюты из списка наблюдения
//	@description	Удаление криптовалюты из списка наблюдения.
//	@router			/currency/remove [delete]
//	@id				disable-observe-coin
//	@tags			currency
//	@param			Coin	body	coinObservedInput	true	"Название криптовалюты"
//	@success		204		"Успешное добавление в список наблюдения"
//	@failure		400		"Невалидное тело запроса"
//	@failure		404		"Криптовалюта с таким названием не была добавлена в список наблюдения ранее"
func (c *Controller) RemoveObserve(ctx *fiber.Ctx) error {
	bodyData := &coinObservedInput{}
	// parse body
	if err := ctx.BodyParser(bodyData); err != nil {
		return fmt.Errorf("parse body: %w", err)
	}
	// validate parsed data
	if err := c.valid.Validate(bodyData); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, validateErrStr(err))
	}

	// observe coin
	if _, err := c.uc.DisableObserveCoin(bodyData.Symbol); err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return err
	}
	return ctx.Status(fiber.StatusNoContent).Send(nil)
}

// GetPrice returns nearest price for coin and timestamp.
//
//	@summary		Получение цены криптовалюты
//	@description	Получение цены криптовалюты.
//	@router			/currency/price [get]
//	@id				get-coin-price
//	@tags			currency
//	@param			coin		query		string	true	"Название криптовалюты и время"
//	@param			timestamp	query		int64	true	"Время в UNIX-формате"
//	@success		200			{object}	coinPriceOutput
//	@failure		400			"Невалидное тело запроса"
//	@failure		404			"Ни одна цена криптовалюты не найдена"
func (c *Controller) GetPrice(ctx *fiber.Ctx) error {
	bodyData := &coinPriceInput{}
	// parse body
	if err := ctx.QueryParser(bodyData); err != nil {
		return fmt.Errorf("parse body: %w", err)
	}
	// validate parsed data
	if err := c.valid.Validate(bodyData); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, validateErrStr(err))
	}

	// get coin price
	price, err := c.uc.GetNearestPrice(bodyData.Symbol, bodyData.Timestamp)
	if err != nil && errors.Is(err, usecase.ErrNotFound) {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}
	if err != nil {
		return fmt.Errorf("get coin price: %w", err)
	}
	outputPrice := coinPriceOutput{
		Symbol:    bodyData.Symbol,
		Timestamp: price.Timestamp,
		Price:     price.Price,
	}
	return ctx.Status(fiber.StatusOK).JSON(outputPrice)
}

// validateErrStr returns string with validate error.
func validateErrStr(err error) string {
	return "validate data: " + err.Error()
}
