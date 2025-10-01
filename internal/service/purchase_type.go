package service

import (
	"github.com/google/uuid"
	"github.com/me/finance/internal/models"
	"github.com/me/finance/internal/repository"
)

type PurchaseTypeUseCase interface {
	CreatePurchaseType(pt models.PurchaseType) error
	UpdatePurchaseType(pt models.PurchaseType) error
	DeletePurchaseType(id uuid.UUID) error
	FindPurchaseTypeByID(id uuid.UUID) (models.PurchaseType, error)
	FindAllPurchaseTypes() ([]models.PurchaseType, error)
}

type PurchaseType struct {
	repository repository.RepositoryPurchaseType
}

func NewPurchaseTypeService(r repository.RepositoryPurchaseType) PurchaseTypeUseCase {
	return &PurchaseType{
		repository: r,
	}
}

func (p *PurchaseType) CreatePurchaseType(pt models.PurchaseType) error {
	if err := p.repository.Create(pt); err != nil {
		return err
	}

	return nil
}

func (p *PurchaseType) UpdatePurchaseType(pt models.PurchaseType) error {
	if err := p.repository.Update(pt); err != nil {
		return err
	}

	return nil
}

func (p *PurchaseType) DeletePurchaseType(id uuid.UUID) error {
	if err := p.repository.Delete(id); err != nil {
		return err
	}

	return nil
}

func (p *PurchaseType) FindPurchaseTypeByID(id uuid.UUID) (models.PurchaseType, error) {
	purchaseType, err := p.repository.FindByID(id)
	if err != nil {
		return models.PurchaseType{}, err
	}

	return purchaseType, nil
}

func (p *PurchaseType) FindAllPurchaseTypes() ([]models.PurchaseType, error) {
	purchaseTypes, err := p.repository.FindAll()
	if err != nil {
		return []models.PurchaseType{}, err
	}

	return purchaseTypes, nil
}
