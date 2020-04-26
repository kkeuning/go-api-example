package users

import (
	"context"

	"github.com/kkeuning/go-api-example/pkg/models"
	"github.com/rs/zerolog"
)

type usersSvc struct {
	// db
	logger *zerolog.Logger
}

// Service ...
type Service interface {
	Show(context.Context, *ShowPayload) (*models.User, error)
	Delete(context.Context, *DeletePayload) error
	Update(context.Context, *UpdatePayload) error
	List(context.Context, *ListPayload) (*models.UserCollection, error)
	Create(context.Context, *models.User) (*models.UserCollection, error)
}

// NewUsersSvc ...
func NewUsersSvc(logger *zerolog.Logger) (Service, error) {
	return &usersSvc{logger}, nil
}

func (us *usersSvc) Create(ctx context.Context, payload *models.User) (*models.UserCollection, error) {
	var uc *models.UserCollection
	return uc, nil
}
func (us *usersSvc) Delete(ctx context.Context, payload *DeletePayload) error {
	return nil
}
func (us *usersSvc) Update(ctx context.Context, payload *UpdatePayload) error {
	return nil
}

// ErrorNotFound ...
type ErrorNotFound struct{}

func (e ErrorNotFound) Error() string {
	return "Not Found"
}

// ListUsers ...
func (us *usersSvc) List(ctx context.Context, payload *ListPayload) (*models.UserCollection, error) {
	if payload != nil {
		for _, u := range models.Users.GetUsers() {
			if u.ID == payload.ID {
				uc := &models.UserCollection{
					Users: []models.User{u},
				}
				return uc, nil
			}
		}
		return nil, &ErrorNotFound{}
	}
	// No id specified, list all users.
	uc := &models.UserCollection{Users: models.Users.GetUsers()}
	return uc, nil
}

// ShowUser ...
func (us *usersSvc) Show(ctx context.Context, payload *ShowPayload) (*models.User, error) {
	result, err := models.Users.GetUserByID(payload.ID)
	if err != nil {
		return nil, &ErrorNotFound{}
	}
	return result, nil
}
