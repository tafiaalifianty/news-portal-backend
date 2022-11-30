package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"final-project-backend/internal/dtos"
	errn "final-project-backend/internal/errors"
	"final-project-backend/internal/helpers"
	"final-project-backend/internal/models"
	"final-project-backend/internal/services"
	"final-project-backend/mocks"

	"github.com/jackc/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestHandler_CreatePost(t *testing.T) {
	invalidRequest := &dtos.CreatePostRequest{}
	validRequest := &dtos.CreatePostRequest{
		Title:      "Title",
		Content:    "Content",
		Summary:    "Summary",
		CategoryID: 1,
		TypeID:     1,
		ImgUrl:     "ImgUrl",
		AuthorName: "AuthorName",
	}
	formattedPostDto := &models.Post{
		Title:        validRequest.Title,
		Slug:         strings.ReplaceAll(validRequest.Title, " ", "-"),
		Content:      validRequest.Content,
		Summary:      validRequest.Summary,
		CategoryID:   validRequest.CategoryID,
		TypeID:       validRequest.TypeID,
		ImgUrl:       validRequest.ImgUrl,
		ImgThumbnail: validRequest.ImgUrl,
		AuthorName:   validRequest.AuthorName,
	}
	mockPostReturn := &models.Post{ID: 1}

	mockValidDataInInterface, err := helpers.StructToMap(&dtos.CreatePostResponse{
		ID: mockPostReturn.ID,
	})
	require.NoError(t, err)

	mockError := fmt.Errorf("error")
	type fields struct {
		postService *mocks.IPostService
	}
	type args struct {
		body io.Reader
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		mock   func(*mocks.IPostService)
		want   helpers.JsonResponse
	}{
		{
			name: "ERROR | Invalid Request Body",
			fields: fields{
				postService: mocks.NewIPostService(t),
			},
			args: args{
				body: helpers.MakeRequestBody(invalidRequest),
			},
			mock: func(s *mocks.IPostService) {
			},
			want: helpers.JsonResponse{
				Code:    http.StatusBadRequest,
				Message: http.StatusText(http.StatusBadRequest),
				Data:    nil,
				IsError: true,
			},
		},
		{
			name: "ERROR | Error from Services.Post.Create (Foreign Key Violation: invalid categoryID or typeID)",
			fields: fields{
				postService: mocks.NewIPostService(t),
			},
			args: args{
				body: helpers.MakeRequestBody(validRequest),
			},
			mock: func(s *mocks.IPostService) {
				s.On("Create", formattedPostDto).Return(nil, &pgconn.PgError{
					Code: errn.ForeignKeyViolation,
				})
			},
			want: helpers.JsonResponse{
				Code:    http.StatusBadRequest,
				Message: errn.ErrInvalidFields.Error(),
				Data:    nil,
				IsError: true,
			},
		},
		{
			name: "ERROR | Error from Services.Post.Create",
			fields: fields{
				postService: mocks.NewIPostService(t),
			},
			args: args{
				body: helpers.MakeRequestBody(validRequest),
			},
			mock: func(s *mocks.IPostService) {
				s.On("Create", formattedPostDto).Return(nil, mockError)
			},
			want: helpers.JsonResponse{
				Code:    http.StatusInternalServerError,
				Message: http.StatusText(http.StatusInternalServerError),
				Data:    nil,
				IsError: true,
			},
		},
		{
			name: "SUCCESS",
			fields: fields{
				postService: mocks.NewIPostService(t),
			},
			args: args{
				body: helpers.MakeRequestBody(validRequest),
			},
			mock: func(s *mocks.IPostService) {
				s.On("Create", formattedPostDto).Return(mockPostReturn, nil)
			},
			want: helpers.JsonResponse{
				Code:    http.StatusCreated,
				Message: http.StatusText(http.StatusCreated),
				Data:    mockValidDataInInterface,
				IsError: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				services: &services.Services{
					Post: tt.fields.postService,
				},
			}

			tt.mock(tt.fields.postService)
			r := helpers.SetUpRouter()
			endpoint := "/admin/posts"
			r.POST(endpoint, h.CreatePost)
			req, _ := http.NewRequest(
				http.MethodPost,
				endpoint,
				tt.args.body,
			)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			var response helpers.JsonResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			assert.Equal(t, tt.want.Code, w.Code)
			assert.Equal(t, tt.want, response)
		})
	}
}

