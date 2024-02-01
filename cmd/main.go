package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	storage "github.com/The-Gleb/EWallet/internal/adapter/db/postgres"
	"github.com/The-Gleb/EWallet/internal/config"
	"github.com/The-Gleb/EWallet/internal/domain/service"
	maketransaction_usecase "github.com/The-Gleb/EWallet/internal/domain/usecase/makeTransaction.go"
	newwallet_usecase "github.com/The-Gleb/EWallet/internal/domain/usecase/newWallet.go"
	txhistory_usecase "github.com/The-Gleb/EWallet/internal/domain/usecase/txHistory.go"
	walletinfogo_usecase "github.com/The-Gleb/EWallet/internal/domain/usecase/walletInfo.go"
	"github.com/The-Gleb/EWallet/internal/handler"
	"github.com/The-Gleb/EWallet/internal/logger"
	postgresql "github.com/The-Gleb/EWallet/pkg/client/postgres"
	"github.com/go-chi/chi"
)

// https://github.com/The-Gleb/EWallet

func main() {

	cfg := config.BuildConfig("config.yml")
	logger.Initialize(cfg.LogLevel)
	slog.Debug("here is config", "config", cfg)

	dbClient, err := postgresql.NewClient(context.Background(), 3, cfg.DB)
	if err != nil {
		slog.Error(err.Error())
	}

	walletStorage := storage.NewWalletStorage(dbClient)
	transactionStorage := storage.NewTransactionStorage(dbClient)

	walletService := service.NewWalletService(walletStorage)
	transactionService := service.NewTransactionService(transactionStorage)

	newWalletUsecase := newwallet_usecase.NewNewWalletUsecase(walletService)
	makeTransactionUsecase := maketransaction_usecase.NewMakeTransactionUsecase(walletService, transactionService)
	walletInfoUsecase := walletinfogo_usecase.NewWalletInfoUsecase(walletService)
	txHistoryUsecase := txhistory_usecase.NewTxHistoryUsecase(transactionService, walletService)

	newWalletHandler := handler.NewNewWalletHandler(newWalletUsecase)
	makeTransactionHandler := handler.NewMakeTransactionHandler(makeTransactionUsecase)
	walletInfoHandler := handler.NewWalletInfoHandler(walletInfoUsecase)
	txHistoryHandler := handler.NewTxHistoryHandler(txHistoryUsecase)

	r := chi.NewRouter()

	newWalletHandler.AddToRouter(r)
	makeTransactionHandler.AddToRouter(r)
	walletInfoHandler.AddToRouter(r)
	txHistoryHandler.AddToRouter(r)

	s := http.Server{
		Addr:    "localhost:8080",
		Handler: r,
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		ServerShutdownSignal := make(chan os.Signal, 1)
		signal.Notify(ServerShutdownSignal, syscall.SIGINT)
		<-ServerShutdownSignal
		s.Shutdown(context.Background())
		slog.Info("server shutdown")
	}()

	slog.Info("starting server")
	if err := s.ListenAndServe(); err != nil {
		slog.Error("error running server", "error", err)
	}

}
