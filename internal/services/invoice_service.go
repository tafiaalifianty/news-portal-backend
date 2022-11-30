package services

import (
	"errors"
	"time"

	"final-project-backend/internal/constants"
	errn "final-project-backend/internal/errors"
	"final-project-backend/internal/helpers"
	"final-project-backend/internal/models"
	"final-project-backend/internal/repositories"

	"gorm.io/gorm"
)

type IInvoiceService interface {
	GetAll() ([]*models.Invoice, error)
	GetUserInvoices(userID int64) ([]*models.Invoice, error)
	GetByID(id int64) (*models.Invoice, error)
	GetByCode(code string) (*models.Invoice, error)
	GetUserInvoiceByCode(code string, userID int64) (*models.Invoice, error)
	Create(invoice *models.Invoice) (*models.Invoice, error)
	UpdateStatus(code string, status models.InvoiceStatus) (*models.Invoice, []*models.Gift, []*models.Voucher, error)
}

type invoiceService struct {
	invoiceRepository          repositories.IInvoiceRepository
	subscriptionRepository     repositories.ISubscriptionRepository
	userSubscriptionRepository repositories.IUserSubscriptionRepository
	giftRepository             repositories.IGiftRepository
	userGiftRepository         repositories.IUserGiftRepository
	userRepository             repositories.IUserRepository
	voucherRepository          repositories.IVoucherRepository
	userVoucherRepository      repositories.IUserVoucherRepository
	userSpendingRepository     repositories.IUserSpendingRepository
}

type InvoiceServiceConfig struct {
	invoiceRepository          repositories.IInvoiceRepository
	subscriptionRepository     repositories.ISubscriptionRepository
	userSubscriptionRepository repositories.IUserSubscriptionRepository
	giftRepository             repositories.IGiftRepository
	userGiftRepository         repositories.IUserGiftRepository
	userRepository             repositories.IUserRepository
	voucherRepository          repositories.IVoucherRepository
	userVoucherRepository      repositories.IUserVoucherRepository
	userSpendingRepository     repositories.IUserSpendingRepository
}

func NewInvoiceService(c *InvoiceServiceConfig) IInvoiceService {
	return &invoiceService{
		invoiceRepository:          c.invoiceRepository,
		subscriptionRepository:     c.subscriptionRepository,
		userSubscriptionRepository: c.userSubscriptionRepository,
		giftRepository:             c.giftRepository,
		userGiftRepository:         c.userGiftRepository,
		userRepository:             c.userRepository,
		voucherRepository:          c.voucherRepository,
		userVoucherRepository:      c.userVoucherRepository,
		userSpendingRepository:     c.userSpendingRepository,
	}
}

func (s *invoiceService) GetAll() ([]*models.Invoice, error) {
	invoices, err := s.invoiceRepository.GetAll(&models.Invoice{})
	if err != nil {
		return nil, err
	}

	return invoices, nil
}

func (s *invoiceService) GetUserInvoices(userID int64) ([]*models.Invoice, error) {
	invoices, err := s.invoiceRepository.GetAll(&models.Invoice{UserID: userID})
	if err != nil {
		return nil, err
	}

	return invoices, nil
}

func (s *invoiceService) GetUserInvoicesCurrentMonth(userID int64) ([]*models.Invoice, error) {
	invoices, err := s.invoiceRepository.GetAllCurrentMonth(&models.Invoice{UserID: userID})
	if err != nil {
		return nil, err
	}

	return invoices, nil
}

func (s *invoiceService) GetByID(id int64) (*models.Invoice, error) {
	invoice, err := s.invoiceRepository.GetByID(id)
	if err != nil {
		return nil, err
	}

	return invoice, nil
}

func (s *invoiceService) GetByCode(code string) (*models.Invoice, error) {
	invoice, err := s.invoiceRepository.GetByCode(code)
	if err != nil {
		return nil, err
	}

	return invoice, nil
}

func (s *invoiceService) GetUserInvoiceByCode(code string, userID int64) (*models.Invoice, error) {
	invoice, err := s.invoiceRepository.GetByCode(code)
	if err != nil {
		return nil, err
	}

	if invoice.UserID != userID {
		return nil, errn.ErrNotAuthorized
	}

	return invoice, nil
}

func (s *invoiceService) Create(invoice *models.Invoice) (*models.Invoice, error) {
	var userVoucherID int64
	if invoice.VoucherCode != "" {
		userVoucher, err := s.userVoucherRepository.GetByCode(invoice.VoucherCode)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errn.ErrInvalidVoucher
		}

		if err != nil {
			return nil, err
		}

		if userVoucher.Status != constants.AVAILABLE || time.Now().After(userVoucher.ValidUntil) {
			return nil, errn.ErrVoucherExpired
		}

		invoice.Total = invoice.OriginalPrice - userVoucher.Voucher.Discount
		if invoice.Total < 0 {
			invoice.Total = 0
		}
		userVoucherID = userVoucher.ID
	}

	createdInvoice, err := s.invoiceRepository.Insert(invoice)
	if err != nil {
		return nil, err
	}

	if userVoucherID != 0 {
		s.userVoucherRepository.Update(&models.UserVoucher{
			ID:     userVoucherID,
			Status: constants.PENDING,
		})
	}

	return createdInvoice, nil
}

