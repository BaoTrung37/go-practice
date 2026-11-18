package repository

import (
	"sync"
	"test/internal/model"
)

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

func (r *MemoryWalletRepository) Transfer(fromID string, toID string, amount int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	fromW, ok := r.wallets[fromID]
	if !ok {
		return model.ErrWalletNotFound
	}

	if fromW.Balance < amount {
		return model.ErrInsufficientBalance
	}

	toW, ok := r.wallets[toID]
	if !ok {
		return model.ErrWalletNotFound
	}

	fromW.Balance -= amount
	toW.Balance += amount
	return nil
}

var _ WalletRepository = (*MemoryWalletRepository)(nil)