func TestHandler_GetAllPosts(t *testing.T) {
	endpoint := "/posts"
	defaultQuery := &dtos.PostsRequestQuery{Sort: "desc", Limit: 10, Page: 1}
	mockPosts := []*models.Post{}
	validResponse := &dtos.GetAllPostResponse{
		Data: dtos.FormatPostsCompact(mockPosts),
		PaginationResponse: dtos.PaginationResponse{
			PerPage:     defaultQuery.Limit,
			CurrentPage: defaultQuery.Page,
			TotalRows:   1,
			TotalPages:  1,
		},
	}
	mockValidDataInInterface, err := helpers.StructToMap(validResponse)
	require.NoError(t, err)

	mockError := fmt.Errorf("error")

	type fields struct {
		postService mocks.IPostService
	}
	tests := []struct {
		name   string
		fields fields
		query  string
		mock   func(*mocks.IPostService)
		want   helpers.JsonResponse
	}{
		{
			name:  "ERROR | Error from category query provided is NaN",
			query: "?category=a",
			fields: fields{
				postService: *mocks.NewIPostService(t),
			},
			mock: func(s *mocks.IPostService) {
			},
			want: helpers.JsonResponse{
				Code:    http.StatusBadRequest,
				Message: http.StatusText(http.StatusBadRequest),
				Data:    nil,
				IsError: true,
			},
		},
		{
			name:  "ERROR | Error from category postType provided is NaN",
			query: "?type=a",
			fields: fields{
				postService: *mocks.NewIPostService(t),
			},
			mock: func(s *mocks.IPostService) {
			},
			want: helpers.JsonResponse{
				Code:    http.StatusBadRequest,
				Message: http.StatusText(http.StatusBadRequest),
				Data:    nil,
				IsError: true,
			},
		},
		{
			name:  "ERROR | Error from limit query provided is NaN",
			query: "?limit=a",
			fields: fields{
				postService: *mocks.NewIPostService(t),
			},
			mock: func(s *mocks.IPostService) {
			},
			want: helpers.JsonResponse{
				Code:    http.StatusBadRequest,
				Message: http.StatusText(http.StatusBadRequest),
				Data:    nil,
				IsError: true,
			},
		},
		{
			name:  "ERROR | Error from page query provided is NaN",
			query: "?page=a",
			fields: fields{
				postService: *mocks.NewIPostService(t),
			},
			mock: func(s *mocks.IPostService) {
			},
			want: helpers.JsonResponse{
				Code:    http.StatusBadRequest,
				Message: http.StatusText(http.StatusBadRequest),
				Data:    nil,
				IsError: true,
			},
		},
		{
			name:  "ERROR | Error from limit query provided is less than 1",
			query: "?limit=0",
			fields: fields{
				postService: *mocks.NewIPostService(t),
			},
			mock: func(s *mocks.IPostService) {
			},
			want: helpers.JsonResponse{
				Code:    http.StatusBadRequest,
				Message: http.StatusText(http.StatusBadRequest),
				Data:    nil,
				IsError: true,
			},
		},
		{
			name:  "ERROR | Error from page query provided is less than 1",
			query: "?page=0",
			fields: fields{
				postService: *mocks.NewIPostService(t),
			},
			mock: func(s *mocks.IPostService) {
			},
			want: helpers.JsonResponse{
				Code:    http.StatusBadRequest,
				Message: http.StatusText(http.StatusBadRequest),
				Data:    nil,
				IsError: true,
			},
		},
		{
			name: "ERROR | Error from PostService.GetAll (posts not found)",
			fields: fields{
				postService: *mocks.NewIPostService(t),
			},
			mock: func(s *mocks.IPostService) {
				s.On("GetAll", defaultQuery).Return(nil, gorm.ErrRecordNotFound)
			},
			want: helpers.JsonResponse{
				Code:    http.StatusNotFound,
				Message: errn.ErrNoPostsFound.Error(),
				Data:    nil,
				IsError: true,
			},
		},
		{
			name: "ERROR | Error from PostService.GetAll (other errors)",
			fields: fields{
				postService: *mocks.NewIPostService(t),
			},
			mock: func(s *mocks.IPostService) {
				s.On("GetAll", defaultQuery).Return(nil, mockError)
			},
			want: helpers.JsonResponse{
				Code:    http.StatusInternalServerError,
				Message: mockError.Error(),
				Data:    nil,
				IsError: true,
			},
		},
		{
			name: "ERROR | Error from PostService.CountByQuery",
			fields: fields{
				postService: *mocks.NewIPostService(t),
			},
			mock: func(s *mocks.IPostService) {
				s.On("GetAll", defaultQuery).Return(mockPosts, nil)
				s.On("CountByQuery", defaultQuery).Return(int64(0), int64(0), mockError)
			},
			want: helpers.JsonResponse{
				Code:    http.StatusInternalServerError,
				Message: mockError.Error(),
				Data:    nil,
				IsError: true,
			},
		},
		{
			name: "SUCCESS",
			fields: fields{
				postService: *mocks.NewIPostService(t),
			},
			mock: func(s *mocks.IPostService) {
				s.On("GetAll", defaultQuery).Return(mockPosts, nil)
				s.On("CountByQuery", defaultQuery).Return(int64(1), int64(1), nil)
			},
			want: helpers.JsonResponse{
				Code:    http.StatusOK,
				Message: http.StatusText(http.StatusOK),
				Data:    mockValidDataInInterface,
				IsError: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				services: &services.Services{
					Post: &tt.fields.postService,
				},
			}

			tt.mock(&tt.fields.postService)
			r := helpers.SetUpRouter()

			r.GET(endpoint, h.GetAllPosts)

			req, _ := http.NewRequest(
				http.MethodGet,
				endpoint+tt.query,
				nil,
			)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			var response helpers.JsonResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			assert.Equal(t, tt.want.Code, w.Code)
			assert.Equal(t, tt.want, response)
		})
	}
}

