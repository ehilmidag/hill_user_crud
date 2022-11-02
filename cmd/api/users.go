package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"hilmi.dag/internal/data"
	"hilmi.dag/internal/validator"
)

func (app *application) addUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJsonHelpler(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r)
		return
	}
	user := &data.UserDTO{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
	}
	v := validator.New()

	if user.ValidateUser(v); !v.IsValid() {
		app.badRequestResponse(w, r)
		return
	}

	userEntity := user.ConvertDTOtoEntity()

	modelfromDb, err := app.models.Users.AddUser(userEntity)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddErrorToMap("email", "UserEntity with that email already exists")
			app.alreadyExistResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	response := modelfromDb.ConvertEntitytoDTO()
	err = app.writeJsonHelper(w, http.StatusCreated, response, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getUserByIdHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDFromParameter(r)

	if err != nil {
		app.badRequestResponse(w, r)
		return
	}
	user, err := app.models.Users.GetUserByID(id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			app.notFoundResponse(w, r)

		default:
			app.serverErrorResponse(w, r, err)

		}
		return
	}
	userResponse := data.UserDTO{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
	err = app.writeJsonHelper(w, http.StatusOK, userResponse, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) editUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDFromParameter(r)

	if err != nil {
		app.badRequestResponse(w, r)
		return
	}
	user, err := app.models.Users.GetUserByID(id)
	fmt.Println(user)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			app.notFoundResponse(w, r)

		default:
			app.serverErrorResponse(w, r, err)

		}
		return
	}

	var input struct {
		Name     *string `json:"name"`
		Password *string `json:"password"`
	}
	err = app.readJsonHelpler(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r)
		return
	}
	if input.Name != nil {
		user.Name = *input.Name
	}

	if input.Password != nil {
		user.Password = *input.Password
	}

	if input.Password != nil || input.Name != nil {
		user.UpdatedAt = time.Now()
	}
	fmt.Println(&user)
	updatedUser, err := app.models.Users.UpdateUser(user)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	response := updatedUser.ConvertEntitytoDTO()

	err = app.writeJsonHelper(w, http.StatusOK, response, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteUserByIdHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDFromParameter(r)

	if err != nil {
		app.badRequestResponse(w, r)
		return
	}
	err = app.models.Users.DeleteUser(id)

	if err != nil {
		if err == data.ErrRecordNotFound {
			app.notFoundResponse(w, r)
			return
		} else {
			app.serverErrorResponse(w, r, err)
			return
		}

	}

	w.WriteHeader(http.StatusOK)

}

func (app *application) getAllUsersHandler(w http.ResponseWriter, r *http.Request) {

	users, err := app.models.Users.GetAllUsers()
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
	}
	usersResponse := data.ConvertEntityListtoDTO(users)
	err = app.writeJsonHelper(w, http.StatusOK, usersResponse, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
