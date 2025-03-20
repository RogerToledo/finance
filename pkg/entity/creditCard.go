package entity

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type CreditCard struct {
	ID    uuid.UUID `json:"id"`
	Owner string    `json:"owner"`
	FinalCardNum string `json:"final_card_num"`
	InvoiceClosingDay int `json:"invoice_closing_day"`
}

func (cc *CreditCard) Validate(removeID bool) error {
	var invalidFields []string

	if !removeID {
		if cc.ID == uuid.Nil {
			invalidFields = append(invalidFields, "ID")
		}
	}
	
	if cc.Owner == "" {
		invalidFields = append(invalidFields, "Owner")
	}
	
	if cc.InvoiceClosingDay == 0 {
		invalidFields = append(invalidFields, "InvoiceClosingDay")
	}
	
	if len(invalidFields) > 0 {
		fields := strings.Join(invalidFields, ", ")

		if len(invalidFields) == 1 {
			return fmt.Errorf("The field %s is required", fields)
		}

		return fmt.Errorf("The fields %s are required", fields)
	}

	return nil
}
