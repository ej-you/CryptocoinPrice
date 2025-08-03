package coinmanage

// @description Input to add/remove coin to/from observed list..
type coinObservedInput struct {
	// Coin short name
	Symbol string `json:"coin" validate:"required,alpha" example:"btc"`
}

// @description Input to get coin price at timestamp.
type coinPriceInput struct {
	// Coin short name
	Symbol string `query:"coin" validate:"required,alpha" example:"btc"`
	// Unix timestamp
	Timestamp int64 `query:"timestamp" validate:"required,min=0" example:"1736500490"`
}

// @description Output for gotten coin price at timestamp.
type coinPriceOutput struct {
	// Coin short name
	Symbol string `json:"coin" validate:"required,alpha" example:"btc"`
	// Unix timestamp
	Timestamp int64 `json:"timestamp" validate:"required,min=0" example:"1754045773"`
	// Coin price
	Price string `json:"price" example:"114818"`
}
