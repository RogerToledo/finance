package model

import "github.com/google/uuid"

type PaymentType struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
