package models

import (
	"fmt"
	"time"
)

// Pre-loaded users for demonstration purposes
var initialUsers = []User{
	{
		FirstName: "Rob",
		LastName:  "Pike",
	},
	{
		FirstName: "Ken",
		LastName:  "Thompson",
	},
	{
		FirstName: "Robert",
		LastName:  "Griesemer",
	},
	{
		FirstName:     "Russ",
		MiddleInitial: "S",
		LastName:      "Cox",
	},
}

// Users is a package level variable acting as an in-memory user database
var Users UserCollection

func init() {
	for _, y := range initialUsers {
		Users.AddUser(y)
	}
}

// User represents a user of the system
type User struct {
	ID            int        `json:"id"`
	FirstName     string     `json:"first_name"`
	MiddleInitial string     `json:"middle_initial,omitempty"`
	LastName      string     `json:"last_name"`
	CreatedAt     *time.Time `json:"-"`
	UpdatedAt     *time.Time `json:"-"`
}

// UserCollection is a collection of user records
type UserCollection struct {
	Users []User
}

// AddUser will add a user if it doesn't already exist or return an error
func (uc *UserCollection) AddUser(u User) (*User, error) {
	nextID := len(uc.Users) + 1 // ID begins with 1
	u.ID = nextID
	for _, y := range uc.Users {
		if y.FirstName == u.FirstName && y.LastName == y.LastName {
			// Not yet supporting multiple users of same name
			return nil, fmt.Errorf("user with that name already exists")
		}
	}
	u.CreatedAt = &[]time.Time{time.Now().UTC()}[0]
	u.UpdatedAt = &[]time.Time{time.Now().UTC()}[0]
	uc.Users = append(uc.Users, u)
	return &u, nil
}

// GetUserByID returns the user record matching privided ID
func (uc UserCollection) GetUserByID(id int) (*User, error) {
	for _, y := range uc.Users {
		if y.ID == id {
			return &y, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

// GetUserByName will return the first user matching firstName and LastName
// This may not work in the real world since names are not unique
func (uc UserCollection) GetUserByName(firstName string, lastName string) (*User, error) {
	for _, y := range uc.Users {
		if y.FirstName == firstName && y.LastName == lastName {
			return &y, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

// GetUsers returns the slice of all users
func (uc UserCollection) GetUsers() []User {
	return uc.Users
}

// UpdateUser will overwrite current user record with new data
func (uc *UserCollection) UpdateUser(u User) error {
	for i := range uc.Users {
		if uc.Users[i].ID == u.ID {
			// Currently no partial updates supported since all struct fields are required
			uc.Users[i] = u
			uc.Users[i].UpdatedAt = &[]time.Time{time.Now().UTC()}[0]
			return nil
		}
	}
	return fmt.Errorf("update failed likely due to missing or incorrect id")
}
