package repository

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/me/finance/pkg/entity"
)

type RepositoryCreditCard interface {
	Create(cc entity.CreditCard) error
	Update(cc entity.CreditCard) error
	Delete(id uuid.UUID) error
	FindByID(id uuid.UUID) (entity.CreditCard, error)
	FindAll() ([]entity.CreditCard, error)
}

type repositoryCreditCard struct {
	db *sql.DB
}

func NewRepositoryCreditCard(db *sql.DB) *repositoryCreditCard {
	return &repositoryCreditCard{db}
}

func (r repositoryCreditCard) Create(cc entity.CreditCard) error {
	query := `INSERT INTO credit_card (id, owner, final_card_num, invoice_closing_day) VALUES ($1, $2, $3, $4)`
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("error trying prepare statment: %v", err)
	}

	id, err := uuid.NewUUID()
	if err != nil {
		return fmt.Errorf("error trying create uuid: %v", err)
	}

	if _, err = stmt.Exec(id, cc.Owner, cc.FinalCardNum, cc.InvoiceClosingDay); err != nil {
		return fmt.Errorf("error trying insert credit card: %v", err)
	}

	if err := stmt.Close(); err != nil {
		return fmt.Errorf("error trying close statment: %v", err)
	}

	return nil
}

func (r repositoryCreditCard) Update(cc entity.CreditCard) error {
	query := `UPDATE credit_card 
				SET owner = $1,
					final_card_num = $2,
					invoice_closing_day = $3
				WHERE id = $4`
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("error trying prepare statment: %v", err)
	}

	if _, err = stmt.Exec(cc.Owner, cc.FinalCardNum, cc.InvoiceClosingDay, cc.ID); err != nil && err != sql.ErrNoRows {
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

func (r repositoryCreditCard) Delete(id uuid.UUID) error {
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

func (r repositoryCreditCard) FindByID(id uuid.UUID) (entity.CreditCard, error) {
	query := "SELECT id, owner, final_card_num, invoice_closing_day FROM credit_card WHERE id = $1"

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return entity.CreditCard{}, fmt.Errorf("error trying prepare statment: %v", err)
	}

	var cc entity.CreditCard
	if err = stmt.QueryRow(id).Scan(&cc.ID, &cc.Owner, &cc.FinalCardNum, &cc.InvoiceClosingDay); err != nil  && err != sql.ErrNoRows{
		return entity.CreditCard{}, fmt.Errorf("error trying find credit card: %v", err)
	}

	if err != nil && err == sql.ErrNoRows {
		return entity.CreditCard{}, fmt.Errorf("does not exist this id")
	}

	if err := stmt.Close(); err != nil {
		return entity.CreditCard{}, fmt.Errorf("error trying close statment: %v", err)
	}

	return cc, nil
}

func (r repositoryCreditCard) FindAll() ([]entity.CreditCard, error) {
	query := "SELECT id, owner, final_card_num, invoice_closing_day FROM credit_card order by owner"

	rows, err := r.db.Query(query)
	if err != nil {
		return []entity.CreditCard{}, fmt.Errorf("error trying find all credit cards: %v", err)
	}

	var creditCards []entity.CreditCard

	for rows.Next() {
		var cc entity.CreditCard
		if err = rows.Scan(&cc.ID, &cc.Owner, &cc.FinalCardNum, &cc.InvoiceClosingDay); err != nil && err != sql.ErrNoRows {
			return []entity.CreditCard{}, fmt.Errorf("error trying scan credit card: %v", err)
		}

		if err != nil && err == sql.ErrNoRows {
			return []entity.CreditCard{}, fmt.Errorf("does not exist credit card")
		}

		creditCards = append(creditCards, cc)
	}

	if rows.Close(); err != nil {
		return []entity.CreditCard{}, fmt.Errorf("error trying close rows: %v", err)
	}

	return creditCards, nil
}
