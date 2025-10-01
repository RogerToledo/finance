package repository

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/me/finance/internal/models"
)

type PurchaseRepository interface {
	BeginTransaction() (*sql.Tx, error)
	Commit(tx *sql.Tx) error
	Rollback(tx *sql.Tx) error
	Create(tx *sql.Tx, p models.Purchase) (uuid.UUID, error)
	Update(tx *sql.Tx, p models.Purchase) error
	Delete(tx *sql.Tx, id uuid.UUID) error
	FindByID(id uuid.UUID) (models.PurchaseResponse, error)
	FindByDate(date string) ([]models.PurchaseResponse, error)
	FindByMonth(date string) ([]models.PurchaseResponse, error)
	FindByPerson(id uuid.UUID) ([]models.PurchaseResponse, error)
	FindAll() ([]models.PurchaseResponse, error)
}

type repositoryPurchase struct {
	db *sql.DB
}

func NewRepositoryPurchase(db *sql.DB) *repositoryPurchase {
	return &repositoryPurchase{db}
}

func (r repositoryPurchase) BeginTransaction() (*sql.Tx, error) {
	return r.db.Begin()
}

func (r repositoryPurchase) Commit(tx *sql.Tx) error {
	return tx.Commit()
}

func (r repositoryPurchase) Rollback(tx *sql.Tx) error {
	return tx.Rollback()
}

func (r repositoryPurchase) Create(tx *sql.Tx, p models.Purchase) (uuid.UUID, error) {
	query := `INSERT INTO purchase(
		id,
		description, 
		amount, 
		"date", 
		place,
		paid, 
		id_payment_type, 
		id_purchase_type, 
		id_credit_card, 
		id_person
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);`

	stmt, err := tx.Prepare(query)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error trying prepare statment: %v", err)
	}

	id, err := uuid.NewUUID()
	if err != nil {
		return uuid.Nil, fmt.Errorf("error trying create uuid: %v", err)
	}

	if _, err = stmt.Exec(
		id,
		p.Description,
		p.Amount,
		p.Date,
		p.Place,
		p.Paid,
		p.IDPaymentType,
		p.IDPurchaseType,
		p.IDCreditCard,
		p.IDPerson,
	); err != nil {
		return uuid.Nil, fmt.Errorf("error trying insert purchase type: %v", err)
	}

	if err := stmt.Close(); err != nil {
		return uuid.Nil, fmt.Errorf("error trying close statment: %v", err)
	}

	return id, nil
}

func (r repositoryPurchase) Update(tx *sql.Tx, p models.Purchase) error {
	query := `UPDATE purchase
		SET description = $1, 
			amount = $2, 
			"date" = $3, 
			place = $4, 
			paid = $5,
			id_payment_type = $6, 
			id_purchase_type = $7, 
			id_credit_card = $8, 
			id_person = $9
		WHERE id = $10;`

	stmt, err := tx.Prepare(query)
	if err != nil {
		return fmt.Errorf("error trying prepare statment: %v", err)
	}

	if _, err = stmt.Exec(
		p.Description,
		p.Amount,
		p.Date,
		p.Place,
		p.Paid,
		p.IDPaymentType,
		p.IDPurchaseType,
		p.IDCreditCard,
		p.IDPerson,
		p.ID,
	); err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error trying update purchase: %v", err)
	}

	if err != nil && err == sql.ErrNoRows {
		return fmt.Errorf("does not exist purchase with this id")
	}

	if err := stmt.Close(); err != nil {
		return fmt.Errorf("error trying close statment: %v", err)
	}

	return nil
}

func (r repositoryPurchase) Delete(tx *sql.Tx, id uuid.UUID) error {
	query := `DELETE FROM purchase WHERE id = $1`

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("error trying prepare statment: %v", err)
	}

	_, err = stmt.Exec(id)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error trying delete purchase: %v", err)
	}

	if err != nil && err == sql.ErrNoRows {
		return fmt.Errorf("does not exist purchase with this id")
	}

	if err := stmt.Close(); err != nil {
		return fmt.Errorf("error trying close statment: %v", err)
	}

	return nil
}

