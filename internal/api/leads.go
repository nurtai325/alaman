package api

import (
	"net/http"

	"github.com/nurtai325/alaman/internal/auth"
)

func (app *app) handleLeadsGet(w http.ResponseWriter, r *http.Request) {
	app.execute(w, tLayout, layoutData{
		Page:  pLeads,
		User:  auth.GetUser(r),
		Pages: pages,
	})
}
