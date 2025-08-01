// Package entity contains all app entities.
package entity

// @description Coin object
type Coin struct {
	// coin uuid
	ID string `gorm:"id;primaryKey;type:uuid"`
	// coin name
	Symbol string `gorm:"symbol;not null;uniqueIndex"`
	// true if coin is observed
	Observed bool `gorm:"observed;not null"`
}
