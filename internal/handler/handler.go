package handler

import (
	"fmt"
	"net/http"

	"github.com/me/finance/internal/config"
	"github.com/me/finance/internal/database"
	"github.com/me/finance/internal/repository"
	"github.com/me/finance/internal/service"
	"github.com/rs/cors"
	"github.com/sagikazarmark/slog-shim"
)

func InitializeHandler() {
	db, err := database.NewDB()
	if err != nil {
		slog.Error("error trying to connect to database")
		return
	}

	mux := http.NewServeMux()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // Allow your frontend origin
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false, // Important for cookies, authorization headers with CORS
		Debug:            false,  // Enable for debugging CORS issues
	})

	personRepo := repository.NewRepositoryPerson(db)
	personService := service.NewPersonService(personRepo)
	personHandler := NewPersonHandler(personService)
	personHandler.RegisterRoutes(mux)

	slog.Info(fmt.Sprintf("Server running on port %s - env: %s", config.ServerPort(), config.Env()))
	http.ListenAndServe(fmt.Sprintf(":%s", config.ServerPort()), c.Handler(mux))
}
