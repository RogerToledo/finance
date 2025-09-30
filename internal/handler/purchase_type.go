package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/me/finance/internal/models"
	"github.com/me/finance/internal/service"
	"github.com/sagikazarmark/slog-shim"
)

type PurchaseTypeHandler interface {
	RegisterRoutes(mux *http.ServeMux)
	CreatePurchaseType(w http.ResponseWriter, r *http.Request)
	UpdatePurchaseType(w http.ResponseWriter, r *http.Request)
	DeletePurchaseType(w http.ResponseWriter, r *http.Request)
	FindPurchaseTypeByID(w http.ResponseWriter, r *http.Request)
	FindAllPurchaseTypes(w http.ResponseWriter, r *http.Request)
}

type purchaseTypeHandler struct {
	service service.PurchaseTypeUseCase
}

func NewPurchaseTypeHandler(s service.PurchaseTypeUseCase) PurchaseTypeHandler {
	return &purchaseTypeHandler{service: s}
}

func (h *purchaseTypeHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /purchaseType", func(w http.ResponseWriter, r *http.Request) {
		h.CreatePurchaseType(w, r)
	})

	mux.HandleFunc("PUT /purchaseType", func(w http.ResponseWriter, r *http.Request) {
		h.UpdatePurchaseType(w, r)
	})

	mux.HandleFunc("DELETE /purchaseType/{id}", func(w http.ResponseWriter, r *http.Request) {
		h.DeletePurchaseType(w, r)
	})

	mux.HandleFunc("GET /purchaseType/{id}", func(w http.ResponseWriter, r *http.Request) {
		h.FindPurchaseTypeByID(w, r)
	})

	mux.HandleFunc("GET /purchaseTypes", func(w http.ResponseWriter, r *http.Request) {
		h.FindAllPurchaseTypes(w, r)
	})
}

func (pt *purchaseTypeHandler) CreatePurchaseType(w http.ResponseWriter, r *http.Request) {
	var purchaseType models.PurchaseType

	err := json.NewDecoder(r.Body).Decode(&purchaseType)
	if err != nil {
		slog.Error(fmt.Sprintf("Error decoding Purchase Type: %v", err))
		http.Error(w, fmt.Sprintf("Error decoding Purchase Type: %v", err), http.StatusBadRequest)
		return
	}

	if err := purchaseType.Validate(true); err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := pt.service.CreatePurchaseType(purchaseType); err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	HTTPResponse(w, "Purchase Type was created with success!", http.StatusCreated)
}

func (pt *purchaseTypeHandler) UpdatePurchaseType(w http.ResponseWriter, r *http.Request) {
	var purchaseType models.PurchaseType

	err := json.NewDecoder(r.Body).Decode(&purchaseType)
	if err != nil {
		slog.Error(fmt.Sprintf("Error decoding Purchase Type: %v", err))
		http.Error(w, fmt.Sprintf("Error decoding Purchase Type: %v", err), http.StatusBadRequest)
		return
	}

	if err := purchaseType.Validate(false); err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := pt.service.UpdatePurchaseType(purchaseType); err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	HTTPResponse(w, "Purchase Type was updated with success", http.StatusOK)
}

func (pt *purchaseTypeHandler) DeletePurchaseType(w http.ResponseWriter, r *http.Request) {
	idRequest := r.PathValue("id")

	id, err := models.ValidateID(idRequest)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = pt.service.DeletePurchaseType(id)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	HTTPResponse(w, "Purchase Type was deleted with success!", http.StatusOK)
}

func (pt *purchaseTypeHandler) FindPurchaseTypeByID(w http.ResponseWriter, r *http.Request) {
	idRequest := r.PathValue("id")

	id, err := models.ValidateID(idRequest)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	purchaseType, err := pt.service.FindPurchaseTypeByID(id)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	HTTPResponse(w, purchaseType, http.StatusOK)
}

func (pt *purchaseTypeHandler) FindAllPurchaseTypes(w http.ResponseWriter, r *http.Request) {
	purchaseTypes, err := pt.service.FindAllPurchaseTypes()
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	HTTPResponse(w, purchaseTypes, http.StatusOK)
}
