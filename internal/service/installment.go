package service

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/me/finance/internal/models"
	"github.com/me/finance/internal/repository"
)

type InstallmentService interface {
	CreateInstallment(tx *sql.Tx, purchase models.Purchase) error
	UpdateInstalment(id uuid.UUID) error
	DeleteInstallment(purchaseID uuid.UUID) error
	FindInstallmentByPurchaseID(id uuid.UUID) (models.InstallmentResponse, error)
	FindInstallmentByMonth(month string) (models.InstallmentResponse, error)
	FindInstallmentByNotPaid() (models.InstallmentResponse, error)
}

type Installment struct {
	installmentRepository repository.InstallmentRepository
	creditCardRepository  repository.CreditCardRepository
}

func NewInstallmentService(r repository.InstallmentRepository, cc repository.CreditCardRepository) InstallmentService {
	return &Installment{
		installmentRepository: r,
		creditCardRepository:  cc,
	}
}

func (i *Installment) CreateInstallment(tx *sql.Tx, purchase models.Purchase) error {
	var (
		installment = purchase.Installment
		first       = true
		month       = purchase.Date
		err         error
	)

	slog.Info(fmt.Sprintf("purchaseID: %s", installment.PurchaseID.String()))

	for j := 1; j <= installment.Number; j++ {
		installment.ID = uuid.New()
		installment.Description = fmt.Sprintf("Parcela %d de %d", j, installment.Number)
		installment.Value = purchase.Amount / float64(installment.Number)

		month, err = calculeteDateNextInvoice(i, first, month, purchase.IDCreditCard)
		if err != nil {
			return err
		}
		installment.Month = month
		installment.Paid = false

		if err := i.installmentRepository.Create(tx, installment); err != nil {
			return err
		}
	}

	return nil
}
func (i *Installment) UpdateInstalment(id uuid.UUID) error {
	if err := i.installmentRepository.Update(id); err != nil {
		return fmt.Errorf("error updating installment: %v", err)
	}

	return nil
}

func (i *Installment) DeleteInstallment(purchaseID uuid.UUID) error {
	if err := i.installmentRepository.Delete(purchaseID); err != nil {
		return fmt.Errorf("error deleting installment: %v", err)
	}

	return nil
}

func (i *Installment) FindInstallmentByPurchaseID(id uuid.UUID) (models.InstallmentResponse, error) {
	installments, err := i.installmentRepository.FindByPurchaseID(id)
	if err != nil {
		return models.InstallmentResponse{}, fmt.Errorf("error finding installment by purchaseID: %v", err)
	}

	response := processInstallmentResponse(installments)

	return response, nil

}

func (i *Installment) FindInstallmentByMonth(month string) (models.InstallmentResponse, error) {
	installments, err := i.installmentRepository.FindByMonth(month)
	if err != nil {
		return models.InstallmentResponse{}, fmt.Errorf("error finding installment by month: %v", err)
	}

	response := processInstallmentResponse(installments)

	return response, nil
}

func (i *Installment) FindInstallmentByNotPaid() (models.InstallmentResponse, error) {
	installments, err := i.installmentRepository.FindByNotPaid()
	if err != nil {
		return models.InstallmentResponse{}, fmt.Errorf("error finding installment by not paid: %v", err)
	}

	response := processInstallmentResponse(installments)

	return response, nil
}

func calculeteDateNextInvoice(i *Installment, first bool, date string, id uuid.UUID) (string, error) {
	var (
		cc  models.CreditCard
		err error
		msg error
	)

	completeDate, err1 := time.Parse("2006-01-02", date)
	if err1 != nil {
		msg = fmt.Errorf("error parsing date: %v", err1)
	}

	if first {
		cc, err = i.creditCardRepository.FindByID(id)
		if err != nil {
			return "", err
		}

		if completeDate.Day() >= cc.InvoiceClosingDay {
			completeDate = completeDate.AddDate(0, 1, 0)

			return completeDate.Format("2006-01-02"), nil
		} else {
			return completeDate.Format("2006-01-02"), nil
		}
	}

	newDate, err2 := time.Parse("2006-01-02", date)
	if err1 != nil && err2 != nil {
		return "", msg
	}

	newDate = newDate.AddDate(0, 1, 0)

	return newDate.Format("2006-01-02"), nil
}

func processInstallmentResponse(installments []models.Installment) models.InstallmentResponse {
	paid, toPay := calculateTotal(installments)

	response := models.InstallmentResponse{
		Response: installments,
		Paid:     paid,
		ToPay:    toPay,
		Total:    paid + toPay}

	return response
}

func calculateTotal(installments []models.Installment) (float64, float64) {
	var toPay, paid float64

	for _, installment := range installments {
		if installment.Paid {
			paid += installment.Value
		} else {
			toPay += installment.Value
		}
	}

	return paid, toPay
}
