package main

import (
	"fmt"
	"net/http"

	"github.com/me/finance/internal/config"
	"github.com/me/finance/internal/config/logger"
	"github.com/me/finance/internal/database"
	"github.com/me/finance/internal/handler"
	"github.com/me/finance/internal/repository"
	"github.com/me/finance/internal/service"
	"github.com/rs/cors"
	"github.com/sagikazarmark/slog-shim"
)

func main() {
	logger.InitLogger()

	if err := config.Load(); err != nil {
		panic(err)
	}

	db, err := database.NewDB()
	if err != nil {
		slog.Error("error trying to connect to database")
		return
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // Allow your frontend origin
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false, // Important for cookies, authorization headers with CORS
		Debug:            false, // Enable for debugging CORS issues
	})

	personRepo := repository.NewRepositoryPerson(db)
	personService := service.NewPersonService(personRepo)
	personHandler := handler.NewPersonHandler(personService)
	personHandler.RegisterRoutes(mux)

	creditcardRepo := repository.NewRepositoryCreditCard(db)
	creditcardService := service.NewCreditCardService(creditcardRepo)
	creditcardHandler := handler.NewCreditCardHandler(creditcardService)
	creditcardHandler.RegisterRoutes(mux)

	paymentTypeRepo := repository.NewRepositoryPaymentType(db)
	paymentTypeService := service.NewPaymentTypeService(paymentTypeRepo)
	paymentTypeHandler := handler.NewPaymentTypeHandler(paymentTypeService)
	paymentTypeHandler.RegisterRoutes(mux)

	purchaseTypeRepo := repository.NewRepositoryPurchaseType(db)
	purchaseTypeService := service.NewPurchaseTypeService(purchaseTypeRepo)
	purchaseTypeHandler := handler.NewPurchaseTypeHandler(purchaseTypeService)
	purchaseTypeHandler.RegisterRoutes(mux)

	installmentRepo := repository.NewInstallmentRepository(db)
	installmentService := service.NewInstallmentService(installmentRepo, creditcardRepo)
	installmentHandler := handler.NewInstallmentHandler(installmentService)
	installmentHandler.RegisterRoutes(mux)

	purchaseRepo := repository.NewRepositoryPurchase(db)
	purchaseService := service.NewPurchaseService(purchaseRepo, installmentRepo, creditcardRepo)
	purchaseHandler := handler.NewPurchaseHandler(purchaseService)
	purchaseHandler.RegisterRoutes(mux)

	slog.Info(fmt.Sprintf("Server running on port %s - env: %s", config.ServerPort(), config.Env()))
	http.ListenAndServe(fmt.Sprintf(":%s", config.ServerPort()), c.Handler(mux))
}
