package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/jorgeAM/go-template/internal/shared/errors"
	"github.com/jorgeAM/go-template/internal/user/app/command"
	"github.com/jorgeAM/go-template/internal/user/app/models"
	"github.com/jorgeAM/go-template/internal/user/app/query"
	"github.com/jorgeAM/go-template/internal/user/domain"
)

var _ StrictServerInterface = (*Server)(nil)

type Server struct {
	createUser *command.CreateUser
	getUser    *query.GetUser
}

func NewServer(userRepository domain.UserRepository) *Server {
	return &Server{
		createUser: command.NewCreateUser(userRepository),
		getUser:    query.NewGetUser(userRepository),
	}
}

func Register(_ context.Context, r chi.Router, userRepository domain.UserRepository) error {
	handler := NewStrictHandlerWithOptions(NewServer(userRepository), nil, StrictHTTPServerOptions{
		RequestErrorHandlerFunc:  jsonErrorHandler(http.StatusBadRequest, errors.BadRequestCode),
		ResponseErrorHandlerFunc: jsonErrorHandler(http.StatusInternalServerError, errors.InternalCode),
	})

	HandlerFromMux(handler, r)

	return nil
}

func (s *Server) CreateUser(ctx context.Context, request CreateUserRequestObject) (CreateUserResponseObject, error) {
	res, err := s.createUser.Handle(ctx, &command.CreateUserCommand{
		Name:     request.Body.Name,
		Email:    request.Body.Email,
		Password: request.Body.Password,
	})
	if err != nil {
		if errors.Is(err, domain.ErrInvalidUser) {
			return CreateUser400JSONResponse{badRequest(err)}, nil
		}

		return CreateUser500JSONResponse{internalError(err)}, nil
	}

	return CreateUser201JSONResponse(toUserInfo(res)), nil
}

func (s *Server) GetUser(ctx context.Context, request GetUserRequestObject) (GetUserResponseObject, error) {
	res, err := s.getUser.Handle(ctx, &query.GetUserQuery{UserID: request.Id})
	if err != nil {
		if errors.Is(err, domain.ErrInvalidUser) {
			return GetUser400JSONResponse{badRequest(err)}, nil
		}

		if errors.Is(err, domain.ErrUserNotFound) {
			return GetUser404JSONResponse{notFound(err)}, nil
		}

		return GetUser500JSONResponse{internalError(err)}, nil
	}

	return GetUser200JSONResponse(toUserInfo(res)), nil
}

func toUserInfo(info *models.UserInfo) UserInfo {
	return UserInfo{
		Id:        info.ID,
		Name:      info.Name,
		Email:     info.Email,
		CreatedAt: info.CreatedAt,
		UpdatedAt: info.UpdatedAt,
	}
}

func badRequest(err error) BadRequestJSONResponse {
	return BadRequestJSONResponse{Code: errors.BadRequestCode.String(), Message: err.(*errors.Error).Message()}
}

func notFound(err error) NotFoundJSONResponse {
	return NotFoundJSONResponse{Code: errors.NotFoundCode.String(), Message: err.(*errors.Error).Message()}
}

func internalError(err error) InternalErrorJSONResponse {
	return InternalErrorJSONResponse{Code: errors.InternalCode.String(), Message: err.Error()}
}

func jsonErrorHandler(status int, code errors.Code) func(http.ResponseWriter, *http.Request, error) {
	return func(w http.ResponseWriter, _ *http.Request, err error) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(Error{Code: code.String(), Message: err.Error()})
	}
}
