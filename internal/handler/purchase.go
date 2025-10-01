package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/me/finance/internal/models"
	"github.com/me/finance/internal/service"
	"github.com/sagikazarmark/slog-shim"
)

type PurchaseHandler interface {
	RegisterRoutes(mux *http.ServeMux)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	FindByID(w http.ResponseWriter, r *http.Request)
	FindByDate(w http.ResponseWriter, r *http.Request)
	FindByMonth(w http.ResponseWriter, r *http.Request)
	FindByPerson(w http.ResponseWriter, r *http.Request)
	FindAll(w http.ResponseWriter, r *http.Request)
}

type purchaseHandler struct {
	service service.PurchaseService
}

func NewPurchaseHandler(svc service.PurchaseService) PurchaseHandler {
	return &purchaseHandler{service: svc}
}

func (h *purchaseHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /purchase", func(w http.ResponseWriter, r *http.Request) {
		h.Create(w, r)
	})

	mux.HandleFunc("PUT /purchase", func(w http.ResponseWriter, r *http.Request) {
		h.Update(w, r)
	})

	mux.HandleFunc("DELETE /purchase/{id}", func(w http.ResponseWriter, r *http.Request) {
		h.Delete(w, r)
	})

	mux.HandleFunc("GET /purchase/{id}", func(w http.ResponseWriter, r *http.Request) {
		h.FindByID(w, r)
	})

	mux.HandleFunc("GET /purchase/date/{date}", func(w http.ResponseWriter, r *http.Request) {
		h.FindByDate(w, r)
	})

	mux.HandleFunc("GET /purchase/month/{date}", func(w http.ResponseWriter, r *http.Request) {
		h.FindByMonth(w, r)
	})

	mux.HandleFunc("GET /purchase/person/{id}", func(w http.ResponseWriter, r *http.Request) {
		h.FindByPerson(w, r)
	})

	mux.HandleFunc("GET /purchases", func(w http.ResponseWriter, r *http.Request) {
		h.FindAll(w, r)
	})
}

func (p *purchaseHandler) Create(w http.ResponseWriter, r *http.Request) {
	var (
		purchase        models.Purchase
		purchaseRequest models.PurchaseRequest
	)

	err := json.NewDecoder(r.Body).Decode(&purchaseRequest)
	if err != nil {
		slog.Error(fmt.Sprintf("Error decoding Purchase: %v", err))
		http.Error(w, fmt.Sprintf("Error decoding Purchase: %v", err), http.StatusBadRequest)
		return
	}

	purchase, err = purchaseRequest.ToEntity()
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := purchase.Validate(); err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := p.service.CreatePurchase(purchase); err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	HTTPResponse(w, "Purchase was created with success!", http.StatusCreated)
}

func (p *purchaseHandler) Update(w http.ResponseWriter, r *http.Request) {
	var (
		purchase        models.Purchase
		purchaseRequest models.PurchaseRequest
	)

	err := json.NewDecoder(r.Body).Decode(&purchaseRequest)
	if err != nil {
		slog.Error(fmt.Sprintf("Error decoding Purchase: %v", err))
		http.Error(w, fmt.Sprintf("Error decoding Purchase: %v", err), http.StatusBadRequest)
		return
	}

	purchase, err = purchaseRequest.ToEntity()
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if err := purchase.Validate(); err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := p.service.UpdatePurchase(purchase); err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	HTTPResponse(w, "Purchase was updated with success!", http.StatusOK)
}

func (p *purchaseHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idRequest := r.PathValue("id")

	id, err := models.ValidateID(idRequest)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = p.service.DeletePurchase(id)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	HTTPResponse(w, "Purchase was deleted with success!", http.StatusOK)
}

func (p *purchaseHandler) FindByID(w http.ResponseWriter, r *http.Request) {
	idRequest := r.PathValue("id")

	id, err := models.ValidateID(idRequest)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	purchase, err := p.service.FindPurchaseByID(id)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	HTTPResponse(w, purchase, http.StatusOK)
}

func (p *purchaseHandler) FindByDate(w http.ResponseWriter, r *http.Request) {
	date := r.PathValue("date")

	if err := models.ValidateDate(date); err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	purchases, err := p.service.FindPurchaseByDate(date)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	HTTPResponse(w, purchases, http.StatusOK)
}

func (p *purchaseHandler) FindByMonth(w http.ResponseWriter, r *http.Request) {
	date := r.PathValue("date")

	if err := models.ValidateYearMonth(date); err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	purchases, err := p.service.FindPurchaseByMonth(date)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	HTTPResponse(w, purchases, http.StatusOK)
}

func (p *purchaseHandler) FindByPerson(w http.ResponseWriter, r *http.Request) {
	idRequest := r.PathValue("id")

	id, err := models.ValidateID(idRequest)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	purchases, err := p.service.FindPurchaseByPerson(id)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	HTTPResponse(w, purchases, http.StatusOK)
}

func (p *purchaseHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	purchases, err := p.service.FindAllPurchases()
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	HTTPResponse(w, purchases, http.StatusOK)
}
