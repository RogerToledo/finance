package handler

import (
    "encoding/json"
    "fmt"
    "log/slog"
    "net/http"

    "github.com/me/finance/internal/models"
    "github.com/me/finance/internal/service"
)

// PersonHandler define a interface para os handlers de pessoa.
type PersonHandler interface {
	RegisterRoutes(mux *http.ServeMux)
    CreatePerson(w http.ResponseWriter, r *http.Request)
    UpdatePerson(w http.ResponseWriter, r *http.Request)
    DeletePerson(w http.ResponseWriter, r *http.Request)
    FindPersonByID(w http.ResponseWriter, r *http.Request)
    FindAllPersons(w http.ResponseWriter, r *http.Request)
}

type personHandler struct {
    service service.PersonService // O useCase agora é chamado de service
}

// NewPersonHandler cria uma nova instância do handler de pessoa.
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

    // Supondo que você tenha uma função helper para respostas HTTP
    // HTTPResponse(w, fmt.Sprint("Person was created with success!"), http.StatusCreated)
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"message": "Person was created with success!"})
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

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Person was updated with success!"})
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

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Person was deleted with success!"})
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

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(person)
}

func (h *personHandler) FindAllPersons(w http.ResponseWriter, r *http.Request) {
    persons, err := h.service.FindAllPersons()
    if err != nil {
        slog.Error(err.Error())
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(persons)
}