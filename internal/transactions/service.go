package transactions

import (
	"errors"
)

type TransactionService interface {
	InitiateTransaction(userID uint, amount float64, txType TransactionType, reference string, description string, additionalInfo AdditionalInfo) error
	UpdateTransaction(reference string, status TransactionStatus) error
	GetTransactionByReference(reference string) (*Transaction, error)
}

type transactionService struct {
	txRepo TransactionRepository
}

func NewTransactionService(repo TransactionRepository) TransactionService {
	return &transactionService{txRepo: repo}
}

func (s *transactionService) InitiateTransaction(userID uint, amount float64, txType TransactionType, reference string, description string, additionalInfo AdditionalInfo) error {
	if amount <= 0 {
		return errors.New("jumlah transaksi tidak valid")
	}

	transaction := Transaction{
		UserID:            userID,
		Amount:            amount,
		TransactionType:   txType,
		TransactionStatus: StatusPending,
		Reference:         reference,
		Description:       description,
		AdditionalInfo:    additionalInfo,
	}

	return s.txRepo.CreateTransaction(&transaction)
}

func (s *transactionService) UpdateTransaction(reference string, status TransactionStatus) error {
	transaction, err := s.txRepo.GetTransactionByReference(reference)
	if err != nil {
		return errors.New("transaksi tidak ditemukan")
	}

	err = s.txRepo.UpdateTransactionStatus(reference, status)
	if err != nil {
		return err
	}

	if status == StatusSuccess {
		err = s.txRepo.AdjustBalance(transaction.UserID, transaction.TransactionType, transaction.Amount, transaction.Reference)
		if err != nil {
			return errors.New("gagal memperbarui saldo user")
		}
	}

	return nil
}

func (s *transactionService) GetTransactionByReference(reference string) (*Transaction, error) {
	return s.txRepo.GetTransactionByReference(reference)
}
