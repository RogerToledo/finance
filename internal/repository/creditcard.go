package repository

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/me/finance/internal/models"
)

type CreditCardRepository interface {
	Create(cc models.CreditCard) error
	Update(cc models.CreditCard) error
	Delete(id uuid.UUID) error
	FindByID(id uuid.UUID) (models.CreditCard, error)
	FindAll() ([]models.CreditCard, error)
}

type creditCardRepository struct {
	db *sql.DB
}

func NewRepositoryCreditCard(db *sql.DB) *creditCardRepository {
	return &creditCardRepository{db}
}

func (r creditCardRepository) Create(cc models.CreditCard) error {
	query := `INSERT INTO credit_card (id, owner, final_card_num, type, invoice_closing_day) VALUES ($1, $2, $3, $4, $5)`
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("error trying prepare statment: %v", err)
	}

	id, err := uuid.NewUUID()
	if err != nil {
		return fmt.Errorf("error trying create uuid: %v", err)
	}

	if _, err = stmt.Exec(id, cc.Owner, cc.FinalCardNum, cc.Type, cc.InvoiceClosingDay); err != nil {
		return fmt.Errorf("error trying insert credit card: %v", err)
	}

	if err := stmt.Close(); err != nil {
		return fmt.Errorf("error trying close statment: %v", err)
	}

	return nil
}

func (r creditCardRepository) Update(cc models.CreditCard) error {
	query := `UPDATE credit_card 
				SET owner = $1,
					final_card_num = $2,
					type = $3,
					invoice_closing_day = $4
				WHERE id = $5`
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("error trying prepare statment: %v", err)
	}

	if _, err = stmt.Exec(
		cc.Owner,
		cc.FinalCardNum,
		cc.Type,
		cc.InvoiceClosingDay,
		cc.ID); err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error trying update credit card: %v", err)
	}

	if err != nil && err == sql.ErrNoRows {
		return fmt.Errorf("does not exist credit card with this id")
	}

	if err := stmt.Close(); err != nil {
		return fmt.Errorf("error trying close statment: %v", err)
	}

	return nil
}

func (r creditCardRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM credit_card WHERE id = $1`

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("error trying prepare statment: %v", err)
	}

	_, err = stmt.Exec(id)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error trying delete credit card: %v", err)
	}

	if err != nil && err == sql.ErrNoRows {
		return fmt.Errorf("does not exist credit card with this id")
	}

	if err := stmt.Close(); err != nil {
		return fmt.Errorf("error trying close statment: %v", err)
	}

	return nil
}

func (r creditCardRepository) FindByID(id uuid.UUID) (models.CreditCard, error) {
	query := "SELECT id, owner, final_card_num, type, invoice_closing_day FROM credit_card WHERE id = $1"

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return models.CreditCard{}, fmt.Errorf("error trying prepare statment: %v", err)
	}

	var cc models.CreditCard
	if err = stmt.QueryRow(id).Scan(&cc.ID, &cc.Owner, &cc.FinalCardNum, &cc.Type, &cc.InvoiceClosingDay); err != nil && err != sql.ErrNoRows {
		return models.CreditCard{}, fmt.Errorf("error trying find credit card: %v", err)
	}

	if err != nil && err == sql.ErrNoRows {
		return models.CreditCard{}, fmt.Errorf("does not exist this id")
	}

	if err := stmt.Close(); err != nil {
		return models.CreditCard{}, fmt.Errorf("error trying close statment: %v", err)
	}

	return cc, nil
}

func (r creditCardRepository) FindAll() ([]models.CreditCard, error) {
	query := `SELECT 
				id, 
				owner, 
				final_card_num, 
				CASE
					WHEN type = 'F' THEN 'Físico'
					WHEN type = 'V' THEN 'Virtual'
					WHEN type = 'VT' THEN 'Virtual Temporário'
				END AS type,
				invoice_closing_day 
			FROM credit_card 
			ORDER BY owner`

	rows, err := r.db.Query(query)
	if err != nil {
		return []models.CreditCard{}, fmt.Errorf("error trying find all credit cards: %v", err)
	}

	var creditCards []models.CreditCard

	for rows.Next() {
		var cc models.CreditCard
		if err = rows.Scan(&cc.ID, &cc.Owner, &cc.FinalCardNum, &cc.Type, &cc.InvoiceClosingDay); err != nil && err != sql.ErrNoRows {
			return []models.CreditCard{}, fmt.Errorf("error trying scan credit card: %v", err)
		}

		if err != nil && err == sql.ErrNoRows {
			return []models.CreditCard{}, fmt.Errorf("does not exist credit card")
		}

		creditCards = append(creditCards, cc)
	}

	if rows.Close(); err != nil {
		return []models.CreditCard{}, fmt.Errorf("error trying close rows: %v", err)
	}

	return creditCards, nil
}
