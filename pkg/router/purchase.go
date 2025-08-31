package router

import (
	"net/http"

	"github.com/me/finance/pkg/repository"

	"github.com/me/finance/pkg/controller"
)

func PurchaseRoutes(mux *http.ServeMux, rep *repository.Repository) {
	c := controller.NewController(rep)
	
	mux.HandleFunc("POST /purchase", func(w http.ResponseWriter, r *http.Request) {
		c.Purchase.Create(w, r)
	})
	
	mux.HandleFunc("PUT /purchase", func(w http.ResponseWriter, r *http.Request) {
		c.Purchase.Update(w, r)
	})

	mux.HandleFunc("DELETE /purchase/{id}", func(w http.ResponseWriter, r *http.Request) {
		c.Purchase.Delete(w, r)
	})

	mux.HandleFunc("GET /purchase/{id}", func(w http.ResponseWriter, r *http.Request) {
		c.Purchase.FindByID(w, r)
	})

	mux.HandleFunc("GET /purchase/date/{date}", func(w http.ResponseWriter, r *http.Request) {
		c.Purchase.FindByDate(w, r)
	})

	mux.HandleFunc("GET /purchase/month/{date}", func(w http.ResponseWriter, r *http.Request) {
		c.Purchase.FindByMonth(w, r)
	})

	mux.HandleFunc("GET /purchase/person/{id}", func(w http.ResponseWriter, r *http.Request) {
		c.Purchase.FindByPerson(w, r)
	})

	mux.HandleFunc("GET /purchases", func(w http.ResponseWriter, r *http.Request) {
		c.Purchase.FindAll(w, r)
	})
}