func TestHandler_GetPostByID(t *testing.T) {
	mockPost := &models.Post{}
	mockJwtUserID := 1
	validResponse := &dtos.GetPostByIDWithHistoryResponse{}
	mockValidDataInInterface, err := helpers.StructToMap(validResponse)
	require.NoError(t, err)
	mockError := fmt.Errorf("error")
	mockIDParams := "1"
	var mockID int64 = 1

	type fields struct {
		postService             mocks.IPostService
		historyService          mocks.IHistoryService
		userSubscriptionService mocks.IUserSubscriptionService
	}
	tests := []struct {
		name                     string
		fields                   fields
		mock                     func(*mocks.IPostService, *mocks.IHistoryService, *mocks.IUserSubscriptionService)
		mockParamsFromMiddleware bool
		mockedUserRole           string
		want                     helpers.JsonResponse
	}{
		{
			name: "ERROR | Error from Context.Get",
			fields: fields{
				postService: *mocks.NewIPostService(t),
			},
			mock: func(as *mocks.IPostService, h *mocks.IHistoryService, us *mocks.IUserSubscriptionService) {
			},
			mockParamsFromMiddleware: false,
			mockedUserRole:           "",
			want: helpers.JsonResponse{
				Code:    http.StatusBadRequest,
				Message: http.StatusText(http.StatusBadRequest),
				Data:    nil,
				IsError: true,
			},
		},
		{
			name: "ERROR | Error from invalid params",
			fields: fields{
				postService: *mocks.NewIPostService(t),
			},
			mock: func(as *mocks.IPostService, h *mocks.IHistoryService, us *mocks.IUserSubscriptionService) {
			},
			mockParamsFromMiddleware: false,
			mockedUserRole:           "member",
			want: helpers.JsonResponse{
				Code:    http.StatusBadRequest,
				Message: http.StatusText(http.StatusBadRequest),
				Data:    nil,
				IsError: true,
			},
		},
		{
			name: "ERROR | Error from Post.GetByID (no posts found)",
			fields: fields{
				postService: *mocks.NewIPostService(t),
			},
			mock: func(s *mocks.IPostService, h *mocks.IHistoryService, us *mocks.IUserSubscriptionService) {
				s.On("GetByID", mockID).Return(nil, gorm.ErrRecordNotFound)
			},
			mockParamsFromMiddleware: true,
			mockedUserRole:           "member",
			want: helpers.JsonResponse{
				Code:    http.StatusNotFound,
				Message: errn.ErrNoPostsFound.Error(),
				Data:    nil,
				IsError: true,
			},
		},
		{
			name: "ERROR | Error from Post.GetByID (other errors)",
			fields: fields{
				postService: *mocks.NewIPostService(t),
			},
			mock: func(s *mocks.IPostService, h *mocks.IHistoryService, us *mocks.IUserSubscriptionService) {
				s.On("GetByID", mockID).Return(nil, mockError)
			},
			mockParamsFromMiddleware: true,
			mockedUserRole:           "member",
			want: helpers.JsonResponse{
				Code:    http.StatusInternalServerError,
				Message: mockError.Error(),
				Data:    nil,
				IsError: true,
			},
		},
		{
			name: "ERROR | when user is member, Error from History.GetByUserAndPostID and not gorm.ErrRecordNotFound",
			fields: fields{
				postService: *mocks.NewIPostService(t),
			},
			mock: func(s *mocks.IPostService, h *mocks.IHistoryService, us *mocks.IUserSubscriptionService) {
				s.On("GetByID", mockID).Return(mockPost, nil)
				h.On("GetByUserAndPostID", int64(mockJwtUserID), mockPost.ID).Return(nil, mockError)
			},
			mockParamsFromMiddleware: true,
			mockedUserRole:           "member",
			want: helpers.JsonResponse{
				Code:    http.StatusInternalServerError,
				Message: http.StatusText(http.StatusInternalServerError),
				Data:    nil,
				IsError: true,
			},
		},
		{
			name: "ERROR | when user is member, Error NotEnoughQuota from UserSubscription.ValidateUserQuota",
			fields: fields{
				postService: *mocks.NewIPostService(t),
			},
			mock: func(s *mocks.IPostService, h *mocks.IHistoryService, us *mocks.IUserSubscriptionService) {
				s.On("GetByID", mockID).Return(mockPost, nil)
				h.On("GetByUserAndPostID", int64(mockJwtUserID), mockPost.ID).Return(nil, gorm.ErrRecordNotFound)
				us.On("ValidateUserQuota", int64(mockJwtUserID), mockPost.Type.Quota).Return(errn.ErrNotEnoughQuota)
			},
			mockParamsFromMiddleware: true,
			mockedUserRole:           "member",
			want: helpers.JsonResponse{
				Code:    http.StatusBadRequest,
				Message: errn.ErrNotEnoughQuota.Error(),
				Data:    nil,
				IsError: true,
			},
		},
		{
			name: "ERROR | when user is member, Error others from UserSubscription.ValidateUserQuota",
			fields: fields{
				postService: *mocks.NewIPostService(t),
			},
			mock: func(s *mocks.IPostService, h *mocks.IHistoryService, us *mocks.IUserSubscriptionService) {
				s.On("GetByID", mockID).Return(mockPost, nil)
				h.On("GetByUserAndPostID", int64(mockJwtUserID), mockPost.ID).Return(nil, gorm.ErrRecordNotFound)
				us.On("ValidateUserQuota", int64(mockJwtUserID), mockPost.Type.Quota).Return(mockError)
			},
			mockParamsFromMiddleware: true,
			mockedUserRole:           "member",
			want: helpers.JsonResponse{
				Code:    http.StatusInternalServerError,
				Message: http.StatusText(http.StatusInternalServerError),
				Data:    nil,
				IsError: true,
			},
		},
		{
			name: "ERROR | Error from History.UpdateOrInsert when user is member",
			fields: fields{
				postService: *mocks.NewIPostService(t),
			},
			mock: func(s *mocks.IPostService, h *mocks.IHistoryService, us *mocks.IUserSubscriptionService) {
				s.On("GetByID", mockID).Return(mockPost, nil)
				h.On("GetByUserAndPostID", int64(mockJwtUserID), mockPost.ID).Return(nil, gorm.ErrRecordNotFound)
				us.On("ValidateUserQuota", int64(mockJwtUserID), mockPost.Type.Quota).Return(nil)
				h.On("UpdateOrInsert", &models.History{
					UserID: int64(mockJwtUserID),
					PostID: mockPost.ID,
				}).Return(nil, mockError)
			},
			mockParamsFromMiddleware: true,
			mockedUserRole:           "member",
			want: helpers.JsonResponse{
				Code:    http.StatusInternalServerError,
				Message: http.StatusText(http.StatusInternalServerError),
				Data:    nil,
				IsError: true,
			},
		},
		{
			name: "SUCCESS",
			fields: fields{
				postService: *mocks.NewIPostService(t),
			},
			mock: func(s *mocks.IPostService, h *mocks.IHistoryService, us *mocks.IUserSubscriptionService) {
				s.On("GetByID", mockID).Return(mockPost, nil)
				h.On("GetByUserAndPostID", int64(mockJwtUserID), mockPost.ID).Return(nil, gorm.ErrRecordNotFound)
				us.On("ValidateUserQuota", int64(mockJwtUserID), mockPost.Type.Quota).Return(nil, nil)
				h.On("UpdateOrInsert", &models.History{
					UserID: int64(mockJwtUserID),
					PostID: mockPost.ID,
				}).Return(nil, nil)
			},
			mockParamsFromMiddleware: true,
			mockedUserRole:           "admin",
			want: helpers.JsonResponse{
				Code:    http.StatusOK,
				Message: http.StatusText(http.StatusOK),
				Data:    mockValidDataInInterface,
				IsError: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				services: &services.Services{
					Post:             &tt.fields.postService,
					History:          &tt.fields.historyService,
					UserSubscription: &tt.fields.userSubscriptionService,
				},
			}

			tt.mock(&tt.fields.postService, &tt.fields.historyService, &tt.fields.userSubscriptionService)
			r := helpers.SetUpRouter()
			endpoint := "/transactions/1"

			if tt.mockedUserRole == "" {
				r.GET(endpoint, h.GetPostByID)
			} else {
				if !tt.mockParamsFromMiddleware {
					r.GET(endpoint, helpers.MiddlewareMockUser(dtos.JwtData{Role: tt.mockedUserRole, ID: int64(mockJwtUserID)}), h.GetPostByID)
				} else {
					r.GET(endpoint, helpers.MiddlewareMockID(mockIDParams), helpers.MiddlewareMockUser(dtos.JwtData{Role: tt.mockedUserRole, ID: int64(mockJwtUserID)}), h.GetPostByID)
				}
			}

			req, _ := http.NewRequest(
				http.MethodGet,
				endpoint,
				nil,
			)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			var response helpers.JsonResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			assert.Equal(t, tt.want.Code, w.Code)
			assert.Equal(t, tt.want, response)
		})
	}
}

