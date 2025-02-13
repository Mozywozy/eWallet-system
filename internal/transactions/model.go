package transactions

import (
	"time"
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type TransactionType string

const (
	TransactionTopUp    TransactionType = "TOPUP"
	TransactionPurchase TransactionType = "PURCHASE"
	TransactionRefund   TransactionType = "REFUND"
)

type TransactionStatus string

const (
	StatusPending  TransactionStatus = "PENDING"
	StatusSuccess  TransactionStatus = "SUCCESS"
	StatusFailed   TransactionStatus = "FAILED"
	StatusReversed TransactionStatus = "REVERSED"
)

type AdditionalInfo map[string]interface{}

func (a AdditionalInfo) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *AdditionalInfo) Scan(value interface{}) error {
	if value == nil {
		*a = make(AdditionalInfo)
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to unmarshal JSON")
	}
	return json.Unmarshal(bytes, a)
}

type Transaction struct {
	ID                uint              `gorm:"primaryKey" json:"id"`
	UserID            uint              `gorm:"not null" json:"user_id"`
	Amount            float64           `gorm:"not null;default:0" json:"amount"`
	TransactionType   TransactionType   `gorm:"type:enum('TOPUP','PURCHASE','REFUND');not null" json:"transaction_type"`
	TransactionStatus TransactionStatus `gorm:"type:enum('PENDING','SUCCESS','FAILED','REVERSED');default:'PENDING'" json:"transaction_status"`
	Reference         string            `gorm:"type:varchar(255);not null" json:"reference"`
	Description       string            `gorm:"type:varchar(255);not null" json:"description"`
	AdditionalInfo    AdditionalInfo    `gorm:"type:json" json:"additional_info,omitempty"`
	CreatedAt         time.Time         `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time         `gorm:"autoUpdateTime" json:"updated_at"`
}
