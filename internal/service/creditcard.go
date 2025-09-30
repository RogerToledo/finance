package service

import (
	"github.com/google/uuid"
	"github.com/me/finance/internal/models"
	"github.com/me/finance/internal/repository"
)

type CreditCardService interface {
	CreateCreditCard(cc models.CreditCard) error
	UpdateCreditCard(cc models.CreditCard) error
	DeleteCreditCard(id uuid.UUID) error
	FindCreditCardByID(id uuid.UUID) (models.CreditCard, error)
	FindAllCreditCards() ([]models.CreditCard, error)
}

type CreditCard struct {
	creditCardRepository repository.CreditCardRepository
}

func NewCreditCardService(r repository.CreditCardRepository) CreditCardService {
	return &CreditCard{
		creditCardRepository: r,
	}
}

func (c *CreditCard) CreateCreditCard(cc models.CreditCard) error {
	if err := c.creditCardRepository.Create(cc); err != nil {
		return err
	}

	return nil
}

func (c *CreditCard) UpdateCreditCard(cc models.CreditCard) error {
	if err := c.creditCardRepository.Update(cc); err != nil {
		return err
	}

	return nil
}

func (c *CreditCard) DeleteCreditCard(id uuid.UUID) error {
	if err := c.creditCardRepository.Delete(id); err != nil {
		return err
	}

	return nil
}

func (c *CreditCard) FindCreditCardByID(id uuid.UUID) (models.CreditCard, error) {
	cc, err := c.creditCardRepository.FindByID(id)
	if err != nil {
		return models.CreditCard{}, err
	}

	return cc, nil
}

func (c *CreditCard) FindAllCreditCards() ([]models.CreditCard, error) {
	cc, err := c.creditCardRepository.FindAll()
	if err != nil {
		return []models.CreditCard{}, err
	}

	return cc, nil
}
