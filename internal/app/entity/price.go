package entity

// Price is a coin price object
type Price struct {
	// price record uuid
	ID string `gorm:"id;primaryKey;type:uuid"`
	// coin uuid
	CoinID string `gorm:"coin_id;type:uuid"`
	// coin price
	Price string `gorm:"price;not null"`
	// created at timestamp
	Timestamp int64 `gorm:"timestamp;not null"`

	// coin instance
	Coin *Coin `gorm:"foreignKey:CoinID;->"`
}

// PriceList is a slice of coins' prices.
type PriceList []Price

// CoinPriceAPI ia a coin prise parsed from API.
type CoinPriceAPI struct {
	// coin symbol
	Symbol string
	// coin price
	Price float64
	// last update time in unix format
	LastUpdate int64
}

// CoinPriceAPIList is a slice of coins' prices from API.
type CoinPriceAPIList []CoinPriceAPI
