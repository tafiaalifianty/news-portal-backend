package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"final-project-backend/internal/dtos"
	errn "final-project-backend/internal/errors"
	"final-project-backend/internal/helpers"
	"final-project-backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgconn"
	"gorm.io/gorm"
)

func (h *Handler) CreatePost(c *gin.Context) {
	var request dtos.CreatePostRequest
	var response dtos.CreatePostResponse

	err := c.ShouldBindJSON(&request)
	if err != nil {
		helpers.SendErrorResponse(c, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	// TODO: Resize image_url for image_thumbnail
	post := &models.Post{
		Title:        request.Title,
		Slug:         strings.ReplaceAll(request.Title, " ", "-"),
		Content:      request.Content,
		Summary:      request.Summary,
		CategoryID:   request.CategoryID,
		TypeID:       request.TypeID,
		ImgUrl:       request.ImgUrl,
		ImgThumbnail: request.ImgUrl,
		AuthorName:   request.AuthorName,
	}

	createdPost, err := h.services.Post.Create(post)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == errn.ForeignKeyViolation {
			helpers.SendErrorResponse(c, http.StatusBadRequest, errn.ErrInvalidFields.Error())

			return
		}
		helpers.SendErrorResponse(c, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))

		return
	}

	response = dtos.CreatePostResponse{
		ID: createdPost.ID,
	}

	helpers.SendSuccessResponse(c, http.StatusCreated, http.StatusText(http.StatusCreated), response)
}

func (h *Handler) GetAllPosts(c *gin.Context) {
	var response dtos.GetAllPostResponse

	search := c.DefaultQuery("s", "")
	sort := c.DefaultQuery("sort", "desc")
	category, err1 := strconv.Atoi(c.DefaultQuery("category", "0"))
	postType, err2 := strconv.Atoi(c.DefaultQuery("type", "0"))
	limit, err3 := strconv.Atoi(c.DefaultQuery("limit", "10"))
	page, err4 := strconv.Atoi(c.DefaultQuery("page", "1"))

	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		helpers.SendErrorResponse(c, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	if limit < 1 || page < 1 {
		helpers.SendErrorResponse(c, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	query := &dtos.PostsRequestQuery{
		CategoryID: int64(category),
		TypeID:     int64(postType),
		Search:     search,
		Sort:       sort,
		Limit:      limit,
		Page:       page,
	}

	posts, err := h.services.Post.GetAll(query)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			helpers.SendErrorResponse(c, http.StatusNotFound, errn.ErrNoPostsFound.Error())
			return
		}

		helpers.SendErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	totalRows, totalPages, err := h.services.Post.CountByQuery(query)
	if err != nil {
		helpers.SendErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	formattedPosts := dtos.FormatPostsCompact(posts)
	response = dtos.GetAllPostResponse{
		Data: formattedPosts,
		PaginationResponse: dtos.PaginationResponse{
			PerPage:     limit,
			CurrentPage: page,
			TotalRows:   totalRows,
			TotalPages:  totalPages,
		},
	}

	helpers.SendSuccessResponse(c, http.StatusOK, http.StatusText(http.StatusOK), response)
}

func (h *Handler) GetPostByID(c *gin.Context) {
	userContext, ok := c.Get("user")
	if !ok {
		helpers.SendErrorResponse(
			c,
			http.StatusBadRequest,
			http.StatusText(http.StatusBadRequest),
		)
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		helpers.SendErrorResponse(c, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))

		return
	}

	post, err := h.services.Post.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			helpers.SendErrorResponse(c, http.StatusNotFound, errn.ErrNoPostsFound.Error())
			return
		}

		helpers.SendErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	user := userContext.(dtos.JwtData)

	if user.Role == string(models.Member) {
		_, err := h.services.History.GetByUserAndPostID(user.ID, post.ID)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = h.services.UserSubscription.ValidateUserQuota(user.ID, post.Type.Quota)
		}

		if err != nil {
			if errors.Is(err, errn.ErrNotEnoughQuota) {
				helpers.SendErrorResponse(c, http.StatusBadRequest, err.Error())
				return
			}

			helpers.SendErrorResponse(c, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		history := &models.History{
			UserID: userContext.(dtos.JwtData).ID,
			PostID: post.ID,
		}

		updatedHistory, err := h.services.History.UpdateOrInsert(history)
		if err != nil {
			helpers.SendErrorResponse(c, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))

			return
		}

		response := dtos.GetPostByIDWithHistoryResponse{
			PostResponse:       *dtos.FormatPost(post),
			HistoryResponseDTO: *dtos.FormatHistory(updatedHistory),
		}

		response.HistoryResponseDTO.Post = nil

		helpers.SendSuccessResponse(c, http.StatusOK, http.StatusText(http.StatusOK), response)
		return
	}

	response := dtos.GetPostByIDWithHistoryResponse{
		PostResponse: *dtos.FormatPost(post),
	}

	helpers.SendSuccessResponse(c, http.StatusOK, http.StatusText(http.StatusOK), response)
}

func (h *Handler) LikePost(c *gin.Context) {
	var request dtos.LikePostRequest
	var response dtos.GetPostByIDWithHistoryResponse

	userContext, ok := c.Get("user")
	if !ok {
		helpers.SendErrorResponse(
			c,
			http.StatusBadRequest,
			http.StatusText(http.StatusBadRequest),
		)
		return
	}
	user := userContext.(dtos.JwtData)

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		helpers.SendErrorResponse(c, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))

		return
	}

	err = c.ShouldBindJSON(&request)
	if err != nil {
		helpers.SendErrorResponse(c, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	post, history, err := h.services.Post.Like(id, user.ID, *request.IsLike)
	if err != nil {
		helpers.SendErrorResponse(c, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	response = dtos.GetPostByIDWithHistoryResponse{
		PostResponse:       *dtos.FormatPost(post),
		HistoryResponseDTO: *dtos.FormatHistory(history),
	}

	helpers.SendSuccessResponse(c, http.StatusOK, http.StatusText(http.StatusOK), response)
}

func (h *Handler) SharePost(c *gin.Context) {
	var response dtos.GetPostByIDWithHistoryResponse

	userContext, ok := c.Get("user")
	if !ok {
		helpers.SendErrorResponse(
			c,
			http.StatusBadRequest,
			http.StatusText(http.StatusBadRequest),
		)
		return
	}
	user := userContext.(dtos.JwtData)

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		helpers.SendErrorResponse(c, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))

		return
	}

	post, history, err := h.services.Post.Share(id, user.ID)
	if err != nil {
		helpers.SendErrorResponse(c, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	response = dtos.GetPostByIDWithHistoryResponse{
		PostResponse:       *dtos.FormatPost(post),
		HistoryResponseDTO: *dtos.FormatHistory(history),
	}

	helpers.SendSuccessResponse(c, http.StatusOK, http.StatusText(http.StatusOK), response)
}

func (h *Handler) DeletePost(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		helpers.SendErrorResponse(c, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))

		return
	}

	err = h.services.Post.Delete(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			helpers.SendErrorResponse(c, http.StatusNotFound, errn.ErrNoPostsFound.Error())
			return
		}

		helpers.SendErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	helpers.SendSuccessResponse(c, http.StatusOK, http.StatusText(http.StatusOK), id)
}

func (h *Handler) GetAllTypes(c *gin.Context) {
	var response dtos.GetAllTypesResponse

	types, err := h.services.Post.GetTypes()
	if err != nil {
		helpers.SendErrorResponse(c, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))

		return
	}

	response = types

	helpers.SendSuccessResponse(c, http.StatusOK, http.StatusText(http.StatusOK), response)
}

