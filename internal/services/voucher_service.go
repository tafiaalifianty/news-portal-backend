package services

import (
	"final-project-backend/internal/models"
	"final-project-backend/internal/repositories"
)

type IVoucherService interface {
	GetAll() ([]*models.Voucher, error)
	GetUserVouchersByID(userID int64) ([]*models.UserVoucher, error)
}

type voucherService struct {
	voucherRepository     repositories.IVoucherRepository
	userVoucherRepository repositories.IUserVoucherRepository
}

type VoucherServiceConfig struct {
	voucherRepository     repositories.IVoucherRepository
	userVoucherRepository repositories.IUserVoucherRepository
}

func NewVoucherService(c *VoucherServiceConfig) IVoucherService {
	return &voucherService{voucherRepository: c.voucherRepository, userVoucherRepository: c.userVoucherRepository}
}

func (s *voucherService) GetAll() ([]*models.Voucher, error) {
	vouchers, err := s.voucherRepository.GetAll()
	if err != nil {
		return nil, err
	}

	return vouchers, nil
}

func (s *voucherService) GetUserVouchersByID(userID int64) ([]*models.UserVoucher, error) {
	userVouchers, err := s.userVoucherRepository.GetAllByUserID(userID)
	if err != nil {
		return nil, err
	}

	return userVouchers, nil
}
