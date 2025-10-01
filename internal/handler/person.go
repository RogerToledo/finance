package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/me/finance/internal/models"
	"github.com/me/finance/internal/service"
)

type PersonHandler interface {
	RegisterRoutes(mux *http.ServeMux)
	CreatePerson(w http.ResponseWriter, r *http.Request)
	UpdatePerson(w http.ResponseWriter, r *http.Request)
	DeletePerson(w http.ResponseWriter, r *http.Request)
	FindPersonByID(w http.ResponseWriter, r *http.Request)
	FindAllPersons(w http.ResponseWriter, r *http.Request)
}

type personHandler struct {
	service service.PersonService
}

func NewPersonHandler(svc service.PersonService) PersonHandler {
	return &personHandler{service: svc}
}

func (h *personHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /person", func(w http.ResponseWriter, r *http.Request) {
		h.CreatePerson(w, r)
	})

	mux.HandleFunc("PUT /person", func(w http.ResponseWriter, r *http.Request) {
		h.UpdatePerson(w, r)
	})

	mux.HandleFunc("DELETE /person/{id}", func(w http.ResponseWriter, r *http.Request) {
		h.DeletePerson(w, r)
	})

	mux.HandleFunc("GET /person/{id}", func(w http.ResponseWriter, r *http.Request) {
		h.FindPersonByID(w, r)
	})

	mux.HandleFunc("GET /persons", func(w http.ResponseWriter, r *http.Request) {
		h.FindAllPersons(w, r)
	})
}

func (h *personHandler) CreatePerson(w http.ResponseWriter, r *http.Request) {
	var person models.Person

	err := json.NewDecoder(r.Body).Decode(&person)
	if err != nil {
		slog.Error(fmt.Sprintf("Error decoding person: %v", err))
		http.Error(w, fmt.Sprintf("Error decoding person: %v", err), http.StatusBadRequest)
		return
	}

	if err := person.Validate(true); err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.CreatePerson(person); err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	HTTPResponse(w, "Person was created with success!", http.StatusCreated)
}

func (h *personHandler) UpdatePerson(w http.ResponseWriter, r *http.Request) {
	var person models.Person

	err := json.NewDecoder(r.Body).Decode(&person)
	if err != nil {
		slog.Error(fmt.Sprintf("Error decoding person: %v", err))
		http.Error(w, fmt.Sprintf("Error decoding person: %v", err), http.StatusBadRequest)
		return
	}

	if err := person.Validate(false); err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.UpdatePerson(person); err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	HTTPResponse(w, "Person was updated with success!", http.StatusOK)
}

func (h *personHandler) DeletePerson(w http.ResponseWriter, r *http.Request) {
	idRequest := r.PathValue("id")

	id, err := models.ValidateID(idRequest)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.DeletePerson(id); err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	HTTPResponse(w, "Person was deleted with success!", http.StatusOK)
}

func (h *personHandler) FindPersonByID(w http.ResponseWriter, r *http.Request) {
	idRequest := r.PathValue("id")

	id, err := models.ValidateID(idRequest)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	person, err := h.service.FindPersonByID(id)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	HTTPResponse(w, person, http.StatusOK)
}

func (h *personHandler) FindAllPersons(w http.ResponseWriter, r *http.Request) {
	persons, err := h.service.FindAllPersons()
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	HTTPResponse(w, persons, http.StatusOK)
}