func (r repositoryPurchase) FindByID(id uuid.UUID) (models.PurchaseResponse, error) {
	query := `SELECT 
				p.id, 
				p.description, 
				p.amount, 
				p."date", 
				i.number as installment_number, 
				i.value as installment,
				p.place,
				p.paid,
				pt."name",
				purt."name", 
				cc."owner", 
				per."name" 
			FROM purchase p
			INNER JOIN payment_type pt 
				ON p.id_payment_type = pt.id 
			INNER JOIN purchase_type purt	
				ON p.id_purchase_type = purt.id 
			INNER JOIN credit_card cc	
				ON p.id_credit_card = cc.id
			INNER JOIN person per	
				ON p.id_person = per.id
			LEFT JOIN installment i
				ON p.id = i.purchase_id		
			WHERE p.id = $1
			LIMIT 1;`

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return models.PurchaseResponse{}, fmt.Errorf("error trying prepare statment: %v", err)
	}

	var (
		pt                models.PurchaseResponse
		installmentNumber = sql.NullInt64{}
		installment       = sql.NullFloat64{}
	)
	if err = stmt.QueryRow(id).Scan(
		&pt.ID,
		&pt.Description,
		&pt.Amount,
		&pt.Date,
		&installmentNumber,
		&installment,
		&pt.Place,
		&pt.Paid,
		&pt.PaymentType,
		&pt.PurchaseType,
		&pt.CreditCard,
		&pt.Person,
	); err != nil && err != sql.ErrNoRows {
		return models.PurchaseResponse{}, fmt.Errorf("error trying find purchase: %v", err)
	}

	if err != nil && err == sql.ErrNoRows {
		return models.PurchaseResponse{}, fmt.Errorf("does not exist purchase with this id")
	}

	if err := stmt.Close(); err != nil {
		return models.PurchaseResponse{}, fmt.Errorf("error trying close statment: %v", err)
	}

	return pt, nil
}

func (r repositoryPurchase) FindByDate(date string) ([]models.PurchaseResponse, error) {
	var purchases []models.PurchaseResponse

	query := `SELECT 
				p.id, 
				p.description, 
				p.amount, 
				p."date", 
				i.number as installment_number, 
				i.value as installment, 
				p.place,
				p.paid, 
				pt."name",
				purt."name", 
				cc."owner", 
				per."name"
			FROM purchase p
			INNER JOIN payment_type pt 
				ON p.id_payment_type = pt.id 
			INNER JOIN purchase_type purt	
				ON p.id_purchase_type = purt.id 
			INNER JOIN credit_card cc	
				ON p.id_credit_card = cc.id
			INNER JOIN person per	
				ON p.id_person = per.id 
			LEFT JOIN installment i
				ON p.id = i.purchase_id
			WHERE "date" = $1
			GROUP BY
				p.id, 
				p.description, 
				p.amount, 
				p."date", 
				i.number, 
				i.value, 
				p.place,
				p.paid, 
				pt."name",
				purt."name", 
				cc."owner", 
				per."name";`

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("error trying prepare statment: %v", err)
	}

	rows, err := stmt.Query(date)
	if err != nil {
		return nil, fmt.Errorf("error trying find purchase by date: %v", err)
	}

	installmentNumber := sql.NullInt64{}
	installment := sql.NullFloat64{}

	for rows.Next() {
		var p models.PurchaseResponse
		if err = rows.Scan(
			&p.ID,
			&p.Description,
			&p.Amount,
			&p.Date,
			&installmentNumber,
			&installment,
			&p.Place,
			&p.Paid,
			&p.PaymentType,
			&p.PurchaseType,
			&p.CreditCard,
			&p.Person,
		); err != nil && err != sql.ErrNoRows {
			return nil, fmt.Errorf("error trying scan purchase: %v", err)
		}

		if err != nil && err == sql.ErrNoRows {
			return nil, fmt.Errorf("does not exist purchase  with this date")
		}

		purchases = append(purchases, p)
	}

	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("error trying close rows: %v", err)
	}

	if err := stmt.Close(); err != nil {
		return nil, fmt.Errorf("error trying close statment: %v", err)
	}

	return purchases, nil
}

func (r repositoryPurchase) FindByMonth(date string) ([]models.PurchaseResponse, error) {
	var purchases []models.PurchaseResponse

	query := `SELECT 
				p.id, 
				p.description, 
				p.amount, 
				p."date", 
				i.number as installment_number, 
				i.value as installment, 
				p.place, 
				p.paid,
				pt."name",
				purt."name", 
				cc."owner", 
				per."name"
			FROM purchase p
			INNER JOIN payment_type pt 
				on p.id_payment_type = pt.id 
			INNER JOIN purchase_type purt	
				on p.id_purchase_type = purt.id 
			INNER JOIN credit_card cc	
				on p.id_credit_card = cc.id
			INNER JOIN person per	
				on p.id_person = per.id 
			LEFT JOIN installment i
				ON p.id = i.purchase_id	
			WHERE to_char(p."date", 'YYYY-MM') = $1
			GROUP BY
				p.id, 
				p.description, 
				p.amount, 
				p."date", 
				i.number, 
				i.value, 
				p.place,
				p.paid, 
				pt."name",
				purt."name", 
				cc."owner", 
				per."name"
			ORDER BY p."date";`

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("error trying prepare statment: %v", err)
	}

	rows, err := stmt.Query(date)
	if err != nil {
		return nil, fmt.Errorf("error trying find purchase by date: %v", err)
	}

	installmentNumber := sql.NullInt64{}
	installment := sql.NullFloat64{}

	for rows.Next() {
		var p models.PurchaseResponse
		if err = rows.Scan(
			&p.ID,
			&p.Description,
			&p.Amount,
			&p.Date,
			&installmentNumber,
			&installment,
			&p.Place,
			&p.Paid,
			&p.PaymentType,
			&p.PurchaseType,
			&p.CreditCard,
			&p.Person,
		); err != nil && err != sql.ErrNoRows {
			return nil, fmt.Errorf("error trying scan purchase: %v", err)
		}

		if err != nil && err == sql.ErrNoRows {
			return nil, fmt.Errorf("does not exist purchase  with this date")
		}

		purchases = append(purchases, p)
	}

	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("error trying close rows: %v", err)
	}

	if err := stmt.Close(); err != nil {
		return nil, fmt.Errorf("error trying close statment: %v", err)
	}

	return purchases, nil
}

