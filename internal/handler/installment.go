package handler

import (
	"log/slog"
	"net/http"

	"github.com/me/finance/internal/models"
	"github.com/me/finance/internal/service"
)

type InstallmentHandler interface {
	RegisterRoutes(mux *http.ServeMux)
	UpdateInstallment(w http.ResponseWriter, r *http.Request)
	FindInstallmentByPurchaseID(w http.ResponseWriter, r *http.Request)
	FindInstallmentByMonth(w http.ResponseWriter, r *http.Request)
	FindInstallmentByNotPaid(w http.ResponseWriter, r *http.Request)
}

type installmentHandler struct {
	service service.InstallmentService
}

func NewInstallmentHandler(svc service.InstallmentService) InstallmentHandler {
	return &installmentHandler{service: svc}
}

func (h *installmentHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("PUT /installment/{id}", func(w http.ResponseWriter, r *http.Request) {
		h.UpdateInstallment(w, r)
	})

	mux.HandleFunc("GET /installment/{id}", func(w http.ResponseWriter, r *http.Request) {
		h.FindInstallmentByPurchaseID(w, r)
	})

	mux.HandleFunc("GET /installment/month/{date}", func(w http.ResponseWriter, r *http.Request) {
		h.FindInstallmentByMonth(w, r)
	})

	mux.HandleFunc("GET /installment/notPaid", func(w http.ResponseWriter, r *http.Request) {
		h.FindInstallmentByNotPaid(w, r)
	})
}

func (i *installmentHandler) UpdateInstallment(w http.ResponseWriter, r *http.Request) {
	idRequest := r.PathValue("id")

	id, err := models.ValidateID(idRequest)
	if err != nil {
		slog.Error(err.Error())
		HTTPResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := i.service.UpdateInstalment(id); err != nil {
		slog.Error(err.Error())
		HTTPResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	HTTPResponse(w, "Installment was updated with success!", http.StatusOK)
}

func (i *installmentHandler) FindInstallmentByPurchaseID(w http.ResponseWriter, r *http.Request) {
	var response models.InstallmentResponse

	idRequest := r.PathValue("id")

	id, err := models.ValidateID(idRequest)
	if err != nil {
		slog.Error(err.Error())
		HTTPResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err = i.service.FindInstallmentByPurchaseID(id)
	if err != nil {
		slog.Error(err.Error())
		HTTPResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	HTTPResponse(w, response, http.StatusOK)
}

func (i *installmentHandler) FindInstallmentByMonth(w http.ResponseWriter, r *http.Request) {
	month := r.PathValue("date")

	if err := models.ValidateYearMonth(month); err != nil {
		slog.Error(err.Error())
		HTTPResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	installments, err := i.service.FindInstallmentByMonth(month)
	if err != nil {
		slog.Error(err.Error())
		HTTPResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	HTTPResponse(w, installments, http.StatusOK)

}

func (i *installmentHandler) FindInstallmentByNotPaid(w http.ResponseWriter, r *http.Request) {
	installments, err := i.service.FindInstallmentByNotPaid()
	if err != nil {
		slog.Error(err.Error())
		HTTPResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	HTTPResponse(w, installments, http.StatusOK)
}
