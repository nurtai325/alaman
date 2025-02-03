package api

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/nurtai325/alaman/internal/service"
)

const (
	reportsModalErrs = "#reports-modal-errors"
	reportsRowErrs   = "#reports-row-errors"
)

type reportsContent struct {
	Rows []service.Report
}

func (app *app) handleReportsGet(w http.ResponseWriter, r *http.Request) {
	reports, err := app.service.GetReports(r.Context())
	if err != nil {
		app.error(w, err)
		return
	}
	app.execute(w, tReports, "/pages/reports", layoutData{
		BarsData: barsData{
			Page:     "reports",
			PageName: "Есептер",
			Pages:    getPage(r),
		},
		User: app.service.GetAuthUser(r),
		Data: reportsContent{
			Rows: reports,
		},
	})
}

func (app *app) handleReportsPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.error(w, err)
		return
	}
	name := r.FormValue("name")
	startAtStr := r.FormValue("start-at")
	startAt, err := time.Parse(time.DateOnly, startAtStr)
	if err != nil {
		app.error(w, err)
		return
	}
	endAtStr := r.FormValue("end-at")
	endAt, err := time.Parse(time.DateOnly, endAtStr)
	if err != nil {
		app.error(w, err)
		return
	}
	report, err := app.service.InsertReport(r.Context(), name, startAt, endAt)
	if err != nil {
		app.error(w, err)
		return
	}
	w.Header().Add("HX-Trigger", closeModalEvent)
	app.execute(w, tReportsRow, "", report)
	return
}

func (app *app) handleReportsEdit(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		app.error(w, err)
		return
	}
	report, err := app.service.GetReport(r.Context(), id)
	if err != nil {
		app.error(w, err)
		return
	}
	app.execute(w, tReportsRowEdit, "", report)
	return
}

func (app *app) handleReportsPut(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.error(w, err)
		return
	}
	name := r.FormValue("name")
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		app.errorHx(w, tAlert, reportsRowErrs, ErrNotNumber.Error())
		return
	}
	report, err := app.service.UpdateReport(r.Context(), id, name)
	if err != nil {
		if errors.Is(err, service.ErrInternal) {
			app.error(w, err)
			return
		}
		app.errorHx(w, tAlert, reportsRowErrs, err.Error())
		return
	}
	app.execute(w, tReportsRow, "", report)
	return
}

func (app *app) handleReportsDelete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		app.error(w, err)
		return
	}
	_, err = app.service.DeleteReport(r.Context(), id)
	if err != nil {
		app.error(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}
