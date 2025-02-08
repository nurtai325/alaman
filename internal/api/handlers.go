package api

import (
	"net/http"
)

func (app *app) registerHandlers() {
	fs := http.FileServer(http.Dir("./assets"))
	app.mux.Handle("/assets/", http.StripPrefix("/assets/", fs))

	app.register("/", app.handleRoot, true)
	app.register("/dashboard", app.handleDashBoard, true)
	app.register("GET /login", app.handleLogin, false)
	app.register("POST /login", app.handleLogin, false)
	app.register("POST /logout", app.handleLogout, true)
	app.register("GET /empty", app.handleEmpty, true)

	app.register("GET /users", app.handleUsersGet, true)
	app.register("POST /users", app.handleUsersPost, true)
	app.register("PUT /users/{id}", app.handleUsersPut, true)
	app.register("GET /users/{id}/edit", app.handleUsersEdit, true)
	app.register("DELETE /users/{id}", app.handleUsersDelete, true)

	app.register("GET /products", app.handleProductsGet, true)
	app.register("POST /products", app.handleProductsPost, true)
	app.register("PUT /products/{id}", app.handleProductsPut, true)
	app.register("GET /products/{id}/edit", app.handleProductsEdit, true)
	app.register("DELETE /products/{id}", app.handleProductsDelete, true)
	app.register("PUT /products/{id}/add", app.handleProductsAdd, true)
	app.register("PUT /products/{id}/remove", app.handleProductsRemove, true)

	app.register("GET /leads", app.handleLeadsGet, true)
	app.register("GET /leads/new", app.handleLeadsNew, true)
	app.register("POST /leads", app.handleLeadsPost, true)
	app.register("PUT /leads/{id}/assign", app.handleLeadsAssign, true)
	app.register("POST /leads/sell", app.handleLeadsSell, true)
	app.register("POST /leads/{id}/complete", app.handleLeadsComplete, true)
	app.register("GET /leads/product", app.handleLeadsProduct, true)

	app.register("GET /reports", app.handleReportsGet, true)
	app.register("POST /reports", app.handleReportsPost, true)
	app.register("PUT /reports/{id}", app.handleReportsPut, true)
	app.register("GET /reports/{id}/edit", app.handleReportsEdit, true)
	app.register("DELETE /reports/{id}", app.handleReportsDelete, true)

	app.register("GET /leadwhs", app.handleLeadwhsGet, true)
	app.register("POST /leadwhs", app.handleLeadWhsPost, true)
	app.register("PUT /leadwhs/{id}", app.handleLeadWhsPut, true)
	app.register("GET /leadwhs/{id}/edit", app.handleLeadWhsEdit, true)
	app.register("DELETE /leadwhs/{id}", app.handleLeadWhsDelete, true)

	app.register("GET /chats", app.handleChatsGet, true)
	app.register("GET /messages/{id}", app.handleMessagesGet, true)
}
