package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"test/internal/middleware"
	"test/internal/model"
	"test/internal/service"
	"test/internal/worker"
)

type WalletHandler struct {
	svc  *service.WalletService
	pool *worker.WorkerPool
}

func NewWalletHandler(svc *service.WalletService, pool *worker.WorkerPool) *WalletHandler {
	return &WalletHandler{svc: svc, pool: pool}
}

func (h *WalletHandler) RegisterRoutes(mux *http.ServeMux) {
	// Public Routes
	mux.HandleFunc("POST /wallets", h.CreateWallet)
	mux.HandleFunc("GET /wallets/{id}", h.GetWallet)

	// Protected Routes
	mux.Handle("POST /wallets/{id}/deposit", middleware.Auth(http.HandlerFunc(h.Deposit)))
	mux.Handle("POST /wallets/transfer", middleware.Auth(http.HandlerFunc(h.Transfer)))
}

func (h *WalletHandler) CreateWallet(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID        string `json:"id"`
		OwnerName string `json:"owner_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	wallet, err := h.svc.CreateWallet(req.ID, req.OwnerName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(wallet)

}

func (h *WalletHandler) GetWallet(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	wallet, err := h.svc.GetWallet(id)

	if err != nil {
		if errors.Is(err, model.ErrWalletNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(wallet)
}

func (h *WalletHandler) Deposit(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var req struct {
		Amount int64 `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	wallet, err := h.svc.Deposit(id, req.Amount)
	if err != nil {
		if errors.Is(err, model.ErrWalletNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if errors.Is(err, model.ErrInvalidAmount) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(wallet)
}

func (h *WalletHandler) Transfer(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		FromID string `json:"from_id"`
		ToID   string `json:"to_id"`
		Amount int64  `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if !strings.EqualFold(userID, req.FromID) {
		http.Error(w, "forbidden: you can only transfer from your own wallet", http.StatusForbidden)
		return
	}

	err := h.svc.Transfer(req.FromID, req.ToID, req.Amount)
	if err != nil {
		if errors.Is(err, model.ErrWalletNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if errors.Is(err, model.ErrInvalidAmount) || errors.Is(err, model.ErrInsufficientBalance) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Chuyển tiền thành công -> Bắn task ngầm gửi thông báo qua Worker Pool!
	if h.pool != nil {
		h.pool.Submit(func() {
			// Giả lập tác vụ gửi email/SMS thông báo biến động số dư chạy ngầm
			fmt.Printf("🔔 [Worker] Notification: Transferred %d from %s to %s successfully!\n",
				req.Amount, req.FromID, req.ToID)
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "transfer successful",
	})
}
