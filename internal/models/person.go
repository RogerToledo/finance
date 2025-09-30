package models

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type Person struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func (p *Person) Validate(removeID bool) error {
	var invalidFields []string

	if !removeID {
		if p.ID == uuid.Nil {
			invalidFields = append(invalidFields, "ID")
		}
	}
	
	if p.Name == "" {
		invalidFields = append(invalidFields, "Name")
	}
	
	if len(invalidFields) > 0 {
		fields := strings.Join(invalidFields, ", ")

		if len(invalidFields) == 1 {
			return fmt.Errorf("the field %s is required", fields)
		}

		return fmt.Errorf("the fields %s are required", fields)
	}

	return nil
}