func TestHandler_DeletePost(t *testing.T) {
	mockError := fmt.Errorf("error")
	mockIDParams := "1"
	var mockID int64 = 1

	type fields struct {
		postService mocks.IPostService
	}
	tests := []struct {
		name                     string
		fields                   fields
		mock                     func(*mocks.IPostService)
		mockParamsFromMiddleware bool
		want                     helpers.JsonResponse
	}{
		{
			name: "ERROR | Error from invalid params",
			fields: fields{
				postService: *mocks.NewIPostService(t),
			},
			mock: func(as *mocks.IPostService) {
			},
			mockParamsFromMiddleware: false,
			want: helpers.JsonResponse{
				Code:    http.StatusBadRequest,
				Message: http.StatusText(http.StatusBadRequest),
				Data:    nil,
				IsError: true,
			},
		},
		{
			name: "ERROR | Error from Post.Delete (no posts found)",
			fields: fields{
				postService: *mocks.NewIPostService(t),
			},
			mock: func(s *mocks.IPostService) {
				s.On("Delete", mockID).Return(gorm.ErrRecordNotFound)
			},
			mockParamsFromMiddleware: true,
			want: helpers.JsonResponse{
				Code:    http.StatusNotFound,
				Message: errn.ErrNoPostsFound.Error(),
				Data:    nil,
				IsError: true,
			},
		},
		{
			name: "ERROR | Error from Post.Delete (other errors)",
			fields: fields{
				postService: *mocks.NewIPostService(t),
			},
			mock: func(s *mocks.IPostService) {
				s.On("Delete", mockID).Return(mockError)
			},
			mockParamsFromMiddleware: true,
			want: helpers.JsonResponse{
				Code:    http.StatusInternalServerError,
				Message: mockError.Error(),
				Data:    nil,
				IsError: true,
			},
		},
		{
			name: "SUCCESS",
			fields: fields{
				postService: *mocks.NewIPostService(t),
			},
			mock: func(s *mocks.IPostService) {
				s.On("Delete", mockID).Return(nil)
			},
			mockParamsFromMiddleware: true,
			want: helpers.JsonResponse{
				Code:    http.StatusOK,
				Message: http.StatusText(http.StatusOK),
				Data:    float64(mockID),
				IsError: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				services: &services.Services{
					Post: &tt.fields.postService,
				},
			}

			tt.mock(&tt.fields.postService)
			r := helpers.SetUpRouter()
			endpoint := "/posts/1"

			if tt.mockParamsFromMiddleware {
				r.DELETE(endpoint, helpers.MiddlewareMockID(mockIDParams), h.DeletePost)
			} else {
				r.DELETE(endpoint, h.DeletePost)
			}

			req, _ := http.NewRequest(
				http.MethodDelete,
				endpoint,
				nil,
			)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			var response helpers.JsonResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			assert.Equal(t, tt.want.Code, w.Code)
			assert.Equal(t, tt.want, response)
		})
	}
}

