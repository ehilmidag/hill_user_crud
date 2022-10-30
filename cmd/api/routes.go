package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) router() *httprouter.Router {
	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, "/healthcheck", app.healthCheckHandler)
	router.HandlerFunc(http.MethodPut, "/users", app.addUserHandler)
	router.HandlerFunc(http.MethodPatch, "/users", app.editUserHandler)
	router.HandlerFunc(http.MethodDelete, "/users/:id", app.deleteUserByIdHandler)
	router.HandlerFunc(http.MethodGet, "/users/:id", app.getUserByIdHandler)
	router.HandlerFunc(http.MethodGet, "/users", app.getAllUsersHandler)

	return router

}
