package api

import (
	"net/http"
)

func (app *app) registerHandlers() {
	fs := http.FileServer(http.Dir("./assets"))
	app.mux.Handle("/assets/", http.StripPrefix("/assets/", fs))

	app.register("/", app.handleLeadsGet, true)
	app.register("GET /login", app.handleLoginGet, false)
	app.register("POST /login", app.handleLoginPost, false)
	app.register("POST /logout", app.handleLogout, true)

	app.register("GET /leads", app.handleLeadsGet, true)
	app.register("GET /sales", app.handleSalesGet, true)
	app.register("GET /reports", app.handleReportsGet, true)
	app.register("GET /users", app.handleUsersGet, true)
}
