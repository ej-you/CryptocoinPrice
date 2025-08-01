package entity

// @description Coin price object
type Price struct {
	// price record uuid
	ID string `gorm:"id;primaryKey;type:uuid"`
	// coin uuid
	CoinID string `gorm:"coin_id;type:uuid"`
	// created at timestamp
	Timestamp int `gorm:"timestamp;not null"`

	// coin instance
	Coin *Coin `gorm:"foreignKey:CoinID;->"`
}
