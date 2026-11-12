package model

import (
	"errors"
	"time"
)

var (
	ErrInvalidAmount       = errors.New("amount must be greater than zero")
	ErrInsufficientBalance = errors.New("insufficient wallet balance")
	ErrWalletNotFound      = errors.New("wallet not found")
)

type Wallet struct {
	ID        string    `json:"id"`
	OwnerName string    `json:"owner_name"`
	Balance   int64     `json:"balance"`
	CreateAt  time.Time `json:"created_at"`
}

func (w *Wallet) Deposit(amount int64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}
	w.Balance += amount
	return nil
}

func (w *Wallet) Withdraw(amount int64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}
	if w.Balance < amount {
		return ErrInsufficientBalance
	}
	w.Balance -= amount
	return nil
}
