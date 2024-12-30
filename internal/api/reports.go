package api

import (
	"net/http"

	"github.com/nurtai325/alaman/internal/auth"
)

func (app *app) handleReportsGet(w http.ResponseWriter, r *http.Request) {
	app.execute(w, tLayout, layoutData{
		Page:  pRepors,
		User:  auth.GetUser(r),
		Pages: pages,
	})
}
