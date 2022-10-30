package main

import (
	"net/http"
)

// her işlemde tek tek hataları dönmemek adına burda topladım
func (app *application) errorLog(r *http.Request, err error) {
	app.logger.Println(err)
}

func (app *application) responseError(w http.ResponseWriter, r *http.Request, statusCode int, message interface{}) {
	env := envelope{"error": message}

	err := app.writeJsonHelper(w, statusCode, env, nil)

	if err != nil {
		app.errorLog(r, err)
		w.WriteHeader(500)
	}
}

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorLog(r, err)

	message := "server error"
	app.responseError(w, r, http.StatusInternalServerError, message)
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "User with that id does not exist"
	app.responseError(w, r, http.StatusNotFound, message)
}

func (app *application) alreadyExistResponse(w http.ResponseWriter, r *http.Request) {
	message := "User with that email already exists"
	app.responseError(w, r, http.StatusForbidden, message)
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request) {
	message := "Bad request"
	app.responseError(w, r, http.StatusBadRequest, message)
}
