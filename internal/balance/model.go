package balance

import (
	"time"
)

type Wallet struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"uniqueIndex" json:"user_id"`
	Balance   float64   `gorm:"not null;default:0" json:"balance"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type WalletTransaction struct {
	ID                   uint      `gorm:"primaryKey" json:"id"`
	WalletID             uint      `gorm:"not null" json:"wallet_id"`
	Amount               float64   `gorm:"not null" json:"amount"`
	WalletTransactionType string    `gorm:"column:wallet_transaction_type;type:enum('CREDIT','DEBIT');not null" json:"wallet_transaction_type"`
	Reference            string    `gorm:"type:varchar(100);not null" json:"reference"`
	CreatedAt            time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt            time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}