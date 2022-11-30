package services

import (
	"math"

	"final-project-backend/internal/dtos"
	"final-project-backend/internal/models"
	"final-project-backend/internal/repositories"

	"gorm.io/gorm"
)

type IPostService interface {
	Create(post *models.Post) (*models.Post, error)
	GetByID(id int64) (*models.Post, error)
	GetAll(query *dtos.PostsRequestQuery) ([]*models.Post, error)
	Delete(id int64) error
	CountByQuery(
		query *dtos.PostsRequestQuery,
	) (int64, int64, error)
	GetRecommendedPosts(userID int64) ([]*models.Post, error)
	GetTrendingPosts(userID int64) ([]*models.Post, error)
	GetTypes() ([]*models.PostType, error)
	GetCategories() ([]*models.Category, error)
	Like(postID int64, userID int64, isLike bool) (*models.Post, *models.History, error)
	Share(postID int64, userID int64) (*models.Post, *models.History, error)
}

type postService struct {
	postRepository    repositories.IPostRepository
	historyRepository repositories.IHistoryRepository
}

type PostServiceConfig struct {
	postRepository    repositories.IPostRepository
	historyRepository repositories.IHistoryRepository
}

func NewPostService(c *PostServiceConfig) IPostService {
	return &postService{
		postRepository:    c.postRepository,
		historyRepository: c.historyRepository,
	}
}

func (s *postService) Create(post *models.Post) (*models.Post, error) {
	result, err := s.postRepository.Insert(post)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *postService) GetByID(id int64) (*models.Post, error) {
	result, err := s.postRepository.GetByID(id)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *postService) GetAll(query *dtos.PostsRequestQuery) ([]*models.Post, error) {
	result, err := s.postRepository.GetAll(query)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *postService) Delete(id int64) error {
	err := s.postRepository.Delete(id)
	if err != nil {
		return err
	}

	return nil
}

func (s *postService) CountByQuery(
	query *dtos.PostsRequestQuery,
) (int64, int64, error) {
	totalPosts, err := s.postRepository.CountByQuery(query)
	if err != nil {
		return totalPosts, 0, err
	}

	totalPages := int64(math.Ceil(float64(totalPosts) / float64(query.Limit)))

	return totalPosts, totalPages, nil
}

func (s *postService) GetTypes() ([]*models.PostType, error) {
	result, err := s.postRepository.GetTypes()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *postService) GetCategories() ([]*models.Category, error) {
	result, err := s.postRepository.GetCategories()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *postService) GetRecommendedPosts(userID int64) ([]*models.Post, error) {
	readCounts, err := s.historyRepository.GetCategoriesReadCountPastMonth(&models.History{UserID: userID})

	if err != nil {
		return nil, err
	}

	if len(readCounts) == 0 {
		readCounts, err = s.historyRepository.GetCategoriesReadCountPastMonth(&models.History{})
		if err != nil {
			return nil, err
		}

		if len(readCounts) == 0 {
			return []*models.Post{}, nil
		}
	}

	topCategoryID := readCounts[0].CategoryID
	posts, err := s.postRepository.GetTopLikedAndSharedPost(&models.Post{CategoryID: topCategoryID})
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (s *postService) GetTrendingPosts(userID int64) ([]*models.Post, error) {
	readCounts, err := s.historyRepository.GetCategoriesReadCountPastWeek(&models.History{UserID: userID})

	if err != nil {
		return nil, err
	}

	if len(readCounts) == 0 {
		readCounts, err = s.historyRepository.GetCategoriesReadCountPastWeek(&models.History{})
		if err != nil {
			return nil, err
		}

		if len(readCounts) == 0 {
			return []*models.Post{}, nil
		}
	}

	topCategoryID := readCounts[0].CategoryID
	posts, err := s.historyRepository.GetTrendingPostsByCategory(topCategoryID)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (s *postService) Like(postID int64, userID int64, isLike bool) (*models.Post, *models.History, error) {
	post, err := s.postRepository.GetByID(postID)
	if err != nil {
		return nil, nil, err
	}

	incrementValue := 1
	if !isLike {
		incrementValue = -1
	}

	post, rowsAffected, err := s.postRepository.UpdateSingleColumn("LikeCount", &models.Post{ID: postID, LikeCount: post.LikeCount + incrementValue})
	if err != nil {
		return nil, nil, err
	}

	if rowsAffected == 0 {
		return nil, nil, gorm.ErrRecordNotFound
	}

	updatedHistory, rowsAffected, err := s.historyRepository.UpdateSingleColumn("IsLiked", &models.History{UserID: userID, PostID: postID, IsLiked: isLike})
	if err != nil {
		return nil, nil, err
	}

	if rowsAffected == 0 {
		return nil, nil, gorm.ErrRecordNotFound
	}

	return post, updatedHistory, err
}

func (s *postService) Share(postID int64, userID int64) (*models.Post, *models.History, error) {
	history, err := s.historyRepository.GetByUserAndPostID(userID, postID)
	if err != nil {
		return nil, nil, err
	}

	post, err := s.postRepository.GetByID(postID)
	if err != nil {
		return nil, nil, err
	}

	if history.IsShared {
		return post, history, nil
	}

	updatedHistory, rowsAffected, err := s.historyRepository.UpdateSingleColumn("IsShared", &models.History{UserID: userID, PostID: postID, IsShared: true})
	if err != nil {
		return nil, nil, err
	}

	if rowsAffected == 0 {
		return nil, nil, gorm.ErrRecordNotFound
	}

	post, rowsAffected, err = s.postRepository.UpdateSingleColumn("ShareCount", &models.Post{ID: postID, ShareCount: post.ShareCount + 1})
	if err != nil {
		return nil, nil, err
	}

	if rowsAffected == 0 {
		return nil, nil, gorm.ErrRecordNotFound
	}

	return post, updatedHistory, err
}
