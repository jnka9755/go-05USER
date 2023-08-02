package user

import (
	"context"
	"log"

	"github.com/jnka9755/go-05DOMAIN/domain"
)

type (
	Business interface {
		Create(ctx context.Context, request *CreateReq) (*domain.User, error)
		GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.User, error)
		Get(ctx context.Context, id string) (*domain.User, error)
		Delete(ctx context.Context, id string) error
		Update(ctx context.Context, request *UpdateReq) error
		Count(ctx context.Context, filters Filters) (int, error)
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
		ID        string
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

func (b business) Create(ctx context.Context, request *CreateReq) (*domain.User, error) {

	user := domain.User{
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Email:     request.Email,
		Phone:     request.Phone,
	}

	if err := b.repository.Create(ctx, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (b business) GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.User, error) {

	users, err := b.repository.GetAll(ctx, filters, offset, limit)

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (b business) Get(ctx context.Context, id string) (*domain.User, error) {

	user, err := b.repository.Get(ctx, id)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (b business) Delete(ctx context.Context, id string) error {

	return b.repository.Delete(ctx, id)
}

func (b business) Update(ctx context.Context, request *UpdateReq) error {

	user := UpdateUser{
		ID:        request.ID,
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Email:     request.Email,
		Phone:     request.Phone,
	}

	return b.repository.Update(ctx, &user)
}

func (b business) Count(ctx context.Context, filters Filters) (int, error) {
	return b.repository.Count(ctx, filters)
}
