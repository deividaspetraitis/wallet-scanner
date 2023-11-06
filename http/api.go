package http

import (
	"context"
	"net/http"
	stdhttp "net/http"
	"os"
	"syscall"
	"time"

	"github.com/deividaspetraitis/wallet-screener"
	db "github.com/deividaspetraitis/wallet-screener/database/immudb"

	"github.com/deividaspetraitis/wallet-screener/http/middleware"
	"github.com/deividaspetraitis/wallet-screener/log"

	immudb "github.com/codenotary/immudb/pkg/client"
	"github.com/gorilla/mux"
)

// App is the entrypoint into our application and what configures our context
// object for each of our http handlers. Feel free to add any configuration
// data/logic on this App struct
type App struct {
	API      *mux.Router
	shutdown chan os.Signal
}

// NewApp creates an App value that handle a set of routes for the application.
func NewApp(shutdown chan os.Signal) *App {
	api := App{
		API:      mux.NewRouter(),
		shutdown: shutdown,
	}
	return &api
}

// ServeHTTP API
func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.API.ServeHTTP(w, r)
}

// SignalShutdown is used to gracefully shutdown the app when an integrity
// issue is identified.
func (a *App) SignalShutdown() {
	a.shutdown <- syscall.SIGTERM
}

// API constructs an http.Handler with all application routes defined.
func API(shutdown chan os.Signal, cfg *Config, logger log.Logger, riskprovider walletscreener.WalletRiskScreeningProvider, immuclient immudb.ImmuClient) stdhttp.Handler {
	// =========================================================================
	// Construct the web app api which holds all routes as well as common Middleware.

	api := NewApp(shutdown)

	// =========================================================================
	// Construct and attach relevant handlers to web app api

	api.API.HandleFunc("/wallet/{address}/categories", GetRiskCategories(func(ctx context.Context, address string) ([]string, error) {
		return walletscreener.ScreenWalletRiskCategories(ctx, riskprovider, func(ctx context.Context, address string, categories []string) error {
			return db.StoreWalletRiskCategories(ctx, immuclient, address, categories)
		}, address)
	})).Methods(http.MethodPost)

	api.API.HandleFunc("/wallet/{address}/categories", GetRiskCategoriesHistory(func(ctx context.Context, address string) ([]*walletscreener.HistoricalRiskCategory, error) {
		return walletscreener.GetWalletRiskCategoriesHistory(ctx, func(ctx context.Context, address string) ([]string, []uint64, error) {
			return db.GetWalletRiskCategories(ctx, immuclient, address)
		}, address)
	})).Methods(http.MethodGet)

	// guard with request rate limiter
	api.API.Use(func(handler http.Handler) http.Handler {
		return middleware.RequestRate(cfg.Middleware.RateLimit, time.Minute, logger, handler)
	})

	router := mux.NewRouter()

	router.PathPrefix("/").Handler(api.API)

	return router
}
