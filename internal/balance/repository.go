package balance

import (
	"errors"

	"gorm.io/gorm"
)

type BalanceRepository interface {
	GetBalance(userID uint) (float64, error)
	AdjustBalance(userID uint, amount float64, txType string, reference string) error
	RecordTransaction(walletID uint, txType string, amount float64, reference string) error
}

type balanceRepository struct {
	DB *gorm.DB
}

func NewBalanceRepository(db *gorm.DB) BalanceRepository {
	return &balanceRepository{DB: db}
}

func (r *balanceRepository) GetBalance(userID uint) (float64, error) {
	var wallet Wallet
	err := r.DB.Where("user_id = ?", userID).First(&wallet).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		return 0, err
	}
	return wallet.Balance, nil
}

func (r *balanceRepository) AdjustBalance(userID uint, amount float64, txType string, reference string) error {
	var wallet Wallet

	err := r.DB.Where("user_id = ?", userID).FirstOrCreate(&wallet, Wallet{UserID: userID}).Error
	if err != nil {
		return err
	}

	if txType == "CREDIT" {
		wallet.Balance += amount
	} else if txType == "DEBIT" {
		if wallet.Balance < amount {
			return errors.New("saldo tidak mencukupi untuk transaksi ini")
		}
		wallet.Balance -= amount
	} else {
		return errors.New("jenis transaksi tidak valid")
	}

	err = r.DB.Save(&wallet).Error
	if err != nil {
		return err
	}

	return r.RecordTransaction(wallet.ID, txType, amount, reference)
}

func (r *balanceRepository) RecordTransaction(walletID uint, txType string, amount float64, reference string) error {
	tx := WalletTransaction{
		WalletID:             walletID,
		WalletTransactionType: txType, 
		Amount:               amount,
		Reference:            reference,
	}
	return r.DB.Create(&tx).Error
}
