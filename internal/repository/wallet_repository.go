package repository

import (
	"test/internal/model"
)

type WalletRepository interface {
	Save(wallet *model.Wallet) error
	FindByID(id string) (*model.Wallet, error)
	Transfer(fromID string, toID string, amount int64) error
}
