package repository

import (
	"sync"
	"test/internal/model"
)

type WalletRepository interface {
	Save(wallet *model.Wallet) error
	FindByID(id string) (*model.Wallet, error)
}

type MemoryWalletRepository struct {
	wallets map[string]*model.Wallet
	mu      sync.RWMutex
}

func NewMemoryWalletRepository() *MemoryWalletRepository {
	return &MemoryWalletRepository{
		wallets: make(map[string]*model.Wallet),
	}
}

func (r *MemoryWalletRepository) Save(w *model.Wallet) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	clone := *w
	r.wallets[w.ID] = &clone
	return nil
}

func (r *MemoryWalletRepository) FindByID(id string) (*model.Wallet, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.wallets[id] == nil {
		return nil, model.ErrWalletNotFound
	}
	clone := *r.wallets[id]
	return &clone, nil
}

var _ WalletRepository = (*MemoryWalletRepository)(nil)
