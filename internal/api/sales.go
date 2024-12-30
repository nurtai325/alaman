package api

import (
	"net/http"

	"github.com/nurtai325/alaman/internal/auth"
)

func (app *app) handleSalesGet(w http.ResponseWriter, r *http.Request) {
	app.execute(w, tLayout, layoutData{
		Page:  pSales,
		User:  auth.GetUser(r),
		Pages: pages,
	})
}
