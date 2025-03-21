package router

import (
	"net/http"

	"github.com/me/finance/pkg/repository"

	"github.com/me/finance/pkg/controller"
)

func PurchaseTypeRoutes(mux *http.ServeMux, rep *repository.Repository) {
	c := controller.NewController(rep)
	
	mux.HandleFunc("POST /purchaseType", func(w http.ResponseWriter, r *http.Request) {
		c.PurchaseType.CreatePurchaseType(w, r)
	})
	
	mux.HandleFunc("PUT /purchaseType", func(w http.ResponseWriter, r *http.Request) {
		c.PurchaseType.UpdatePurchaseType(w, r)
	})

	mux.HandleFunc("DELETE /purchaseType/{id}", func(w http.ResponseWriter, r *http.Request) {
		c.PurchaseType.DeletePurchaseType(w, r)
	})

	mux.HandleFunc("GET /purchaseType/{id}", func(w http.ResponseWriter, r *http.Request) {
		c.PurchaseType.FindPurchaseTypeByID(w, r)
	})

	mux.HandleFunc("GET /purchaseTypes", func(w http.ResponseWriter, r *http.Request) {
		c.PurchaseType.FindAllPurchaseTypes(w, r)
	})
}
