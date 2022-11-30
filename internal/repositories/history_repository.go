package repositories

import (
	"time"

	"final-project-backend/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IHistoryRepository interface {
	Insert(history *models.History) (*models.History, error)
	Update(history *models.History) (*models.History, int, error)
	UpdateSingleColumn(
		key string, history *models.History,
	) (*models.History, int, error)
	GetByUserID(userID int64) ([]*models.History, error)
	GetByUserAndPostID(userID int64, postID int64) (*models.History, error)
	GetCategoriesReadCountPastMonth(structConditions *models.History) ([]*models.PostCategoryReadCount, error)
	GetCategoriesReadCountPastWeek(structConditions *models.History) ([]*models.PostCategoryReadCount, error)
	GetTrendingPostsByCategory(categoryID int64) ([]*models.Post, error)
}

type historyRepository struct {
	db *gorm.DB
}

type HistoryRepositoryConfig struct {
	db *gorm.DB
}

func NewHistoryRepository(c *HistoryRepositoryConfig) IHistoryRepository {
	return &historyRepository{
		db: c.db,
	}
}

func (r *historyRepository) Insert(history *models.History) (*models.History, error) {
	result := r.db.Create(&history)
	if result.Error != nil {
		return nil, result.Error
	}

	return history, nil
}

func (r *historyRepository) Update(
	history *models.History,
) (*models.History, int, error) {
	result := r.db.Model(&history).Clauses(clause.Returning{}).Updates(history)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return history, int(result.RowsAffected), nil
}

func (r *historyRepository) UpdateSingleColumn(
	key string, history *models.History,
) (*models.History, int, error) {
	result := r.db.Model(&history).Select(key).Clauses(clause.Returning{}).Updates(&history)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return history, int(result.RowsAffected), nil
}

func (r *historyRepository) GetByUserID(userID int64) ([]*models.History, error) {
	var histories []*models.History
	result := r.db.Preload("Post.Category").Preload("Post.Type").Where(&models.History{UserID: userID}).Joins("Post").Find(&histories)

	if result.Error != nil {
		return nil, result.Error
	}

	return histories, nil
}

func (r *historyRepository) GetByUserAndPostID(userID int64, postID int64) (*models.History, error) {
	var history *models.History
	result := r.db.Where(&models.History{UserID: userID, PostID: postID}).Joins("Post").First(&history)

	if result.Error != nil {
		return nil, result.Error
	}

	return history, nil
}

func (r *historyRepository) GetCategoriesReadCountPastMonth(structConditions *models.History) ([]*models.PostCategoryReadCount, error) {
	var readCounts []*models.PostCategoryReadCount

	currentMonthNum := int(time.Now().Month())

	result := r.db.Model(&models.History{}).
		Select("COUNT(post_id) AS read_count, posts.category_id").
		Joins("LEFT JOIN posts ON posts.id = histories.post_id").
		Where(structConditions).
		Where("EXTRACT(MONTH FROM histories.last_accessed AT TIME ZONE '+7') = ? OR EXTRACT(MONTH FROM histories.last_accessed AT TIME ZONE '+7') = ?", currentMonthNum, currentMonthNum-1).
		Group("posts.category_id").
		Order("read_count DESC").
		Scan(&readCounts)

	if result.Error != nil {
		return nil, result.Error
	}

	return readCounts, nil
}

func (r *historyRepository) GetCategoriesReadCountPastWeek(structConditions *models.History) ([]*models.PostCategoryReadCount, error) {
	var readCounts []*models.PostCategoryReadCount

	result := r.db.Model(&models.History{}).
		Select("COUNT(post_id) AS read_count, posts.category_id").
		Joins("LEFT JOIN posts ON posts.id = histories.post_id").
		Where(structConditions).
		Where("histories.last_accessed BETWEEN NOW() - INTERVAL '7 DAY' AND NOW() ").
		Group("posts.category_id").
		Order("read_count DESC").
		Scan(&readCounts)

	if result.Error != nil {
		return nil, result.Error
	}

	return readCounts, nil
}

func (r *historyRepository) GetTrendingPostsByCategory(categoryID int64) ([]*models.Post, error) {
	var trendingPosts []*models.Post

	subQuery := r.db.Debug().Table("histories").Select("COUNT(user_id) AS reader_count, post_id").Joins("LEFT JOIN posts ON posts.id = histories.post_id").Where("posts.category_id = ?", categoryID).Group("post_id")

	result := r.db.Debug().Joins("JOIN (?) as t1 ON posts.id = t1.post_id", subQuery).Joins("Category").Joins("Type").Order("t1.reader_count DESC").Limit(5).Find(&trendingPosts)

	if result.Error != nil {
		return nil, result.Error
	}

	return trendingPosts, nil

}
