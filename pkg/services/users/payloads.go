package users

import (
	"github.com/kkeuning/go-api-example/pkg/models"
)

// ShowPayload ...
type ShowPayload struct {
	ID int
}

// ListPayload ...
type ListPayload struct {
	ID int
}

// DeletePayload ...
type DeletePayload struct {
	ID int
}

// UpdatePayload ...
type UpdatePayload struct {
	User models.User
}
