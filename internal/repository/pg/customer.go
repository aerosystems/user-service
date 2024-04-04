package pg

import (
	"context"
	"errors"
	"github.com/aerosystems/customer-service/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CustomerRepo struct {
	db *gorm.DB
}

func NewCustomerRepo(db *gorm.DB) *CustomerRepo {
	return &CustomerRepo{
		db: db,
	}
}

func (r *CustomerRepo) GetByUuid(ctx context.Context, uuid uuid.UUID) (*models.Customer, error) {
	var user models.Customer
	result := r.db.Find(&user, "uuid = ?", uuid.String())
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

func (r *CustomerRepo) Create(ctx context.Context, user *models.Customer) error {
	result := r.db.Create(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *CustomerRepo) Update(ctx context.Context, user *models.Customer) error {
	result := r.db.Save(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *CustomerRepo) Delete(ctx context.Context, user *models.Customer) error {
	result := r.db.Delete(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
