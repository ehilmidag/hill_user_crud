package main

import (
	"fmt"
	"net/http"

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
	user := &data.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
	}
	v := validator.New()

	if data.ValidateUser(v, user); !v.IsValid() {
		app.badRequestResponse(w, r)
		return
	}

	fmt.Fprintf(w, "%+v", input)
}

func (app *application) getUserByIdHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDFromParameter(r)

	if err != nil {
		app.badRequestResponse(w, r)
		return
	}
	user := data.User{
		ID:       id,
		Name:     "Test",
		Email:    "ehdag@gmail.com",
		Password: "1234",
	}
	err = app.writeJsonHelper(w, http.StatusOK, envelope{"user": user}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) editUserHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "sa gardaş")
}

func (app *application) deleteUserByIdHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "sa gardaş")
}

func (app *application) getAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "sa gardaş")
}
