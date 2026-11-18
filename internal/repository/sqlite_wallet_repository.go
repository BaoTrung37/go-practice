package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"test/internal/model"
	"time"

	_ "modernc.org/sqlite"
)

type SQLiteWalletRepository struct {
	db *sql.DB
}

func NewSqliteWalletRepository(dbPath string) (*SQLiteWalletRepository, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite db : %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping sqlite db: %w", err)
	}

	repo := &SQLiteWalletRepository{db: db}

	// Tự động tạo bảng nếu chưa có
	if err := repo.initSchema(); err != nil {
		return nil, err
	}
	return repo, nil

}

func (r *SQLiteWalletRepository) initSchema() error {
	query := `
	create table if not exists wallets(
		id text primary key,
		owner_name text not null,
		balance integer not null,
		created_at datetime not null
	);`

	_, err := r.db.Exec(query)
	return err
}

func (r *SQLiteWalletRepository) Save(w *model.Wallet) error {
	query := `
	insert into wallets (id, owner_name, balance, created_at)
	values (?,?,?,?)
	on conflict(id) do update set
		balance = excluded.balance;`

	_, err := r.db.Exec(query, w.ID, w.OwnerName, w.Balance, w.CreateAt)
	return err
}

func (r *SQLiteWalletRepository) FindByID(id string) (*model.Wallet, error) {
	query := `select id, owner_name, balance, created_at from wallets where id = ?`
	var w model.Wallet
	err := r.db.QueryRow(query, id).Scan(&w.ID, &w.OwnerName, &w.Balance, &w.CreateAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrWalletNotFound
		}
		return nil, err
	}
	return &w, nil
}

func (r *SQLiteWalletRepository) Transfer(fromID string, toID string, amount int64) error {
	if amount <= 0 {
		return model.ErrInvalidAmount
	}

	// 1. Mo transaction
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback() // Dam bao rollback neu co bat ki dong nao ben duoi bi loi

	// 2. Kiem tra vi gui

	var fromBalance int64
	err = tx.QueryRow("select balance from wallets where id = ?", fromID).Scan(&fromBalance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.ErrWalletNotFound
		}
		return err
	}
	if fromBalance < amount {
		return model.ErrInsufficientBalance
	}

	// 3. Kiem tra vi dich

	var toBalance int64
	err = tx.QueryRow("select balance from wallets where id = ?", toID).Scan(&toBalance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.ErrWalletNotFound
		}
		return err
	}

	// 4. Tru tien
	_, err = tx.Exec("update wallets set balance = balance - ? where id = ?", amount, fromID)
	if err != nil {
		return err
	}

	// 5. Cong tien

	_, err = tx.Exec("update wallets set balance = balance + ? where id = ?", amount, toID)
	if err != nil {
		return err
	}
	return tx.Commit()
}

var _ WalletRepository = (*SQLiteWalletRepository)(nil)
