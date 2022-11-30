package repositories

import (
	"final-project-backend/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IUserRepository interface {
	Insert(user *models.User) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByReferralCode(code string) (*models.User, error)
	GetByID(id int64) (*models.User, error)
	GetUsersReferredByID(id int64) ([]*models.User, error)
	GetUsersReferredTotalSpending(id int64) ([]*models.UserSpendingTotalAggregates, error)
	Update(user *models.User) (*models.User, int, error)
}

type userRepository struct {
	db *gorm.DB
}

type UserRepositoryConfig struct {
	db *gorm.DB
}

func NewUserRepository(c *UserRepositoryConfig) IUserRepository {
	return &userRepository{
		db: c.db,
	}
}

func (r *userRepository) Insert(user *models.User) (*models.User, error) {
	result := r.db.Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (r *userRepository) GetByID(id int64) (*models.User, error) {
	var user *models.User

	result := r.db.
		Where("id = ?", id).
		First(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	var user *models.User

	result := r.db.
		Where("email = ?", email).
		First(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (r *userRepository) GetByReferralCode(code string) (*models.User, error) {
	var user *models.User

	result := r.db.
		Where("referral_code = ?", code).
		First(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (r *userRepository) GetUsersReferredByID(id int64) ([]*models.User, error) {
	var users []*models.User

	result := r.db.Where("users.referred_user_id = ?", id).Find(&users)

	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}

func (r *userRepository) GetUsersReferredTotalSpending(id int64) ([]*models.UserSpendingTotalAggregates, error) {
	var users []*models.UserSpendingTotalAggregates = []*models.UserSpendingTotalAggregates{}

	subQuery := r.db.Table("user_spendings").Select("user_id, SUM(total_spending) AS total_spending").Group("user_id")

	result := r.db.Table("users").Joins("JOIN (?) as t1 ON users.id = t1.user_id", subQuery).Where("users.referred_user_id = ?", id).Scan(&users)

	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}

func (r *userRepository) Update(user *models.User) (*models.User, int, error) {
	result := r.db.Model(&user).Clauses(clause.Returning{}).Updates(user)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return user, int(result.RowsAffected), nil
}
