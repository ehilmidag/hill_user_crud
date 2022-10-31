package data

import "hilmi.dag/internal/validator"

type User struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.Name != "", "name", "Bad request")
	v.Check(user.Password != "", "password", "Bad request")
	v.Check(user.Email != "", "email", "Bad request")
	v.IsMailValid(user.Email, "notvalidmail", "Bad request")
}
