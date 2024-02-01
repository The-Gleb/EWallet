package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/The-Gleb/EWallet/internal/domain/entity"
	"github.com/The-Gleb/EWallet/internal/domain/service"
	"github.com/go-chi/chi"
)

const (
	txHistoryURL = "/api/v1/wallet/{walletId}/history"
)

type TxHistoryUsecase interface {
	GetTxHistory(ctx context.Context, walletID int64) ([]entity.Transaction, error)
}

type txHistoryHandler struct {
	usecase TxHistoryUsecase
}

func NewTxHistoryHandler(usecase TxHistoryUsecase) *txHistoryHandler {
	return &txHistoryHandler{usecase: usecase}
}

func (h *txHistoryHandler) AddToRouter(r *chi.Mux) {
	r.Get(txHistoryURL, h.GetTxHistory)
}

func (h *txHistoryHandler) GetTxHistory(rw http.ResponseWriter, r *http.Request) {

	walletIDString := chi.URLParam(r, "walletId")
	walletID, err := strconv.ParseInt(walletIDString, 10, 64)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	txHistory, err := h.usecase.GetTxHistory(r.Context(), walletID)
	if err != nil {
		if errors.Is(err, service.ErrWalletNotFound) {
			http.Error(rw, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(txHistory)
	if err != nil {
		slog.Error(err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	// set application/json
	rw.Write(b)

}
