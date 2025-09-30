package repository

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/me/finance/internal/models"
)

type PersonRepository interface {
	Create(p models.Person) error
	Update(p models.Person) error
	Delete(id uuid.UUID) error
	FindByID(id uuid.UUID) (models.Person, error)
	FindAll() ([]models.Person, error)
}

type personRepository struct {
	db *sql.DB
}

func NewRepositoryPerson(db *sql.DB) *personRepository {
	return &personRepository{db}
}

func (r personRepository) Create(p models.Person) error {
	query := `INSERT INTO person (id, name) VALUES ($1, $2)`
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("error trying prepare statment: %v", err)
	}

	id, err := uuid.NewUUID()
	if err != nil {
		return fmt.Errorf("error trying create uuid: %v", err)
	}

	if _, err = stmt.Exec(id, p.Name); err != nil {
		return fmt.Errorf("error trying insert person: %v", err)
	}

	if err := stmt.Close(); err != nil {
		return fmt.Errorf("error trying close stmt: %v", err)
	}

	return nil
}

func (r personRepository) Update(p models.Person) error {
	query := `UPDATE person SET name = $1 WHERE id = $2`
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("error trying prepare statment: %v", err)
	}

	if _, err = stmt.Exec(p.Name, p.ID); err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error trying update person: %v", err)
	}

	if err != nil && err == sql.ErrNoRows {
		return fmt.Errorf("does not exist person with this id")
	}

	if err := stmt.Close(); err != nil {
		return fmt.Errorf("error trying close stmt: %v", err)
	}

	return nil
}

func (r personRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM person WHERE id = $1`

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("error trying prepare statment: %v", err)
	}

	_, err = stmt.Exec(id)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error trying delete person: %v", err)
	}

	if err != nil && err == sql.ErrNoRows {
		return fmt.Errorf("does not exist person with this id")
	}

	if err := stmt.Close(); err != nil {
		return fmt.Errorf("error trying close stmt: %v", err)
	}

	return nil
}

func (r personRepository) FindByID(id uuid.UUID) (models.Person, error) {
	query := "SELECT id, name FROM person WHERE id = $1"

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return models.Person{}, fmt.Errorf("error trying prepare statment: %v", err)
	}

	var p models.Person
	if err = stmt.QueryRow(id).Scan(&p.ID, &p.Name); err != nil && err != sql.ErrNoRows {
		return models.Person{}, fmt.Errorf("error trying find person: %v", err)
	}

	if err != nil && err == sql.ErrNoRows {
		return models.Person{}, fmt.Errorf("does not exist person with this id")
	}

	if err := stmt.Close(); err != nil {
		return models.Person{}, fmt.Errorf("error trying close stmt: %v", err)
	}

	return p, nil
}

func (r personRepository) FindAll() ([]models.Person, error) {
	query := "SELECT id, name FROM person ORDER BY name"

	rows, err := r.db.Query(query)
	if err != nil {
		return []models.Person{}, fmt.Errorf("error trying find all persons: %v", err)
	}

	var persons []models.Person

	for rows.Next() {
		var p models.Person
		if err = rows.Scan(&p.ID, &p.Name); err != nil && err != sql.ErrNoRows {
			return []models.Person{}, fmt.Errorf("error trying scan person: %v", err)
		}

		if err != nil && err == sql.ErrNoRows {
			return []models.Person{}, fmt.Errorf("does not exist person with this name")
		}

		persons = append(persons, p)
	}

	if err := rows.Close(); err != nil {
		return []models.Person{}, fmt.Errorf("error trying close rows: %v", err)
	}

	return persons, nil
}
