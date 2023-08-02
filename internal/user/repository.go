package user

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/jnka9755/go-05DOMAIN/domain"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, user *domain.User) error
	GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.User, error)
	Get(ctx context.Context, id string) (*domain.User, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, user *UpdateUser) error
	Count(ctx context.Context, filters Filters) (int, error)
}

type repository struct {
	log *log.Logger
	db  *gorm.DB
}

func NewRepository(log *log.Logger, db *gorm.DB) Repository {

	return &repository{
		log: log,
		db:  db,
	}
}

func (r *repository) Create(ctx context.Context, user *domain.User) error {

	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		r.log.Println("Error-Repository CreateUser ->", err)
		return err
	}

	r.log.Println("Repository -> Create user with id: ", user.ID)

	return nil
}

func (r *repository) GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.User, error) {

	var users []domain.User

	tx := r.db.WithContext(ctx).Model(&users)
	tx = applyFilters(tx, filters)
	tx = tx.Limit(limit).Offset(offset)

	result := tx.Order("created_at desc").Find(&users)

	if result.Error != nil {
		r.log.Println("Error-Repository GetAllUser ->", result.Error)
		return nil, result.Error
	}

	return users, nil
}

func (r *repository) Get(ctx context.Context, id string) (*domain.User, error) {

	user := domain.User{ID: id}

	if err := r.db.WithContext(ctx).First(&user).Error; err != nil {
		r.log.Println("Error-Repository GetUser ->", err)

		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound{id}
		}

		return nil, err
	}

	return &user, nil
}

func (r *repository) Delete(ctx context.Context, id string) error {

	user := domain.User{ID: id}

	result := r.db.WithContext(ctx).Delete(&user)

	if result.Error != nil {
		r.log.Println("Error-Repository DeleteUser ->", result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		r.log.Printf("User with ID -> '%s' doesn't exist", user.ID)
		return ErrNotFound{id}
	}

	r.log.Println("Repository -> Delete user with id: ", user.ID)

	return nil
}

func (r *repository) Update(ctx context.Context, user *UpdateUser) error {

	values := make(map[string]interface{})

	if user.FirstName != nil {
		values["first_name"] = *user.FirstName
	}

	if user.LastName != nil {
		values["last_name"] = *user.LastName
	}

	if user.Email != nil {
		values["email"] = *user.Email
	}

	if user.Phone != nil {
		values["phone"] = *user.Phone
	}

	result := r.db.WithContext(ctx).Model(&domain.User{}).Where("id = ?", user.ID).Updates(values)

	if result.Error != nil {
		r.log.Println("Error-Repository UdateUser ->", result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		r.log.Printf("User with ID -> '%s' doesn't exist", user.ID)
		return ErrNotFound{user.ID}
	}

	return nil
}

func (r *repository) Count(ctx context.Context, filters Filters) (int, error) {

	var count int64
	tx := r.db.WithContext(ctx).Model(domain.User{})
	tx = applyFilters(tx, filters)

	if err := tx.Count(&count).Error; err != nil {
		r.log.Println("Error-Repository CountUser ->", err)
		return 0, err
	}

	return int(count), nil
}

func applyFilters(tx *gorm.DB, filters Filters) *gorm.DB {

	if filters.FirstName != "" {
		filters.FirstName = fmt.Sprintf("%%%s%%", strings.ToLower(filters.FirstName))
		tx = tx.Where("lower(first_name) like ?", filters.FirstName)
	}

	if filters.LastName != "" {
		filters.LastName = fmt.Sprintf("%%%s%%", strings.ToLower(filters.LastName))
		tx = tx.Where("lower(last_name) like ?", filters.LastName)
	}

	return tx
}
