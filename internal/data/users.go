package data

import (
	"context"
	"database/sql"
	"errors"
	"hilmi.dag/internal/validator"
	"time"
)

var (
	ErrDuplicateEmail = errors.New("UserEntity with this email aldready exist")
	ErrEditConflict   = errors.New("conflict")
)

type UserEntity struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserDTO struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
}

type UserModel struct {
	DB *sql.DB
}

func (u *UserDTO) ValidateUser(v *validator.Validator) {
	v.Check(u.Name != "", "name", "Bad request")
	v.Check(u.Password != "", "password", "Bad request")
	v.Check(u.Email != "", "email", "Bad request")
	v.IsMailValid(u.Email, "notvalidmail", "Bad request")
}

func (u *UserModel) AddUser(user *UserEntity) (*UserEntity, error) {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	var query = `
		INSERT INTO users (name, email, password,created_at,updated_at) VALUES ($1, $2, $3,$4,$5)
		RETURNING id`
	args := []interface{}{user.Name, user.Email, user.Password, user.CreatedAt, user.UpdatedAt}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := u.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return nil, ErrDuplicateEmail
		default:
			return nil, err
		}
	}
	return user, nil
}

func (u *UserModel) GetUserByID(id int64) (*UserEntity, error) {
	query := `
SELECT id, name, email FROM users
WHERE id = $1`
	var user UserEntity
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := u.DB.QueryRowContext(ctx, query, id).Scan(&user.ID,
		&user.Name, &user.Email,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UserModel) GetAllUsers() error {
	return nil
}

func (u *UserModel) UpdateUser(user *UserEntity) (*UserEntity, error) {
	query := ` UPDATE users
SET name = $2, password = $3 WHERE id = $1
RETURNING id,name,email,password`
	args := []interface{}{user.ID, user.Name, user.Password}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := u.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return nil, ErrDuplicateEmail
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrEditConflict
		default:
			return nil, err
		}
	}
	return user, nil
}

func (u *UserModel) DeleteUser(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	query := `DELETE FROM users WHERE id = $1`

	result, err := u.DB.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil

}

func (u *UserDTO) ConvertDTOtoEntity() *UserEntity {
	return &UserEntity{
		Name:     u.Name,
		Email:    u.Email,
		Password: u.Password,
	}
}

func (u *UserEntity) ConvertEntitytoDTO() UserDTO {
	return UserDTO{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	}
}
