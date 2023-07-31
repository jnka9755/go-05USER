package user

import (
	"fmt"
	"log"
	"strings"

	"github.com/jnka9755/go-05DOMAIN/domain"

	"gorm.io/gorm"
)

type Repository interface {
	Create(user *domain.User) error
	GetAll(filters Filters, offset, limit int) ([]domain.User, error)
	Get(id string) (*domain.User, error)
	Delete(id string) error
	Update(id string, user *UpdateUser) error
	Count(filters Filters) (int, error)
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

func (r *repository) Create(user *domain.User) error {

	if err := r.db.Create(user).Error; err != nil {
		r.log.Println("Repository ->", err)
		return err
	}

	r.log.Println("Repository -> Create user with id: ", user.ID)

	return nil
}

func (r *repository) GetAll(filters Filters, offset, limit int) ([]domain.User, error) {

	r.log.Println("GetAll user Repository")
	var users []domain.User

	tx := r.db.Model(&users)
	tx = applyFilters(tx, filters)
	tx = tx.Limit(limit).Offset(offset)

	result := tx.Order("created_at desc").Find(&users)

	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}

func (r *repository) Get(id string) (*domain.User, error) {

	r.log.Println("Get user Repository")
	user := domain.User{ID: id}

	if err := r.db.First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *repository) Delete(id string) error {

	r.log.Println("Delete user Repository")
	user := domain.User{ID: id}

	if err := r.db.Delete(&user).Error; err != nil {
		return err
	}

	r.log.Println("Repository -> Delete user with id: ", user.ID)

	return nil
}

func (r *repository) Update(id string, user *UpdateUser) error {

	r.log.Println("Udate user Repository")

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

	if err := r.db.Model(&domain.User{}).Where("id = ?", id).Updates(values).Error; err != nil {
		return err
	}

	return nil
}

func (r *repository) Count(filters Filters) (int, error) {

	var count int64
	tx := r.db.Model(domain.User{})
	tx = applyFilters(tx, filters)

	if err := tx.Count(&count).Error; err != nil {
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