func TestHandler_GetAllTypes(t *testing.T) {
	mockError := fmt.Errorf("error")
	mockTypes := []*models.PostType{}

	type fields struct {
		postService mocks.IPostService
	}
	tests := []struct {
		name   string
		fields fields
		mock   func(*mocks.IPostService)
		want   helpers.JsonResponse
	}{
		{
			name: "ERROR | Error from Post.GetTypes",
			fields: fields{
				postService: *mocks.NewIPostService(t),
			},
			mock: func(s *mocks.IPostService) {
				s.On("GetTypes").Return(nil, mockError)
			},
			want: helpers.JsonResponse{
				Code:    http.StatusInternalServerError,
				Message: http.StatusText(http.StatusInternalServerError),
				Data:    nil,
				IsError: true,
			},
		},
		{
			name: "SUCCESS",
			fields: fields{
				postService: *mocks.NewIPostService(t),
			},
			mock: func(s *mocks.IPostService) {
				s.On("GetTypes").Return(mockTypes, nil)
			},
			want: helpers.JsonResponse{
				Code:    http.StatusOK,
				Message: http.StatusText(http.StatusOK),
				Data:    []interface{}{},
				IsError: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				services: &services.Services{
					Post: &tt.fields.postService,
				},
			}

			tt.mock(&tt.fields.postService)
			r := helpers.SetUpRouter()

			endpoint := "/types"
			r.GET(endpoint, h.GetAllTypes)

			req, _ := http.NewRequest(
				http.MethodGet,
				endpoint,
				nil,
			)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			var response helpers.JsonResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			assert.Equal(t, tt.want.Code, w.Code)
			assert.Equal(t, tt.want, response)
		})
	}
}

