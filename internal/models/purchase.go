package models

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type Purchase struct {
	ID             uuid.UUID `json:"id"`	
	Description    string  `json:"description"`
	Amount         float64 `json:"amount"`
	Date           string  `json:"date"` 
	Installment    Installment `json:"installment"`
	Place	       string  `json:"place"`
	Paid		   bool    `json:"paid"`
	IDPaymentType  uuid.UUID `json:"id_payment_type"`
	IDCreditCard   uuid.UUID `json:"id_credit_card"`
	IDPurchaseType uuid.UUID `json:"id_purchase_type"`
	IDPerson	   uuid.UUID `json:"id_person"`
}

type PurchaseRequest struct {
	ID                uuid.UUID `json:"id"`
	Description       string  `json:"description"`
	Amount            float64 `json:"amount"`
	Date              string  `json:"date"` 
	InstallmentNumber int     `json:"installment_number"`
	Installment       float64 `json:"installment"`
	Place	          string  `json:"place"`
	Paid			  bool    `json:"paid"`
	IDPaymentType     uuid.UUID `json:"id_payment_type"`
	IDCreditCard	  uuid.UUID `json:"id_credit_card"`
	IDPurchaseType    uuid.UUID `json:"id_purchase_type"`
	IDPerson	      uuid.UUID `json:"id_person"`
}

type PurchaseResponse struct {
	ID                uuid.UUID `json:"id"`
	Description       string  `json:"description"`
	Amount            float64 `json:"amount"`
	Date              string  `json:"date"` 
	InstallmentNumber int     `json:"installment_number"`
	Installment       float64 `json:"installment"`
	Place	          string  `json:"place"`
	Paid			  bool    `json:"paid"`
	PaymentType       string  `json:"payment_type"`
	CreditCard	      string  `json:"credit_card"`
	PurchaseType      string  `json:"purchase_type"`
	Person	          string  `json:"person"`
}

type PurchaseResponseTotal struct {
	Responses []PurchaseResponse `json:"responses"`
	Quantity  int                `json:"quantity"`
	Total     float64            `json:"total"`
}

func (p *PurchaseRequest) ToEntity() (Purchase, error) {
	var installment Installment

	installment.Number = p.InstallmentNumber
	installment.Value  = p.Installment

	purchase := Purchase{
		ID:             p.ID,
		Description:    p.Description,
		Amount:         p.Amount,
		Date:           p.Date,
		Installment:    installment,
		Place:	        p.Place,
		Paid:	        p.Paid,
		IDPaymentType:  p.IDPaymentType,
		IDCreditCard:	p.IDCreditCard,
		IDPurchaseType: p.IDPurchaseType,
		IDPerson:	    p.IDPerson,
	}

	return purchase, nil
}

func (p *Purchase) Validate() error {
	var invalidFields []string

	if p.Amount <= 0{
		invalidFields = append(invalidFields, "Amount")
	}

	if p.Date == "" {
		invalidFields = append(invalidFields, "Data")
	}

	if err := ValidateDate(p.Date); err != nil {
		invalidFields = append(invalidFields, "Date")
	}

	if p.IDPaymentType == uuid.Nil {
		invalidFields = append(invalidFields, "ID of Payment Type")
	}

	if p.IDCreditCard == uuid.Nil {
		invalidFields = append(invalidFields, "ID of Credit Card")
	}

	if p.IDPurchaseType == uuid.Nil {
		invalidFields = append(invalidFields, "ID of Purchase Type")
	}

	if p.IDPerson == uuid.Nil {
		invalidFields = append(invalidFields, "ID of Person")
	}

	if len(invalidFields) > 0 {
		fields := strings.Join(invalidFields, ", ")

		if len(invalidFields) == 1 {
			return fmt.Errorf("the field %s is required", fields)
		} else {
			return fmt.Errorf("the fields %s are required", fields)
		}
	}

	return nil
}
