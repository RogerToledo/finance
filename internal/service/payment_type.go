package service

import (
	"github.com/google/uuid"
	"github.com/me/finance/internal/models"
	"github.com/me/finance/internal/repository"
)

type PaymentTypeService interface {
	CreatePaymentType(paymentType models.PaymentType) error
	UpdatePaymentType(paymentType models.PaymentType) error
	DeletePaymentType(id uuid.UUID) error
	FindPaymentTypeByID(id uuid.UUID) (models.PaymentType, error)
	FindAllPaymentTypes() ([]models.PaymentType, error)
}

type PaymentType struct {
	paymentTypeRepository repository.PaymentTypeRepository
}

func NewPaymentTypeService(r repository.PaymentTypeRepository) PaymentTypeService {
	return &PaymentType{
		paymentTypeRepository: r,
	}
}

func (p *PaymentType) CreatePaymentType(paymentType models.PaymentType) error {
	if err := p.paymentTypeRepository.Create(paymentType); err != nil {
		return err
	}

	return nil
}

func (p *PaymentType) UpdatePaymentType(paymentType models.PaymentType) error {
	if err := p.paymentTypeRepository.Update(paymentType); err != nil {
		return err
	}

	return nil
}

func (p *PaymentType) DeletePaymentType(id uuid.UUID) error {
	if err := p.paymentTypeRepository.Delete(id); err != nil {
		return err
	}

	return nil
}

func (p *PaymentType) FindPaymentTypeByID(id uuid.UUID) (models.PaymentType, error) {
	paymentType, err := p.paymentTypeRepository.FindByID(id)
	if err != nil {
		return models.PaymentType{}, err
	}

	return paymentType, nil
}

func (p *PaymentType) FindAllPaymentTypes() ([]models.PaymentType, error) {
	paymentTypes, err := p.paymentTypeRepository.FindAll()
	if err != nil {
		return nil, err
	}

	return paymentTypes, nil
}
