package repositories

import (
	"final-project-backend/internal/models"

	"gorm.io/gorm"
)

type IVoucherRepository interface {
	GetAll() ([]*models.Voucher, error)
}

type voucherRepository struct {
	db *gorm.DB
}

type VoucherRepositoryConfig struct {
	db *gorm.DB
}

func NewVoucherRepository(c *VoucherRepositoryConfig) IVoucherRepository {
	return &voucherRepository{
		db: c.db,
	}
}

func (r *voucherRepository) GetAll() ([]*models.Voucher, error) {
	var vouchers []*models.Voucher
	result := r.db.Order("id asc").Find(&vouchers)

	if result.Error != nil {
		return nil, result.Error
	}

	return vouchers, nil
}
