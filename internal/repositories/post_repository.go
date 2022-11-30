package repositories

import (
	"fmt"

	"final-project-backend/internal/dtos"
	"final-project-backend/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IPostRepository interface {
	Insert(post *models.Post) (*models.Post, error)
	Update(post *models.Post) (*models.Post, int, error)
	UpdateSingleColumn(key string, post *models.Post) (*models.Post, int, error)
	GetByID(id int64) (*models.Post, error)
	GetAll(query *dtos.PostsRequestQuery) ([]*models.Post, error)
	Delete(id int64) error
	CountByQuery(
		query *dtos.PostsRequestQuery,
	) (int64, error)
	GetTypes() ([]*models.PostType, error)
	GetCategories() ([]*models.Category, error)
	GetTopLikedAndSharedPost(structConditions *models.Post) ([]*models.Post, error)
}

type postRepository struct {
	db *gorm.DB
}

type PostRepositoryConfig struct {
	db *gorm.DB
}

func NewPostRepository(c *PostRepositoryConfig) IPostRepository {
	return &postRepository{
		db: c.db,
	}
}

func (r *postRepository) Insert(post *models.Post) (*models.Post, error) {
	result := r.db.Create(&post)
	if result.Error != nil {
		return nil, result.Error
	}

	return post, nil
}

func (r *postRepository) Update(
	post *models.Post,
) (*models.Post, int, error) {
	result := r.db.Model(&post).Clauses(clause.Returning{}).Updates(&post)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return post, int(result.RowsAffected), nil
}

func (r *postRepository) UpdateSingleColumn(
	key string, post *models.Post,
) (*models.Post, int, error) {
	result := r.db.Model(&post).Select(key).Clauses(clause.Returning{}).Updates(&post)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return post, int(result.RowsAffected), nil
}

func (r *postRepository) GetAll(query *dtos.PostsRequestQuery) ([]*models.Post, error) {
	var posts []*models.Post

	result := r.db.
		Select("posts.id, title, slug, summary, img_thumbnail, author_name, share_count, like_count, created_at").
		Where("posts.title ILIKE ?", "%"+query.Search+"%").
		Where(&models.Post{CategoryID: query.CategoryID, TypeID: query.TypeID}).
		Order(fmt.Sprintf("posts.created_at %s", query.Sort)).
		Offset((query.Page - 1) * query.Limit).
		Limit(query.Limit).
		Joins("Category").
		Joins("Type").
		Find(&posts)

	if result.Error != nil {
		return nil, result.Error
	}

	return posts, nil
}

func (r *postRepository) GetByID(id int64) (*models.Post, error) {
	var post *models.Post

	result := r.db.
		Where("posts.id = ?", id).
		Joins("Category").
		Joins("Type").
		First(&post)

	if result.Error != nil {
		return nil, result.Error
	}

	return post, nil
}

func (r *postRepository) Delete(id int64) error {
	result := r.db.Delete(&models.Post{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *postRepository) CountByQuery(
	query *dtos.PostsRequestQuery,
) (int64, error) {
	var totalRows int64

	result := r.db.Model(&models.Post{}).
		Where("posts.title ILIKE ?", "%"+query.Search+"%").
		Where(&models.Post{CategoryID: query.CategoryID, TypeID: query.TypeID}).
		Count(&totalRows)

	if result.Error != nil {
		return 0, result.Error
	}

	return totalRows, nil
}

func (r *postRepository) GetTypes() ([]*models.PostType, error) {
	var postTypes []*models.PostType

	result := r.db.
		Find(&postTypes)

	if result.Error != nil {
		return nil, result.Error
	}

	return postTypes, nil
}

func (r *postRepository) GetCategories() ([]*models.Category, error) {
	var Categories []*models.Category

	result := r.db.
		Find(&Categories)

	if result.Error != nil {
		return nil, result.Error
	}

	return Categories, nil
}

func (r *postRepository) GetTopLikedAndSharedPost(structConditions *models.Post) ([]*models.Post, error) {
	var posts []*models.Post

	result := r.db.Select("DISTINCT ON (type_id) *").Where(structConditions).Joins("Category").Joins("Type").Order("type_id ASC, like_count + share_count DESC").
		Find(&posts)

	if result.Error != nil {
		return nil, result.Error
	}

	return posts, nil
}
