package users

import (
	"context"

	"github.com/kkeuning/go-api-example/pkg/models"
	"github.com/rs/zerolog"
)

type usersSvc struct {
	db     *models.UserStorage
	logger *zerolog.Logger
}

// Service ...
type Service interface {
	Show(context.Context, *ShowPayload) (*models.User, error)
	Delete(context.Context, *DeletePayload) error
	Update(context.Context, *UpdatePayload) error
	List(context.Context, *ListPayload) ([]models.User, error)
	Create(context.Context, *models.User) ([]models.User, error)
}

// NewUsersSvc ...
func NewUsersSvc(logger *zerolog.Logger, db *models.UserStorage) (Service, error) {
	return &usersSvc{
		db:     db,
		logger: logger,
	}, nil
}

func (us *usersSvc) Create(ctx context.Context, payload *models.User) ([]models.User, error) {
	return us.db.Users, nil
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
func (us *usersSvc) List(ctx context.Context, payload *ListPayload) ([]models.User, error) {
	if payload != nil {
		for _, u := range us.db.GetUsers() {
			if u.ID == payload.ID {
				uc := []models.User{u}
				return uc, nil
			}
		}
		return nil, &ErrorNotFound{}
	}
	// No id specified, list all users.
	uc := us.db.GetUsers()
	return uc, nil
}

// ShowUser ...
func (us *usersSvc) Show(ctx context.Context, payload *ShowPayload) (*models.User, error) {
	result, err := us.db.GetUserByID(payload.ID)
	if err != nil {
		return nil, &ErrorNotFound{}
	}
	return result, nil
}
