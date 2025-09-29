package repository

// import (
// 	"database/sql"
// 	"sync"
// )

// type Repository struct {
// 	Person       *repositoryPerson
// 	// CreditCard   *repositoryCreditCard
// 	// PaymentType  *repositoryPaymentType
// 	// PurchaseType *repositoryPurchaseType
// 	// Purchase     *repositoryPurchase
// 	// Installment  *repositoryInstallment
// 	// All          *repositoryAll
// }

// var (
// 	once     sync.Once
// 	instance *Repository
// )

// func NewRepository(db *sql.DB) *Repository {
// 	if instance == nil {
// 		once.Do(func() {
// 			person := NewRepositoryPerson(db)
// 			// creditCard := NewRepositoryCreditCard(db)
// 			// paymentType := NewRepositoryPaymentType(db)
// 			// purchaseType := NewRepositoryPurchaseType(db)
// 			// purchase := NewRepositoryPurchase(db)
// 			// installment := NewRepositoryInstallment(db)
// 			// all := NewRepositoryAll(db)

// 			instance = &Repository{
// 				// Person:       person,
// 				// CreditCard:   creditCard,
// 				// PaymentType:  paymentType,
// 				// PurchaseType: purchaseType,
// 				// Purchase:     purchase,
// 				// Installment:  installment,
// 				// All:          all,
// 			}
// 		})
// 	}

// 	return instance
// }