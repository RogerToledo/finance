package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/me/finance/internal/models"
	"github.com/me/finance/internal/service"
	"github.com/sagikazarmark/slog-shim"
)

type PaymentTypeHandler interface {
	RegisterRoutes(mux *http.ServeMux)
	CreatePaymentType(w http.ResponseWriter, r *http.Request)
	UpdatePaymentType(w http.ResponseWriter, r *http.Request)
	DeletePaymentType(w http.ResponseWriter, r *http.Request)
	FindPaymentTypeByID(w http.ResponseWriter, r *http.Request)
	FindAllPaymentTypes(w http.ResponseWriter, r *http.Request)
}

type paymentTypeHandler struct {
	service service.PaymentTypeService
}

func NewPaymentTypeHandler(svc service.PaymentTypeService) PaymentTypeHandler {
	return &paymentTypeHandler{service: svc}
}

func (h *paymentTypeHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /paymentType", func(w http.ResponseWriter, r *http.Request) {
		h.CreatePaymentType(w, r)
	})
	
	mux.HandleFunc("PUT /paymentType", func(w http.ResponseWriter, r *http.Request) {
		h.UpdatePaymentType(w, r)
	})

	mux.HandleFunc("DELETE /paymentType/{id}", func(w http.ResponseWriter, r *http.Request) {
		h.DeletePaymentType(w, r)
	})

	mux.HandleFunc("GET /paymentType/{id}", func(w http.ResponseWriter, r *http.Request) {
		h.FindPaymentTypeByID(w, r)
	})

	mux.HandleFunc("GET /paymentTypes", func(w http.ResponseWriter, r *http.Request) {
		h.FindAllPaymentTypes(w, r)
	})
}

func (pt *paymentTypeHandler) CreatePaymentType(w http.ResponseWriter, r *http.Request) {
	var paymentType models.PaymentType

	err := json.NewDecoder(r.Body).Decode(&paymentType)
	if err != nil {
		slog.Error(fmt.Sprintf("Error decoding payment: %v", err))
		http.Error(w, fmt.Sprintf("Error decoding payment: %v", err), http.StatusBadRequest)
		return
	}

	if err := paymentType.Validate(true); err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := pt.service.CreatePaymentType(paymentType); err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	HTTPResponse(w, fmt.Sprintf("Payment Type was created with success!"), http.StatusCreated)
}

func (pt *paymentTypeHandler) UpdatePaymentType(w http.ResponseWriter, r *http.Request) {
	var paymentType models.PaymentType

	err := json.NewDecoder(r.Body).Decode(&paymentType)
	if err != nil {
		slog.Error(fmt.Sprintf("Error decoding Payment Type: %v", err))
		http.Error(w, fmt.Sprintf("Error decoding Payment Type: %v", err), http.StatusBadRequest)
		return
	}

	if err := paymentType.Validate(false); err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := pt.service.UpdatePaymentType(paymentType); err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	HTTPResponse(w, fmt.Sprint("Payment Type was updated with success!"), http.StatusOK)
}

func (pt *paymentTypeHandler) DeletePaymentType(w http.ResponseWriter, r *http.Request) {
	idRequest := r.PathValue("id")
	
	id, err := models.ValidateID(idRequest)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = pt.service.DeletePaymentType(id)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	HTTPResponse(w, fmt.Sprintf("Payment Type was deleted with success!"), http.StatusOK)
}

func (pt *paymentTypeHandler) FindPaymentTypeByID(w http.ResponseWriter, r *http.Request) {
	idRequest := r.PathValue("id")

	id, err := models.ValidateID(idRequest)
	if err != nil {
		fmt.Sprintf(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	paymentType, err := pt.service.FindPaymentTypeByID(id)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	HTTPResponse(w, paymentType, http.StatusOK)
}

func (pt *paymentTypeHandler) FindAllPaymentTypes(w http.ResponseWriter, r *http.Request) {
	paymentTypes, err := pt.service.FindAllPaymentTypes()
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	HTTPResponse(w, paymentTypes, http.StatusOK)
}