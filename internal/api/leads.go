package api

import "net/http"

func (app *app) handleLeadsGet(w http.ResponseWriter, r *http.Request) {
	app.execute(w, tEmpty, "", nil)
}
