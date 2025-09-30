package repository

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/me/finance/internal/models"
)

type RepositoryPurchaseType interface {
	Create(p models.PurchaseType) error
	Update(pt models.PurchaseType) error
	Delete(id uuid.UUID) error
	FindByID(id uuid.UUID) (models.PurchaseType, error)
	FindAll() ([]models.PurchaseType, error)
}

type repositoryPurchaseType struct {
	db *sql.DB
}

func NewRepositoryPurchaseType(db *sql.DB) *repositoryPurchaseType {
	return &repositoryPurchaseType{db}
}

func (r repositoryPurchaseType) Create(p models.PurchaseType) error {
	query := `INSERT INTO purchase_type (id, name) VALUES ($1, $2)`
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("error trying prepare statment: %v", err)
	}

	id, err := uuid.NewUUID()
	if err != nil {
		return fmt.Errorf("error trying create uuid: %v", err)
	}
	if _, err = stmt.Exec(id, p.Name); err != nil {
		return fmt.Errorf("error trying insert purchase type type: %v", err)
	}

	if err := stmt.Close(); err != nil {
		return fmt.Errorf("error trying close statment: %v", err)
	}

	return nil
}

func (r repositoryPurchaseType) Update(pt models.PurchaseType) error {
	query := `UPDATE purchase_type SET name = $1 WHERE id = $2`
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("error trying prepare statment: %v", err)
	}

	if _, err = stmt.Exec(pt.Name, pt.ID); err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error trying update purchase type: %v", err)
	}

	if err != nil && err == sql.ErrNoRows {
		return fmt.Errorf("does not exist purchase type with this id")
	}

	if err := stmt.Close(); err != nil {
		return fmt.Errorf("error trying close statment: %v", err)
	}

	return nil
}

func (r repositoryPurchaseType) Delete(id uuid.UUID) error {
	query := `DELETE FROM purchase_type WHERE id = $1`

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("error trying prepare statment: %v", err)
	}

	_, err = stmt.Exec(id)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error trying delete purchase type: %v", err)
	}

	if err != nil && err == sql.ErrNoRows {
		return fmt.Errorf("does not exist purchase type with this id")
	}

	if err := stmt.Close(); err != nil {
		return fmt.Errorf("error trying close statment: %v", err)
	}

	return nil
}

func (r repositoryPurchaseType) FindByID(id uuid.UUID) (models.PurchaseType, error) {
	query := "SELECT id, name FROM purchase_type WHERE id = $1"

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return models.PurchaseType{}, fmt.Errorf("error trying prepare statment: %v", err)
	}

	var pt models.PurchaseType
	if err = stmt.QueryRow(id).Scan(&pt.ID, &pt.Name); err != nil && err != sql.ErrNoRows {
		return models.PurchaseType{}, fmt.Errorf("error trying find purchase type: %v", err)
	}

	if err != nil && err == sql.ErrNoRows {
		return models.PurchaseType{}, fmt.Errorf("does not exist purchase type with this id")
	}

	if err := stmt.Close(); err != nil {
		return models.PurchaseType{}, fmt.Errorf("error trying close statment: %v", err)
	}

	return pt, nil
}

func (r repositoryPurchaseType) FindAll() ([]models.PurchaseType, error) {
	query := "SELECT id, name FROM purchase_type ORDER BY name"

	rows, err := r.db.Query(query)
	if err != nil {
		slog.Error("error trying find all purchase type: %v", err)
		return []models.PurchaseType{}, fmt.Errorf("error trying find all purchase type: %v", err)
	}

	var purchases []models.PurchaseType

	for rows.Next() {
		var pt models.PurchaseType
		if err = rows.Scan(&pt.ID, &pt.Name); err != nil && err != sql.ErrNoRows {
			return []models.PurchaseType{}, fmt.Errorf("error trying scan purchase type: %v", err)
		}

		if err != nil && err == sql.ErrNoRows {
			return []models.PurchaseType{}, fmt.Errorf("does not exist purchase type with this name")
		}

		purchases = append(purchases, pt)
	}

	if err := rows.Close(); err != nil {
		return []models.PurchaseType{}, fmt.Errorf("error trying close rows: %v", err)
	}

	return purchases, nil
}
