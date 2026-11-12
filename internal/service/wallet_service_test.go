package service

import (
	"errors"
	"test/internal/model"
	"test/internal/repository"
	"testing"
)

func TestWalletService_DepositAndWithdraw(t *testing.T) {
	repo := repository.NewMemoryWalletRepository()
	svc := NewWalletService(repo)

	// 1: Create Wallet
	_, err := svc.CreateWallet("w1", "Trung")
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	// 2: Deposit 100_000
	w, err := svc.Deposit("w1", 100000)
	if err != nil || w.Balance != 100000 {
		t.Fatalf("expect balance 100000, got %d (err: %v)", w.Balance, err)
	}

	// 3: Withdraw 40000

	w, err = svc.Withdraw("w1", 40000)
	if err != nil || w.Balance != 60000 {
		t.Fatalf("expect balance 60000, got %d (err: %v)", w.Balance, err)
	}

	// 4: Withdraw 100000
	w, err = svc.Withdraw("w1", 100000)
	if !errors.Is(err, model.ErrInsufficientBalance) {
		t.Fatalf("expect error %v, got %v", model.ErrInsufficientBalance, err)
	}
}
