package transactions

import (
	"errors"
	"ewallet-engine/internal/balance"
	"log"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	CreateTransaction(tx *Transaction) error
	UpdateTransactionStatus(reference string, status TransactionStatus) error
	GetTransactionByReference(reference string) (*Transaction, error)
	AdjustBalance(userID uint, txType TransactionType, amount float64, reference string) error
}

type transactionRepository struct {
	DB *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{DB: db}
}

func (r *transactionRepository) CreateTransaction(tx *Transaction) error {
	return r.DB.Create(tx).Error
}

func (r *transactionRepository) UpdateTransactionStatus(reference string, status TransactionStatus) error {
	return r.DB.Model(&Transaction{}).Where("reference = ?", reference).Update("transaction_status", status).Error
}

func (r *transactionRepository) GetTransactionByReference(reference string) (*Transaction, error) {
	var tx Transaction
	err := r.DB.Where("reference = ?", reference).First(&tx).Error
	if err != nil {
		return nil, err
	}
	return &tx, nil
}

func (r *transactionRepository) AdjustBalance(userID uint, txType TransactionType, amount float64, reference string) error {
	var wallet balance.Wallet

	err := r.DB.Where("user_id = ?", userID).First(&wallet).Error
	if err != nil {
		log.Printf("ERROR: Wallet tidak ditemukan untuk user_id %d, error: %v", userID, err)
		return err
	}

	var walletTxType string
	if txType == TransactionTopUp || txType == TransactionRefund {
		wallet.Balance += amount
		walletTxType = "CREDIT"
	} else if txType == TransactionPurchase {
		if wallet.Balance < amount {
			log.Printf("ERROR: Saldo tidak mencukupi untuk user_id %d", userID)
			return errors.New("saldo tidak mencukupi")
		}
		wallet.Balance -= amount
		walletTxType = "DEBIT"
	}

	err = r.DB.Save(&wallet).Error
	if err != nil {
		log.Printf("ERROR: Gagal menyimpan saldo user_id %d, error: %v", userID, err)
		return err
	}

	walletTransaction := balance.WalletTransaction{
		WalletID:             wallet.ID,
		Amount:               amount,
		WalletTransactionType: walletTxType,
		Reference:            reference,
	}

	err = r.DB.Create(&walletTransaction).Error
	if err != nil {
		log.Printf("ERROR: Gagal menyimpan transaksi wallet untuk user_id %d, error: %v", userID, err)
		return errors.New("gagal menyimpan transaksi saldo")
	}

	log.Printf("SUCCESS: Saldo user_id %d berhasil diperbarui, transaksi disimpan.", userID)
	return nil
}
