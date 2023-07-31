package user

import (
	"log"

	"github.com/jnka9755/go-05DOMAIN/domain"
)

type (
	Business interface {
		Create(request *CreateReq) (*domain.User, error)
		GetAll(filters Filters, offset, limit int) ([]domain.User, error)
		Get(id string) (*domain.User, error)
		Delete(id string) error
		Update(id string, request *UpdateReq) error
		Count(filters Filters) (int, error)
	}

	business struct {
		log        *log.Logger
		repository Repository
	}

	Filters struct {
		FirstName string
		LastName  string
	}

	UpdateUser struct {
		FirstName *string
		LastName  *string
		Email     *string
		Phone     *string
	}
)

func NewBusiness(log *log.Logger, repository Repository) Business {
	return &business{
		log:        log,
		repository: repository,
	}
}

func (b business) Create(request *CreateReq) (*domain.User, error) {

	user := domain.User{
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Email:     request.Email,
		Phone:     request.Phone,
	}

	b.log.Println("Create user Business")
	if err := b.repository.Create(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (b business) GetAll(filters Filters, offset, limit int) ([]domain.User, error) {

	b.log.Println("GetAll user Business")
	users, err := b.repository.GetAll(filters, offset, limit)

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (b business) Get(id string) (*domain.User, error) {

	b.log.Println("Get user Business")
	user, err := b.repository.Get(id)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (b business) Delete(id string) error {

	b.log.Println("Delete user Business")
	return b.repository.Delete(id)
}

func (b business) Update(id string, request *UpdateReq) error {

	b.log.Println("Update user Business")

	user := UpdateUser{
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Email:     request.Email,
		Phone:     request.Phone,
	}

	return b.repository.Update(id, &user)
}

func (b business) Count(filters Filters) (int, error) {
	return b.repository.Count(filters)
}
