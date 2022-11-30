package repositories

import (
	"final-project-backend/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IUserVoucherRepository interface {
	GetAll() ([]*models.UserVoucher, error)
	GetAllByUserID(userID int64) ([]*models.UserVoucher, error)
	GetByCode(code string) (*models.UserVoucher, error)
	Insert(userVoucher *models.UserVoucher) (*models.UserVoucher, error)
	Update(userVoucher *models.UserVoucher) (*models.UserVoucher, int, error)
}

type userVoucherRepository struct {
	db *gorm.DB
}

type UserVoucherRepositoryConfig struct {
	db *gorm.DB
}

func NewUserVoucherRepository(c *UserVoucherRepositoryConfig) IUserVoucherRepository {
	return &userVoucherRepository{
		db: c.db,
	}
}

func (r *userVoucherRepository) GetAll() ([]*models.UserVoucher, error) {
	var UserVouchers []*models.UserVoucher
	result := r.db.Order("id asc").Find(&UserVouchers)

	if result.Error != nil {
		return nil, result.Error
	}

	return UserVouchers, nil
}

func (r *userVoucherRepository) GetAllByUserID(userID int64) ([]*models.UserVoucher, error) {
	var UserVouchers []*models.UserVoucher
	result := r.db.Where("user_id = ?", userID).Joins("ReceivedFromUser").Joins("Voucher").Order("status asc,valid_until desc, id asc").Find(&UserVouchers)

	if result.Error != nil {
		return nil, result.Error
	}

	return UserVouchers, nil
}

func (r *userVoucherRepository) GetByCode(code string) (*models.UserVoucher, error) {
	var userVoucher *models.UserVoucher
	result := r.db.Where("code = ?", code).Joins("Voucher").First(&userVoucher)

	if result.Error != nil {
		return nil, result.Error
	}

	return userVoucher, nil
}

func (r *userVoucherRepository) Insert(userVoucher *models.UserVoucher) (*models.UserVoucher, error) {
	result := r.db.Create(userVoucher)
	if result.Error != nil {
		return nil, result.Error
	}

	return userVoucher, nil
}

func (r *userVoucherRepository) Update(userVoucher *models.UserVoucher) (*models.UserVoucher, int, error) {
	result := r.db.Model(&userVoucher).Clauses(clause.Returning{}).Updates(userVoucher)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return userVoucher, int(result.RowsAffected), nil
}
