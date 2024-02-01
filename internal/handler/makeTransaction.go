package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/The-Gleb/EWallet/internal/domain/service"
	maketransaction_usecase "github.com/The-Gleb/EWallet/internal/domain/usecase/makeTransaction.go"
	"github.com/The-Gleb/EWallet/internal/handler/dto"
	"github.com/go-chi/chi"
)

const (
	makeTransactionURL = "/api/v1/wallet/{walletId}/send"
)

type MakeTransactionUsecase interface {
	MakeTransaction(ctx context.Context, dto maketransaction_usecase.MakeTransactionDTO) error
}

type makeTransactionHandler struct {
	usecase MakeTransactionUsecase
}

func NewMakeTransactionHandler(usecase MakeTransactionUsecase) *makeTransactionHandler {
	return &makeTransactionHandler{usecase: usecase}
}

func (h *makeTransactionHandler) AddToRouter(r *chi.Mux) {
	r.Post(makeTransactionURL, h.MakeTransaction)

}

func (h *makeTransactionHandler) MakeTransaction(rw http.ResponseWriter, r *http.Request) {

	newTransactionDTO := dto.NewTransactionDTO{}

	err := json.NewDecoder(r.Body).Decode(&newTransactionDTO)
	if err != nil {
		slog.Error(err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	senderString := chi.URLParam(r, "walletId")
	sender, err := strconv.ParseInt(senderString, 10, 64)
	if err != nil {
		slog.Error(err.Error())
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	usecaseDTO := maketransaction_usecase.MakeTransactionDTO{
		From:   sender,
		To:     newTransactionDTO.To,
		Amount: newTransactionDTO.Amount,
	}

	err = h.usecase.MakeTransaction(r.Context(), usecaseDTO)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrSenderWalletNotFound):
			http.Error(rw, err.Error(), http.StatusNotFound)
			return
		case errors.Is(err, service.ErrInsufficientFunds) ||
			errors.Is(err, service.ErrReceiverWalletNotFound):
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		default:
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	rw.WriteHeader(http.StatusOK)

}
