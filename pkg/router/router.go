package router

import (
	"fmt"
	"net/http"

	"github.com/me/finance/config"
	"github.com/me/finance/pkg/db"
	"github.com/me/finance/pkg/repository"
	"github.com/rs/cors"
	"github.com/sagikazarmark/slog-shim"
)

func InitializeRoutes() {
	db, err := db.NewDB()
	if err != nil {
		slog.Error("error trying to connect to database")
		return
	}

	rep := repository.NewRepository(db)

	mux := http.NewServeMux()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // Allow your frontend origin
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false, // Important for cookies, authorization headers with CORS
		Debug:            true, // Enable for debugging CORS issues
	})

	PersonRoutes(mux, rep)
	CreditCardRoutes(mux, rep)
	PaymentTypeRoutes(mux, rep)
	PurchaseTypeRoutes(mux, rep)
	PurchaseRoutes(mux, rep)
	InstallmentRoutes(mux, rep)

	slog.Info(fmt.Sprintf("Server running on port %s - env: %s", config.ServerPort(), config.Env()))
	http.ListenAndServe(fmt.Sprintf(":%s", config.ServerPort()), c.Handler(mux))
}