func TestHandler_GetAllCategories(t *testing.T) {
	mockError := fmt.Errorf("error")
	mockCategories := []*models.Category{}

	type fields struct {
		postService mocks.IPostService
	}
	tests := []struct {
		name   string
		fields fields
		mock   func(*mocks.IPostService)
		want   helpers.JsonResponse
	}{
		{
			name: "ERROR | Error from Post.GetCategories",
			fields: fields{
				postService: *mocks.NewIPostService(t),
			},
			mock: func(s *mocks.IPostService) {
				s.On("GetCategories").Return(nil, mockError)
			},
			want: helpers.JsonResponse{
				Code:    http.StatusInternalServerError,
				Message: http.StatusText(http.StatusInternalServerError),
				Data:    nil,
				IsError: true,
			},
		},
		{
			name: "SUCCESS",
			fields: fields{
				postService: *mocks.NewIPostService(t),
			},
			mock: func(s *mocks.IPostService) {
				s.On("GetCategories").Return(mockCategories, nil)
			},
			want: helpers.JsonResponse{
				Code:    http.StatusOK,
				Message: http.StatusText(http.StatusOK),
				Data:    []interface{}{},
				IsError: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				services: &services.Services{
					Post: &tt.fields.postService,
				},
			}

			tt.mock(&tt.fields.postService)
			r := helpers.SetUpRouter()

			endpoint := "/categories"
			r.GET(endpoint, h.GetAllCategories)

			req, _ := http.NewRequest(
				http.MethodGet,
				endpoint,
				nil,
			)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			var response helpers.JsonResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			assert.Equal(t, tt.want.Code, w.Code)
			assert.Equal(t, tt.want, response)
		})
	}
}
