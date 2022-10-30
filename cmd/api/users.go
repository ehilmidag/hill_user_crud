package main

import (
	"fmt"
	"net/http"

	"hilmi.dag/internal/data"
)

func (app *application) addUserHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "sa gardaş")
}

func (app *application) getUserByIdHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDFromParameter(r)

	if err != nil {
		app.badRequestResponse(w, r)
		return
	}
	user := data.User{
		ID:    id,
		Name:  "Test",
		Email: "ehdag@gmail.com",
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
