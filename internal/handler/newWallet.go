package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/The-Gleb/EWallet/internal/domain/entity"
	"github.com/go-chi/chi"
)

const (
	newWalletURL = "/api/v1/wallet"
)

type NewWalletUsecase interface {
	NewWallet(ctx context.Context) (entity.Wallet, error)
}

type newWalletHandler struct {
	usecase NewWalletUsecase
}

func NewNewWalletHandler(usecase NewWalletUsecase) *newWalletHandler {
	return &newWalletHandler{usecase: usecase}
}

func (h *newWalletHandler) AddToRouter(r *chi.Mux) {
	r.Post(newWalletURL, h.NewWallet)
}

func (h *newWalletHandler) NewWallet(rw http.ResponseWriter, r *http.Request) {

	walletInfo, err := h.usecase.NewWallet(r.Context())
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(walletInfo)
	if err != nil {
		slog.Error(err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Write(body)

}
