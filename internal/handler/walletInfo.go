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
	walletInfoURL = "/api/v1/wallet/{walletId}"
)

type WalletInfoUsecase interface {
	GetWalletInfo(ctx context.Context, ID int64) (entity.Wallet, error)
}

type walletInfoHandler struct {
	usecase WalletInfoUsecase
}

func NewWalletInfoHandler(usecase WalletInfoUsecase) *walletInfoHandler {
	return &walletInfoHandler{usecase: usecase}
}

func (h *walletInfoHandler) AddToRouter(r *chi.Mux) {
	r.Route(walletInfoURL, func(r chi.Router) {
		r.Get("/", h.GetWalletInfo)
		// r.Route(registerURL, func(r chi.Router) {

		// })
	})

}

func (h *walletInfoHandler) GetWalletInfo(rw http.ResponseWriter, r *http.Request) {

	walletIDString := chi.URLParam(r, "walletId")
	walletID, err := strconv.ParseInt(walletIDString, 10, 64)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	walletInfo, err := h.usecase.GetWalletInfo(r.Context(), walletID)
	if err != nil {
		if errors.Is(err, service.ErrWalletNotFound) {
			http.Error(rw, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(walletInfo)
	if err != nil {
		slog.Error(err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
	// set application/json
	rw.Write(b)

}
