package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

func ValidateID(idRequest string) (uuid.UUID, error) {
	id, err := uuid.Parse(idRequest)
	if err != nil {
		return uuid.Nil, fmt.Errorf("Error converting ID to UUID: %v", err)
	}

	return id, nil
}

func ConverDateDB(date string) (string, error) {
	t1, err := time.Parse("2006-01-02", date)
	if err != nil {
		return "", fmt.Errorf("Error converting date: %v", err)
	}

	dateFormated := t1.Format("2006-01-02")

	return dateFormated, nil
}

func ConverDate(date string) (string, error) {
	t1, err := time.Parse("2006-01-02", date)
	if err != nil {
		return "", fmt.Errorf("Error converting date: %v", err)
	}

	dateFormated := t1.Format("02/01/2006")

	return dateFormated, nil
}

func ValidateDate(date string) error {
	if _, err := time.Parse("2006-01-02", date); err != nil {
		return fmt.Errorf("The date is invalid")
	}

	return nil
}

func ValidateYearMonth(date string) error {
	if _, err := time.Parse("2006-01", date); err != nil {
		return fmt.Errorf("The date is invalid")
	}

	return nil
}