package api

import (
	"net/http"

	"github.com/nurtai325/alaman/internal/auth"
	"github.com/nurtai325/alaman/internal/service"
)

func (app *app) handleRoot(w http.ResponseWriter, r *http.Request) {
	if auth.GetUser(r).Role == auth.AdminRole {
		app.handleDashBoard(w, r)
	} else if auth.GetUser(r).Id == 0 {
		app.execute(w, "index.html", "", nil)
	} else {
		app.handleLeadsGet(w, r)
	}
}

type dashboardData struct {
	ChartsData service.ChartsData
}

func (app *app) handleDashBoard(w http.ResponseWriter, r *http.Request) {
	data, err := app.service.GetSalesData(r.Context())
	if err != nil {
		app.error(w, err)
		return
	}
	app.execute(w, tDashboard, "/pages/dashboard", layoutData{
		BarsData: barsData{
			Page:     "dashboard",
			PageName: "Негізгі",
			Pages:    getPage(r),
		},
		User: app.service.GetAuthUser(r),
		Data: data,
	})
}

func (app *app) handleEmpty(w http.ResponseWriter, r *http.Request) {
	app.execute(w, tEmpty, "", nil)
}

func getPage(r *http.Request) []string {
	user := auth.GetUser(r)
	if user.Role == auth.AdminRole {
		return adminPages
	} else if user.Role == auth.RopRole {
		return ropPages
	} else {
		return normPages
	}
}