func (r repositoryPurchase) FindByPerson(id uuid.UUID) ([]models.PurchaseResponse, error) {

	query := `SELECT 
				p.id, 
				p.description, 
				p.amount, 
				p."date", 
				i.number as installment_number, 
				i.value as installment, 
				p.place, 
				p.paid,
				pt."name",
				purt."name", 
				cc."owner", 
				per."name"
			FROM purchase p
			INNER JOIN payment_type pt 
				ON p.id_payment_type = pt.id 
			INNER JOIN purchase_type purt	
				ON p.id_purchase_type = purt.id 
			INNER JOIN credit_card cc	
				ON p.id_credit_card = cc.id
			INNER JOIN person per	
				ON p.id_person = per.id
			LEFT JOIN installment i
				ON p.id = i.purchase_id	
			WHERE p.id_person = $1
			GROUP BY
				p.id, 
				p.description, 
				p.amount, 
				p."date", 
				i.number, 
				i.value, 
				p.place,
				p.paid, 
				pt."name",
				purt."name", 
				cc."owner", 
				per."name"
			ORDER BY p."date";`

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("error trying prepare statment: %v", err)
	}

	rows, err := stmt.Query(id)
	if err != nil {
		return nil, fmt.Errorf("error trying find purchase by person: %v", err)
	}

	var (
		purchases         []models.PurchaseResponse
		installmentNumber = sql.NullInt64{}
		installment       = sql.NullFloat64{}
	)

	for rows.Next() {
		var p models.PurchaseResponse
		if err = rows.Scan(
			&p.ID,
			&p.Description,
			&p.Amount,
			&p.Date,
			&installmentNumber,
			&installment,
			&p.Place,
			&p.Paid,
			&p.PaymentType,
			&p.PurchaseType,
			&p.CreditCard,
			&p.Person,
		); err != nil && err != sql.ErrNoRows {
			return nil, fmt.Errorf("error trying scan purchase: %v", err)
		}

		if err != nil && err == sql.ErrNoRows {
			return nil, fmt.Errorf("does not exist purchase with this person")
		}

		purchases = append(purchases, p)
	}

	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("error trying close rows: %v", err)
	}

	if err := stmt.Close(); err != nil {
		return nil, fmt.Errorf("error trying close statment: %v", err)
	}

	return purchases, nil
}

func (r repositoryPurchase) FindAll() ([]models.PurchaseResponse, error) {
	query := `SELECT 
				p.id, 
				p.description, 
				p.amount, 
				p."date", 
				i.number as installment_number, 
				i.value as installment, 
				p.place, 
				p.paid,
				pt."name",
				purt."name", 
				cc."owner", 
				per."name"
			FROM purchase p
			INNER JOIN payment_type pt 
				ON p.id_payment_type = pt.id 
			INNER JOIN purchase_type purt	
				ON p.id_purchase_type = purt.id 
			INNER JOIN credit_card cc	
				ON p.id_credit_card = cc.id
			INNER JOIN person per	
				ON p.id_person = per.id
			LEFT JOIN installment i
				ON p.id = i.purchase_id
			GROUP BY
				p.id, 
				p.description, 
				p.amount, 
				p."date", 
				i.number, 
				i.value, 
				p.place,
				p.paid, 
				pt."name",
				purt."name", 
				cc."owner", 
				per."name"	
			ORDER BY "date";`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error trying find all purchase: %v", err)
	}

	var (
		purchases         []models.PurchaseResponse
		installmentNumber = sql.NullInt64{}
		installment       = sql.NullFloat64{}
	)

	for rows.Next() {
		var p models.PurchaseResponse
		if err = rows.Scan(
			&p.ID,
			&p.Description,
			&p.Amount,
			&p.Date,
			&installmentNumber,
			&installment,
			&p.Place,
			&p.Paid,
			&p.PaymentType,
			&p.PurchaseType,
			&p.CreditCard,
			&p.Person,
		); err != nil && err != sql.ErrNoRows {
			return nil, fmt.Errorf("error trying scan purchase: %v", err)
		}

		if err != nil && err == sql.ErrNoRows {
			return nil, fmt.Errorf("does not exist purchase")
		}

		purchases = append(purchases, p)
	}

	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("error trying close rows: %v", err)
	}

	return purchases, nil
}
