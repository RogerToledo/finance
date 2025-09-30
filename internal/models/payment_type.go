package models

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type PaymentType struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func (pt *PaymentType) Validate(removeID bool) error {
	var invalidFields []string

	if !removeID {
		if pt.ID == uuid.Nil {
			invalidFields = append(invalidFields, "ID")
		}
	}
	
	if pt.Name == "" {
		invalidFields = append(invalidFields, "Name")
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
