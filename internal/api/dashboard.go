package api

import (
	"net/http"
)

func (app *app) handleDashBoard(w http.ResponseWriter, r *http.Request) {
	app.execute(w, tDashboard, "/pages/dashboard", layoutData{
		BarsData: barsData{
			Page:     "dashboard",
			PageName: "Негізгі",
			Pages:    adminPages,
		},
		User: app.service.GetAuthUser(r),
		Data: "",
	})
}

func (app *app) handleEmpty(w http.ResponseWriter, r *http.Request) {
	app.execute(w, tEmpty, "", nil)
}
