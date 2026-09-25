package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/jorgeAM/go-template/internal/shared/errors"
	"github.com/jorgeAM/go-template/internal/user/domain"
	usermock "github.com/jorgeAM/go-template/internal/user/mocks"
)

func TestServer(t *testing.T) {
	t.Parallel()

	user, err := domain.NewUser("Jorge", "jorge@example.com", "s3cure-pass")
	assert.NoError(t, err)

	tests := []struct {
		name       string
		method     string
		path       string
		body       string
		mock       func(repo *usermock.MockUserRepository)
		wantStatus int
		wantCode   errors.Code
	}{
		{
			name:   "creates a user",
			method: http.MethodPost,
			path:   "/api/v1/user",
			body:   `{"name":"Jorge","email":"jorge@example.com","password":"s3cure-pass"}`,
			mock: func(repo *usermock.MockUserRepository) {
				repo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)
			},
			wantStatus: http.StatusCreated,
		},
		{
			name:       "rejects invalid user data",
			method:     http.MethodPost,
			path:       "/api/v1/user",
			body:       `{"name":"Jorge","email":"not-an-email","password":"s3cure-pass"}`,
			mock:       func(repo *usermock.MockUserRepository) {},
			wantStatus: http.StatusBadRequest,
			wantCode:   errors.BadRequestCode,
		},
		{
			name:       "rejects a malformed body as JSON",
			method:     http.MethodPost,
			path:       "/api/v1/user",
			body:       `{`,
			mock:       func(repo *usermock.MockUserRepository) {},
			wantStatus: http.StatusBadRequest,
			wantCode:   errors.BadRequestCode,
		},
		{
			name:   "returns a user",
			method: http.MethodGet,
			path:   "/api/v1/user/" + user.ID().String(),
			mock: func(repo *usermock.MockUserRepository) {
				repo.EXPECT().FindByID(gomock.Any(), user.ID()).Return(user, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "rejects an invalid id",
			method:     http.MethodGet,
			path:       "/api/v1/user/not-a-uuid",
			mock:       func(repo *usermock.MockUserRepository) {},
			wantStatus: http.StatusBadRequest,
			wantCode:   errors.BadRequestCode,
		},
		{
			name:   "maps a missing user to 404",
			method: http.MethodGet,
			path:   "/api/v1/user/" + user.ID().String(),
			mock: func(repo *usermock.MockUserRepository) {
				repo.EXPECT().FindByID(gomock.Any(), user.ID()).
					Return(nil, errors.New(domain.ErrUserNotFound, "user not found"))
			},
			wantStatus: http.StatusNotFound,
			wantCode:   errors.NotFoundCode,
		},
		{
			name:   "maps a repository failure to 500",
			method: http.MethodGet,
			path:   "/api/v1/user/" + user.ID().String(),
			mock: func(repo *usermock.MockUserRepository) {
				repo.EXPECT().FindByID(gomock.Any(), user.ID()).Return(nil, assert.AnError)
			},
			wantStatus: http.StatusInternalServerError,
			wantCode:   errors.InternalCode,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			repo := usermock.NewMockUserRepository(ctrl)
			tt.mock(repo)

			router := chi.NewRouter()
			assert.NoError(t, Register(context.Background(), router, repo))

			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, httptest.NewRequest(tt.method, tt.path, strings.NewReader(tt.body)))

			assert.Equal(t, tt.wantStatus, rec.Code)
			assert.NotContains(t, rec.Body.String(), "password")
			assert.NotContains(t, rec.Body.String(), user.Password())

			if tt.wantCode != "" {
				var body Error
				assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
				assert.Equal(t, tt.wantCode.String(), body.Code)
				return
			}

			var body UserInfo
			assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
			assert.NotEmpty(t, body.Id)
			assert.Equal(t, "jorge@example.com", body.Email)
		})
	}
}
