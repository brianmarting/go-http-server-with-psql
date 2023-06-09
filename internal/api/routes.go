package api

import (
	"go-backend/internal/api/handlers"
	"go-backend/internal/facade/queue"
	"go-backend/internal/persistence/db"
	"go-backend/internal/service"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler interface {
	ServeHTTP(http.ResponseWriter, *http.Request)

	CreateAllRoutes() Handler
}

type handler struct {
	*chi.Mux

	cryptoHandler     handlers.CryptoHandler
	walletHandler     handlers.WalletHandler
	withdrawalHandler handlers.WithdrawalHandler
}

func NewHandler(store *db.Store) Handler {
	return handler{
		Mux: chi.NewMux(),
		cryptoHandler: handlers.NewCryptoHandler(
			service.NewCryptoService(store.CryptoStore),
		),
		walletHandler: handlers.NewWalletHandler(
			service.NewWalletService(store.WalletStore),
		),
		withdrawalHandler: handlers.NewWithdrawalHandler(
			queue.NewPublisher(),
			store.CryptoStore,
			store.WalletStore,
		),
	}
}

func (h handler) CreateAllRoutes() Handler {
	h.Route("/crypto", func(router chi.Router) {
		router.Get("/{id}", h.cryptoHandler.Get())
		router.Post("/", h.cryptoHandler.Create())
		router.Post("/{id}/delete", h.cryptoHandler.Delete())

		router.Route("/wallet", func(router chi.Router) {
			router.Get("/{id}", h.walletHandler.Get())
			router.Post("/", h.walletHandler.Create())
		})

		router.Route("/withdraw", func(router chi.Router) {
			router.Post("/", h.withdrawalHandler.Withdraw())
		})
	})

	return h
}
