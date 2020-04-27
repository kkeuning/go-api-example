package models

import (
	"fmt"
	"time"

	"github.com/rs/zerolog"
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

// UserDB is a package level variable acting as an in-memory user database
var UserDB UserStorage

func init() {
	for _, y := range initialUsers {
		UserDB.AddUser(y)
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

// UserStorage ...
type UserStorage struct {
	Users []User
	Log   *zerolog.Logger
}

// AddUser will add a user if it doesn't already exist or return an error
func (us *UserStorage) AddUser(u User) (*User, error) {
	nextID := len(us.Users) + 1 // ID begins with 1
	u.ID = nextID
	for _, y := range us.Users {
		if y.FirstName == u.FirstName && y.LastName == y.LastName {
			// Not yet supporting multiple users of same name
			return nil, fmt.Errorf("user with that name already exists")
		}
	}
	u.CreatedAt = &[]time.Time{time.Now().UTC()}[0]
	u.UpdatedAt = &[]time.Time{time.Now().UTC()}[0]
	us.Users = append(us.Users, u)
	return &u, nil
}

// GetUserByID returns the user record matching privided ID
func (us UserStorage) GetUserByID(id int) (*User, error) {
	for _, y := range us.Users {
		if y.ID == id {
			return &y, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

// GetUserByName will return the first user matching firstName and LastName
// This may not work in the real world since names are not unique
func (us UserStorage) GetUserByName(firstName string, lastName string) (*User, error) {
	for _, y := range us.Users {
		if y.FirstName == firstName && y.LastName == lastName {
			return &y, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

// GetUsers returns the slice of all users
func (us UserStorage) GetUsers() []User {
	us.Log.Debug().Msg("Getting all users from collection.")
	return us.Users
}

// UpdateUser will overwrite current user record with new data
func (us *UserStorage) UpdateUser(u User) error {
	for i := range us.Users {
		if us.Users[i].ID == u.ID {
			// Currently no partial updates supported since all struct fields are required
			us.Users[i] = u
			us.Users[i].UpdatedAt = &[]time.Time{time.Now().UTC()}[0]
			return nil
		}
	}
	return fmt.Errorf("update failed likely due to missing or incorrect id")
}