func (h *Handler) GetAllCategories(c *gin.Context) {
	var response dtos.GetAllCategoriesResponse

	categories, err := h.services.Post.GetCategories()
	if err != nil {
		helpers.SendErrorResponse(c, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))

		return
	}

	response = categories

	helpers.SendSuccessResponse(c, http.StatusOK, http.StatusText(http.StatusOK), response)
}

func (h *Handler) GetRecommendedPost(c *gin.Context) {
	var response []*dtos.PostResponseCompact

	userContext, ok := c.Get("user")
	if !ok {
		helpers.SendErrorResponse(
			c,
			http.StatusBadRequest,
			http.StatusText(http.StatusBadRequest),
		)
		return
	}

	posts, err := h.services.Post.GetRecommendedPosts(userContext.(dtos.JwtData).ID)
	if err != nil {
		helpers.SendErrorResponse(c, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))

		return
	}

	response = dtos.FormatPostsCompact(posts)

	helpers.SendSuccessResponse(c, http.StatusOK, http.StatusText(http.StatusOK), response)
}

func (h *Handler) GetTrendingPosts(c *gin.Context) {
	var response []*dtos.PostResponseCompact

	userContext, ok := c.Get("user")
	if !ok {
		helpers.SendErrorResponse(
			c,
			http.StatusBadRequest,
			http.StatusText(http.StatusBadRequest),
		)
		return
	}

	posts, err := h.services.Post.GetTrendingPosts(userContext.(dtos.JwtData).ID)
	if err != nil {
		helpers.SendErrorResponse(c, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))

		return
	}

	response = dtos.FormatPostsCompact(posts)

	helpers.SendSuccessResponse(c, http.StatusOK, http.StatusText(http.StatusOK), response)
}
