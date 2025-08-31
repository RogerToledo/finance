package dto

import (
	"github.com/google/uuid"
	"github.com/me/finance/pkg/entity"
)

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

func (p *PurchaseRequest) ToEntity() (entity.Purchase, error) {
	var installment entity.Installment

	installment.Number = p.InstallmentNumber
	installment.Value  = p.Installment

	purchase := entity.Purchase{
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