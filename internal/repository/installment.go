package repository

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/me/finance/internal/models"
)

type InstallmentRepository interface {
	Create(tx *sql.Tx, installment models.Installment) error
	Update(id uuid.UUID) error
	Delete(id uuid.UUID) error
	FindByPurchaseID(id uuid.UUID) ([]models.Installment, error)
	FindByMonth(month string) ([]models.Installment, error)
	FindByNotPaid() ([]models.Installment, error)
}

type installmentRepository struct {
	db *sql.DB
}

func NewInstallmentRepository(db *sql.DB) *installmentRepository {
	return &installmentRepository{db}
}

func (r *installmentRepository) Create(tx *sql.Tx, installment models.Installment) error {
	sql := `INSERT INTO installment (id, description, number, value, month, paid, purchase_id) 
			VALUES 
			($1, $2, $3, $4, $5, $6, $7)`

	installment.ID = uuid.New()

	_, err := tx.Exec(sql,
		installment.ID,
		installment.Description,
		installment.Number,
		installment.Value,
		installment.Month,
		installment.Paid,
		installment.PurchaseID,
	)
	if err != nil {
		return fmt.Errorf("error executing statement: %v", err)
	}

	return nil
}

func (r *installmentRepository) Update(id uuid.UUID) error {
	sql := `UPDATE installment SET paid = true WHERE id = $1`

	stmt, err := r.db.Prepare(sql)
	if err != nil {
		return fmt.Errorf("error preparing statement: %v", err)
	}

	_, err = stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("error executing statement: %v", err)
	}

	return nil
}

func (r *installmentRepository) Delete(id uuid.UUID) error {
	sql := `DELETE FROM installment WHERE purchase_id = $1`

	stmt, err := r.db.Prepare(sql)
	if err != nil {
		return fmt.Errorf("error preparing statement: %v", err)
	}

	_, err = stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("error executing statement: %v", err)
	}

	return nil
}

func (r *installmentRepository) FindByPurchaseID(id uuid.UUID) ([]models.Installment, error) {
	sql := `SELECT id, description, number, value, month, paid, purchase_id
			 FROM installment 
			 WHERE purchase_id = $1`

	stmt, err := r.db.Prepare(sql)
	if err != nil {
		return nil, fmt.Errorf("error preparing statement: %v", err)
	}

	rows, err := stmt.Query(id)
	if err != nil {
		return nil, fmt.Errorf("error executing statement: %v", err)
	}

	var installments []models.Installment
	for rows.Next() {
		var installment models.Installment
		err = rows.Scan(
			&installment.ID,
			&installment.Description,
			&installment.Number,
			&installment.Value,
			&installment.Month,
			&installment.Paid,
			&installment.PurchaseID,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning rows: %v", err)
		}

		installments = append(installments, installment)
	}

	return installments, nil
}

func (r *installmentRepository) FindByMonth(month string) ([]models.Installment, error) {
	sql := `SELECT id, description, number, value, month, paid, purchase_id
			 FROM installment 
			 WHERE to_char(month, 'YYYY-MM') = $1`

	stmt, err := r.db.Prepare(sql)
	if err != nil {
		return nil, fmt.Errorf("error preparing statement: %v", err)
	}

	rows, err := stmt.Query(month)
	if err != nil {
		return nil, fmt.Errorf("error executing statement: %v", err)
	}

	var installments []models.Installment
	for rows.Next() {
		var installment models.Installment
		err = rows.Scan(
			&installment.ID,
			&installment.Description,
			&installment.Number,
			&installment.Value,
			&installment.Month,
			&installment.Paid,
			&installment.PurchaseID,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning rows: %v", err)
		}

		installments = append(installments, installment)
	}

	return installments, nil
}

func (r *installmentRepository) FindByNotPaid() ([]models.Installment, error) {
	sql := `SELECT id, description, number, value, month, paid, purchase_id 
			FROM installment 
			WHERE paid = false`

	stmt, err := r.db.Prepare(sql)
	if err != nil {
		return nil, fmt.Errorf("error preparing statement: %v", err)
	}

	rows, err := stmt.Query()
	if err != nil {
		return nil, fmt.Errorf("error executing statement: %v", err)
	}

	var installments []models.Installment
	for rows.Next() {
		var installment models.Installment
		err = rows.Scan(
			&installment.ID,
			&installment.Description,
			&installment.Number,
			&installment.Value,
			&installment.Month,
			&installment.Paid,
			&installment.PurchaseID,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning rows: %v", err)
		}

		installments = append(installments, installment)
	}

	return installments, nil
}
