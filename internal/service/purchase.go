package service

import (
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/me/finance/internal/models"
	"github.com/me/finance/internal/repository"
)

type PurchaseService interface {
	CreatePurchase(purchase models.Purchase) error
	UpdatePurchase(purchase models.Purchase) error
	DeletePurchase(id uuid.UUID) error
	FindPurchaseByID(id uuid.UUID) (models.PurchaseResponse, error)
	FindPurchaseByDate(date string) (models.PurchaseResponseTotal, error)
	FindPurchaseByMonth(date string) (models.PurchaseResponseTotal, error)
	FindPurchaseByPerson(id uuid.UUID) (models.PurchaseResponseTotal, error)
	FindAllPurchases() ([]models.PurchaseResponse, error)
}

type Purchase struct {
	purchaseRepository    repository.PurchaseRepository
	installmentRepository repository.InstallmentRepository
	creditCardRepository  repository.CreditCardRepository
}

func NewPurchaseService(p repository.PurchaseRepository, i repository.InstallmentRepository, cc repository.CreditCardRepository) PurchaseService {
	return &Purchase{
		purchaseRepository:    p,
		installmentRepository: i,
		creditCardRepository:  cc,
	}
}

func (p *Purchase) CreatePurchase(purchase models.Purchase) error {
	tx, err := p.purchaseRepository.BeginTransaction()
	if err != nil {
		return fmt.Errorf("error on begin transaction: %v", err)
	}

	var savedID uuid.UUID

	if savedID, err = p.purchaseRepository.Create(tx, purchase); err != nil {
		return err
	}

	slog.Info(fmt.Sprintf("savedID: %s", savedID.String()))

	purchase.Installment.PurchaseID = savedID

	i := NewInstallmentService(p.installmentRepository, p.creditCardRepository)

	if err := i.CreateInstallment(tx, purchase); err != nil {
		p.purchaseRepository.Rollback(tx)

		return err
	}

	p.purchaseRepository.Commit(tx)

	return nil
}

func (p *Purchase) UpdatePurchase(purchase models.Purchase) error {
	tx, err := p.purchaseRepository.BeginTransaction()
	if err != nil {
		return fmt.Errorf("error on begin transaction: %v", err)
	}

	if err := p.purchaseRepository.Update(tx, purchase); err != nil {
		return err
	}

	purchase.Installment.PurchaseID = purchase.ID

	if err := p.installmentRepository.Delete(purchase.ID); err != nil {
		p.purchaseRepository.Rollback(tx)

		return err
	}

	i := NewInstallmentService(p.installmentRepository, p.creditCardRepository)

	if err := i.CreateInstallment(tx, purchase); err != nil {
		p.purchaseRepository.Rollback(tx)

		return err
	}

	p.purchaseRepository.Commit(tx)

	return nil
}

func (p *Purchase) DeletePurchase(id uuid.UUID) error {
	tx, err := p.purchaseRepository.BeginTransaction()
	if err != nil {
		return fmt.Errorf("error on begin transaction: %v", err)
	}

	if err := p.purchaseRepository.Delete(tx, id); err != nil {
		return err
	}

	if err := p.installmentRepository.Delete(id); err != nil {
		p.purchaseRepository.Rollback(tx)

		return err
	}

	p.purchaseRepository.Commit(tx)

	return nil
}

func (p *Purchase) FindPurchaseByID(id uuid.UUID) (models.PurchaseResponse, error) {
	purchase, err := p.purchaseRepository.FindByID(id)
	if err != nil {
		return models.PurchaseResponse{}, err
	}

	return purchase, err
}

func (p *Purchase) FindPurchaseByDate(date string) (models.PurchaseResponseTotal, error) {
	purchases, err := p.purchaseRepository.FindByDate(date)
	if err != nil {
		return models.PurchaseResponseTotal{}, err
	}

	response := processPurchaseResponse(purchases)

	return response, err
}

func (p *Purchase) FindPurchaseByMonth(date string) (models.PurchaseResponseTotal, error) {
	purchases, err := p.purchaseRepository.FindByMonth(date)
	if err != nil {
		return models.PurchaseResponseTotal{}, err
	}

	response := processPurchaseResponse(purchases)

	return response, err
}

func (p *Purchase) FindPurchaseByPerson(personID uuid.UUID) (models.PurchaseResponseTotal, error) {
	purchases, err := p.purchaseRepository.FindByPerson(personID)
	if err != nil {
		return models.PurchaseResponseTotal{}, err
	}

	response := processPurchaseResponse(purchases)

	return response, err
}

func (p *Purchase) FindAllPurchases() ([]models.PurchaseResponse, error) {
	purchases, err := p.purchaseRepository.FindAll()
	if err != nil {
		return nil, err
	}

	return purchases, err
}

func processPurchaseResponse(purchases []models.PurchaseResponse) models.PurchaseResponseTotal {
	total := 0.0

	for _, purchase := range purchases {
		total += purchase.Amount
	}

	response := models.PurchaseResponseTotal{
		Responses: purchases,
		Quantity:  len(purchases),
		Total:     total,
	}

	return response
}
