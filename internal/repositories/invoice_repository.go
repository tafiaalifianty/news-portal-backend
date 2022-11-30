package repositories

import (
	"strconv"
	"time"

	"final-project-backend/internal/constants"
	"final-project-backend/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IInvoiceRepository interface {
	GetAll(*models.Invoice) ([]*models.Invoice, error)
	GetAllCurrentMonth(invoice *models.Invoice) ([]*models.Invoice, error)
	Insert(invoice *models.Invoice) (*models.Invoice, error)
	GetByID(id int64) (*models.Invoice, error)
	GetByCode(code string) (*models.Invoice, error)
	Update(invoice *models.Invoice) (*models.Invoice, int, error)
}

type invoiceRepository struct {
	db *gorm.DB
}

type InvoiceRepositoryConfig struct {
	db *gorm.DB
}

func NewInvoiceRepository(c *InvoiceRepositoryConfig) IInvoiceRepository {
	return &invoiceRepository{
		db: c.db,
	}
}

func (r *invoiceRepository) Insert(invoice *models.Invoice) (*models.Invoice, error) {
	result := r.db.Create(&invoice)
	if result.Error != nil {
		return nil, result.Error
	}

	invoiceNumber := strconv.Itoa(int(invoice.ID) + constants.INVOICE_STARTING_NUMBER)
	subscriptionID := strconv.Itoa(int(invoice.SubscriptionID))

	result = r.db.Model(&invoice).
		Update("code", constants.INVOICE_PREFIX+subscriptionID+"-"+invoiceNumber)

	if result.Error != nil {
		return nil, result.Error
	}

	return invoice, nil
}

func (r *invoiceRepository) GetAll(invoice *models.Invoice) ([]*models.Invoice, error) {
	var invoices []*models.Invoice
	result := r.db.Where(invoice).Order("created_at desc").Joins("User").Joins("Subscription").Find(&invoices)

	if result.Error != nil {
		return nil, result.Error
	}

	return invoices, nil
}

func (r *invoiceRepository) GetAllCurrentMonth(invoice *models.Invoice) ([]*models.Invoice, error) {
	var invoices []*models.Invoice

	currentMonthNum := int(time.Now().Month())

	result := r.db.
		Where(invoice).
		Where("EXTRACT(MONTH FROM invoices.paid_at AT TIME ZONE '+7') = ?", currentMonthNum).
		Order("created_at asc").
		Joins("User").
		Joins("Subscription").
		Find(&invoices)

	if result.Error != nil {
		return nil, result.Error
	}

	return invoices, nil
}

func (r *invoiceRepository) GetByID(id int64) (*models.Invoice, error) {
	var invoice *models.Invoice

	result := r.db.
		Where("invoices.id = ?", id).
		Joins("User").
		Joins("Subscription").
		First(&invoice)

	if result.Error != nil {
		return nil, result.Error
	}

	return invoice, nil
}

func (r *invoiceRepository) GetByCode(code string) (*models.Invoice, error) {
	var invoice *models.Invoice

	result := r.db.
		Where("invoices.code = ?", code).
		Joins("User").
		Joins("Subscription").
		First(&invoice)

	if result.Error != nil {
		return nil, result.Error
	}

	return invoice, nil
}

func (r *invoiceRepository) Update(invoice *models.Invoice) (*models.Invoice, int, error) {
	result := r.db.Model(&invoice).Clauses(clause.Returning{}).Updates(invoice)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return invoice, int(result.RowsAffected), nil
}
