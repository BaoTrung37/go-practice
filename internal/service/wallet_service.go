package service

import (
	"test/internal/model"
	"test/internal/repository"
	"time"
)

type WalletService struct {
	repo repository.WalletRepository
}

func NewWalletService(repo repository.WalletRepository) *WalletService {
	return &WalletService{
		repo: repo,
	}
}

func (s *WalletService) GetWallet(id string) (*model.Wallet, error) {
	w, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (s *WalletService) CreateWallet(id string, ownerName string) (*model.Wallet, error) {
	w := &model.Wallet{
		ID:        id,
		OwnerName: ownerName,
		Balance:   0,
		CreateAt:  time.Now(),
	}
	err := s.repo.Save(w)
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (s *WalletService) Deposit(id string, amount int64) (*model.Wallet, error) {
	w, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	err = w.Deposit(amount)
	if err != nil {
		return nil, err
	}

	err = s.repo.Save(w)
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (s *WalletService) Withdraw(id string, amount int64) (*model.Wallet, error) {
	w, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	err = w.Withdraw(amount)
	if err != nil {
		return nil, err
	}

	err = s.repo.Save(w)
	if err != nil {
		return nil, err
	}
	return w, nil
}
