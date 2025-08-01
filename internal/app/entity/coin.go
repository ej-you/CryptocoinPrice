// Package entity contains all app entities.
package entity

// Coin is a coin DB object.
type Coin struct {
	// coin uuid
	ID string `gorm:"id;primaryKey;type:uuid"`
	// coin name
	Symbol string `gorm:"symbol;not null;uniqueIndex"`
	// true if coin is observed
	Observed bool `gorm:"observed;not null"`
}

// CoinPartial is a coin object with all optional fields.
type CoinPartial struct {
	// coin name
	Symbol *string `gorm:"symbol"`
	// true if coin is observed
	Observed *bool `gorm:"observed"`
}

// CoinList is a slice of coins.
type CoinList []Coin
