package balance

import "errors"

type BalanceService interface {
	GetUserBalance(userID uint) (float64, error)
	ProcessBalanceTransaction(userID uint, amount float64, txType string, reference string) error
}

type balanceService struct {
	repo BalanceRepository
}

func NewBalanceService(repo BalanceRepository) BalanceService {
	return &balanceService{repo: repo}
}

func (s *balanceService) GetUserBalance(userID uint) (float64, error) {
	return s.repo.GetBalance(userID)
}

func (s *balanceService) ProcessBalanceTransaction(userID uint, amount float64, txType string, reference string) error {
	if amount <= 0 {
		return errors.New("jumlah transaksi tidak valid")
	}

	return s.repo.AdjustBalance(userID, amount, txType, reference)
}
