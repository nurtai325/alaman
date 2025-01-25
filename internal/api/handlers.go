package api

import (
	"net/http"
)

func (app *app) registerHandlers() {
	fs := http.FileServer(http.Dir("./assets"))
	app.mux.Handle("/assets/", http.StripPrefix("/assets/", fs))

	app.register("/", app.handleDashBoard, true)
	app.register("/dashboard", app.handleDashBoard, true)
	app.register("GET /login", app.handleLogin, false)
	app.register("POST /login", app.handleLogin, false)
	app.register("POST /logout", app.handleLogout, true)

	app.register("GET /users", app.handleUsersGet, true)
	app.register("POST /users", app.handleUsersPost, true)
	app.register("PUT /users/{id}", app.handleUsersPut, true)
	app.register("GET /users/{id}/edit", app.handleUsersEdit, true)
	app.register("DELETE /users/{id}", app.handleUsersDelete, true)
}
