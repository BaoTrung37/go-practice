package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"test/internal/model"
	"test/internal/service"
)

type WalletHandler struct {
	svc *service.WalletService
}

func NewWalletHandler(svc *service.WalletService) *WalletHandler {
	return &WalletHandler{svc: svc}
}

func (h *WalletHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /wallets", h.CreateWallet)
	mux.HandleFunc("GET /wallets/{id}", h.GetWallet)
	mux.HandleFunc("POST /wallets/{id}/deposit", h.Deposit)
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
