package repository

import (
	"errors"
	"path/filepath"
	"testing"
	"time"

	"test/internal/model"
)

func TestSQLiteWalletRepository_Transfer_Rollback(t *testing.T) {
	// 1. Tạo file SQLite tạm thời trong thư mục temp của bài test
	dbPath := filepath.Join(t.TempDir(), "test_wallets.db")
	repo, err := NewSqliteWalletRepository(dbPath)
	if err != nil {
		t.Fatalf("failed to create repo: %v", err)
	}

	// 2. Tạo ví A (100.000đ) và ví B (50.000đ)
	wA := &model.Wallet{ID: "A", OwnerName: "Alice", Balance: 100000, CreateAt: time.Now()}
	wB := &model.Wallet{ID: "B", OwnerName: "Bob", Balance: 50000, CreateAt: time.Now()}
	_ = repo.Save(wA)
	_ = repo.Save(wB)

	// Kịch bản 1: A chuyển cho một ví KHÔNG TỒN TẠI ("C")
	err = repo.Transfer("A", "C", 40000)
	if !errors.Is(err, model.ErrWalletNotFound) {
		t.Fatalf("expected ErrWalletNotFound, got %v", err)
	}

	// CHỨNG MINH ROLLBACK: Kiểm tra ví A xem có bị trừ tiền oan không?
	checkA, _ := repo.FindByID("A")
	if checkA.Balance != 100000 {
		t.Fatalf("Rollback failed! Wallet A lost money: expected 100000, got %d", checkA.Balance)
	}

	// Kịch bản 2: A chuyển quá số dư (chuyển 150.000đ trong khi chỉ có 100.000đ)
	err = repo.Transfer("A", "B", 150000)
	if !errors.Is(err, model.ErrInsufficientBalance) {
		t.Fatalf("expected ErrInsufficientBalance, got %v", err)
	}

	// CHỨNG MINH ROLLBACK: Cả A và B phải giữ nguyên số dư ban đầu
	checkA, _ = repo.FindByID("A")
	checkB, _ := repo.FindByID("B")
	if checkA.Balance != 100000 || checkB.Balance != 50000 {
		t.Fatalf("Rollback failed on insufficient balance! A: %d, B: %d", checkA.Balance, checkB.Balance)
	}

	// Kịch bản 3: Chuyển tiền hợp lệ (A chuyển cho B 40.000đ)
	err = repo.Transfer("A", "B", 40000)
	if err != nil {
		t.Fatalf("expected transfer success, got %v", err)
	}

	checkA, _ = repo.FindByID("A")
	checkB, _ = repo.FindByID("B")
	if checkA.Balance != 60000 || checkB.Balance != 90000 {
		t.Fatalf("Transfer amounts incorrect! A: %d, B: %d", checkA.Balance, checkB.Balance)
	}
}