func (s *invoiceService) UpdateStatus(code string, status models.InvoiceStatus) (*models.Invoice, []*models.Gift, []*models.Voucher, error) {
	giftsSent := []*models.Gift{}
	vouchersSent := []*models.Voucher{}
	currentInvoice, err := s.invoiceRepository.GetByCode(code)
	if err != nil {
		return nil, nil, nil, err
	}

	if (status != models.REJECTED && int(status)-int(currentInvoice.Status) != 1) || (status == models.REJECTED && currentInvoice.Status != models.PROCESSED) {
		return nil, nil, nil, errn.ErrInvalidStatusUpdate
	}

	now := time.Now()
	updatedInvoice, rowsAffected, err := s.invoiceRepository.Update(&models.Invoice{
		Model: gorm.Model{
			ID: currentInvoice.ID,
		},
		Code:   code,
		Status: status,
		PaidAt: &now,
	})

	if rowsAffected == 0 {
		return nil, nil, nil, gorm.ErrRecordNotFound
	}

	if err != nil {
		return nil, nil, nil, err
	}

	if updatedInvoice.VoucherCode != "" {
		userVoucher, err := s.userVoucherRepository.GetByCode(updatedInvoice.VoucherCode)
		if err != nil {
			newVoucherStatus := constants.USED
			if updatedInvoice.Status == models.REJECTED {
				newVoucherStatus = constants.AVAILABLE
			}
			s.userVoucherRepository.Update(&models.UserVoucher{
				ID:     userVoucher.ID,
				Status: newVoucherStatus,
			})
		}
	}

	if updatedInvoice.Status == models.COMPLETED {
		subscription, err := s.subscriptionRepository.GetByID(updatedInvoice.SubscriptionID)

		if err != nil {
			return nil, nil, nil, err
		}

		userSubscription := &models.UserSubscriptions{
			UserID:         updatedInvoice.UserID,
			SubscriptionID: updatedInvoice.SubscriptionID,
			RemainingQuota: subscription.Quota,
			DateStarted:    time.Now(),
			DateEnded:      time.Now().AddDate(0, 1, 0),
		}

		_, err = s.userSubscriptionRepository.Insert(userSubscription)
		if err != nil {
			return nil, nil, nil, err
		}

		currentTime := time.Now()
		currentMonthNum := int(currentTime.Month())
		currentYearNum := int(currentTime.Year())

		userSpending, err := s.userSpendingRepository.GetByPrimaryKeys(updatedInvoice.UserID, currentMonthNum, currentYearNum)
		if err != nil {
			return updatedInvoice, giftsSent, vouchersSent, nil
		}

		giftsSent, _ = s.SendGiftIfAble(updatedInvoice.UserID, updatedInvoice, userSpending)

		vouchersSent, _ = s.SendVoucherIfAble(updatedInvoice.UserID, updatedInvoice, userSpending)

		s.userSpendingRepository.Update(&models.UserSpending{
			UserID:        updatedInvoice.UserID,
			Month:         currentMonthNum,
			Year:          currentYearNum,
			TotalSpending: userSpending.TotalSpending + updatedInvoice.Total,
		})

		return updatedInvoice, giftsSent, vouchersSent, nil
	}

	return updatedInvoice, giftsSent, vouchersSent, nil
}

func (s *invoiceService) SendVoucherIfAble(userID int64, latestInvoice *models.Invoice, userSpending *models.UserSpending) ([]*models.Voucher, error) {
	vouchersSent := []*models.Voucher{}

	user, err := s.userRepository.GetByID(userID)
	if err != nil {
		return vouchersSent, err
	}

	if user.ReferredUserID == 0 {
		return vouchersSent, nil
	}

	currentTime := time.Now()

	prevTotalPaid := userSpending.TotalSpending
	currentTotalPaid := prevTotalPaid + latestInvoice.Total

	vouchers, err := s.voucherRepository.GetAll()
	if err != nil {
		return vouchersSent, err
	}

	for _, voucher := range vouchers {
		if prevTotalPaid <= voucher.MinimumSpending && currentTotalPaid > voucher.MinimumSpending {
			_, err := s.userVoucherRepository.Insert(&models.UserVoucher{
				UserID:             user.ReferredUserID,
				VoucherID:          int64(voucher.ID),
				Code:               helpers.RandSeq(6),
				ReceivedFromUserID: userID,
				DateReceived:       currentTime,
				ValidUntil:         currentTime.AddDate(0, 1, 0),
			})
			if err == nil {
				vouchersSent = append(vouchersSent, voucher)
			}
		}
	}

	return vouchersSent, err
}

func (s *invoiceService) SendGiftIfAble(userID int64, latestInvoice *models.Invoice, userSpending *models.UserSpending) ([]*models.Gift, error) {
	giftsSent := []*models.Gift{}

	currentTime := time.Now()
	currentMonthNum := int(currentTime.Month())
	currentYearNum := int(currentTime.Year())

	prevTotalPaid := userSpending.TotalSpending
	currentTotalPaid := prevTotalPaid + latestInvoice.Total

	gifts, err := s.giftRepository.GetAll()
	if err != nil {
		return giftsSent, err
	}

	for _, gift := range gifts {
		if prevTotalPaid <= gift.MinimumSpending && currentTotalPaid > gift.MinimumSpending {
			if gift.Stock > 0 {
				_, err := s.userGiftRepository.Insert(&models.UserGift{
					UserID: userID,
					GiftID: int64(gift.ID),
					Month:  currentMonthNum,
					Year:   currentYearNum,
					Status: constants.PROCESSED,
				})

				if err == nil {
					updatedGift, rowsAffected, err := s.giftRepository.UpdateStock(&models.Gift{
						Stock: gift.Stock - 1,
						Model: gorm.Model{
							ID: gift.ID,
						},
					})
					if rowsAffected != 0 && err == nil {
						giftsSent = append(giftsSent, updatedGift)
					}

				}
			}
		}
	}

	return giftsSent, err
}
