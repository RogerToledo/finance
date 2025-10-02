package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/me/finance/internal/models"
	"github.com/me/finance/internal/service"
	"github.com/sagikazarmark/slog-shim"
)

type CreditCardHandler interface {
	RegisterRoutes(mux *http.ServeMux)
	CreateCreditCard(w http.ResponseWriter, r *http.Request)
	UpdateCreditCard(w http.ResponseWriter, r *http.Request)
	DeleteCreditCard(w http.ResponseWriter, r *http.Request)
	FindCreditCardByID(w http.ResponseWriter, r *http.Request)
	FindAllCreditCards(w http.ResponseWriter, r *http.Request)
}

type creditCardHandler struct {
	service service.CreditCardService
}

func NewCreditCardHandler(svc service.CreditCardService) CreditCardHandler {
	return &creditCardHandler{service: svc}
}

func (h *creditCardHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /v1/creditCards", func(w http.ResponseWriter, r *http.Request) {
		h.CreateCreditCard(w, r)
	})

	mux.HandleFunc("PUT /v1/creditCards", func(w http.ResponseWriter, r *http.Request) {
		h.UpdateCreditCard(w, r)
	})

	mux.HandleFunc("DELETE /v1/creditCards/{id}", func(w http.ResponseWriter, r *http.Request) {
		h.DeleteCreditCard(w, r)
	})

	mux.HandleFunc("GET /v1/creditCards/{id}", func(w http.ResponseWriter, r *http.Request) {
		h.FindCreditCardByID(w, r)
	})

	mux.HandleFunc("GET /v1/creditCards", func(w http.ResponseWriter, r *http.Request) {
		h.FindAllCreditCards(w, r)
	})
}

func (c *creditCardHandler) CreateCreditCard(w http.ResponseWriter, r *http.Request) {
	var creditCard models.CreditCard

	if err := json.NewDecoder(r.Body).Decode(&creditCard); err != nil {
		slog.Error(fmt.Sprintf("Error decoding credit card: %v", err))
		http.Error(w, fmt.Sprintf("Error decoding credit card: %v", err), http.StatusBadRequest)
		return
	}

	if err := creditCard.Validate(true); err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := c.service.CreateCreditCard(creditCard); err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	HTTPResponse(w, fmt.Sprintf("Credit card created to %s", creditCard.Owner), http.StatusOK)
}

func (c *creditCardHandler) UpdateCreditCard(w http.ResponseWriter, r *http.Request) {
	var creditCard models.CreditCard

	if err := json.NewDecoder(r.Body).Decode(&creditCard); err != nil {
		slog.Error(fmt.Sprintf("Error decoding credit card: %v", err))
		http.Error(w, fmt.Sprintf("Error decoding credit card: %v", err), http.StatusBadRequest)
		return
	}

	if err := creditCard.Validate(false); err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := c.service.UpdateCreditCard(creditCard); err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	HTTPResponse(w, "Credit card was updated with success!", http.StatusOK)
}

func (c *creditCardHandler) DeleteCreditCard(w http.ResponseWriter, r *http.Request) {
	idRequest := r.PathValue("id")

	id, err := models.ValidateID(idRequest)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := c.service.DeleteCreditCard(id); err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	HTTPResponse(w, fmt.Sprint("Credit card deleted with sucess!"), http.StatusOK)
}

func (c *creditCardHandler) FindCreditCardByID(w http.ResponseWriter, r *http.Request) {
	idRequest := r.PathValue("id")

	id, err := models.ValidateID(idRequest)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	creditCard, err := c.service.FindCreditCardByID(id)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	HTTPResponse(w, creditCard, http.StatusOK)
}

func (c *creditCardHandler) FindAllCreditCards(w http.ResponseWriter, r *http.Request) {
	creditCard, err := c.service.FindAllCreditCards()
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	HTTPResponse(w, creditCard, http.StatusOK)
}
