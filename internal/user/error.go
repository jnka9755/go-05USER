package user

import (
	"errors"
	"fmt"
)

var ErrFirstNameRequired = errors.New("first_name is required")
var ErrLastNameRequired = errors.New("last_name is required")

type ErrNotFound struct {
	UserID string
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("User with ID -> '%s' doesn't exist", e.UserID)
}
