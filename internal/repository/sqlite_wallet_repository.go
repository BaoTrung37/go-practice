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

var _ WalletRepository = (*SQLiteWalletRepository)(nil)
