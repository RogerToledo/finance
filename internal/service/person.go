package service

import (
	"github.com/google/uuid"
	"github.com/me/finance/internal/models"
	"github.com/me/finance/internal/repository"
)

type PersonService interface {
	CreatePerson(person models.Person) error
	UpdatePerson(person models.Person) error
	DeletePerson(id uuid.UUID) error
	FindPersonByID(id uuid.UUID) (models.Person, error)
	FindAllPersons() ([]models.Person, error)
}

type Person struct {
	repositoryPerson repository.PersonRepository
}

func NewPersonService(r repository.PersonRepository) PersonService {
	return &Person{
		repositoryPerson: r,
	}
}

func (p *Person) CreatePerson(person models.Person) error {
	if err := p.repositoryPerson.Create(person); err != nil {
		return err
	}

	return nil
}

func (p *Person) UpdatePerson(person models.Person) error {
	if err := p.repositoryPerson.Update(person); err != nil {
		return err
	}

	return nil
}

func (p *Person) DeletePerson(id uuid.UUID) error {
	if err := p.repositoryPerson.Delete(id); err != nil {
		return err
	}

	return nil
}

func (p *Person) FindPersonByID(id uuid.UUID) (models.Person, error) {
	person, err := p.repositoryPerson.FindByID(id)
	if err != nil {
		return models.Person{}, err
	}

	return person, nil
}

func (p *Person) FindAllPersons() ([]models.Person, error) {
	persons, err := p.repositoryPerson.FindAll()
	if err != nil {
		return nil, err
	}
	return persons, nil
}
