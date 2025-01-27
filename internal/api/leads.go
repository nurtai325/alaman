package api

import "net/http"

type leadsContent struct {
}

func (app *app) handleLeadsGet(w http.ResponseWriter, r *http.Request) {
	app.execute(w, tLeads, "/pages/leads", layoutData{
		BarsData: barsData{
			Page:     "leads",
			PageName: "Лидтер",
			Pages:    adminPages,
		},
		User: app.service.GetAuthUser(r),
		Data: leadsContent{},
	})
}
