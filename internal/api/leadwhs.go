package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/nurtai325/alaman/internal/service"
	"github.com/nurtai325/alaman/internal/wh"
)

const (
	leadWhModalErrs = "#leadwh-modal-errors"
	leadWhRowErrs   = "#leadwh-row-errors"
)

type leadWhContent struct {
	Rows []service.LeadWh
}

func (app *app) handleLeadwhsGet(w http.ResponseWriter, r *http.Request) {
	leadWhs, err := app.service.GetLeadWhs(r.Context(), 0, pagesLimit)
	if err != nil {
		app.error(w, err)
		return
	}
	app.execute(w, tLeadWhs, "/pages/leadwhs", layoutData{
		BarsData: barsData{
			Page:     "leadwhs",
			PageName: "Лид номерлері",
			Pages:    getPage(r),
		},
		User: app.service.GetAuthUser(r),
		Data: leadWhContent{
			Rows: leadWhs,
		},
	})
}

func (app *app) handleLeadWhsPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.error(w, err)
		return
	}
	phone := r.FormValue("phone")
	imagePath, err := app.service.GetLeadWhQr(phone)
	if err != nil {
		if !errors.Is(err, service.ErrAlreadyPaired) {
			app.error(w, err)
			return
		}
	} else {
		w.Header().Add("HX-Retarget", "#qr-section")
		w.Header().Add("HX-Reswap", "innerHTML")
		app.execute(w, tQrImage, "", imagePath)
		return
	}
	name := r.FormValue("name")
	leadWh, err := app.service.InsertLeadWh(r.Context(), name, phone)
	if err != nil {
		if errors.Is(err, service.ErrInternal) {
			app.error(w, err)
			return
		}
		app.errorHx(w, tText, leadWhModalErrs, err.Error())
		return
	}
	_, err = app.service.ConnectLeadWh(r.Context(), leadWh.Id, wh.GetJid(phone[1:]))
	if err != nil {
		app.error(w, err)
		return
	}
	w.Header().Add("HX-Trigger", closeModalEvent)
	app.execute(w, tLeadWhRow, "", leadWh)
	return
}

func (app *app) handleLeadWhsPut(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.error(w, err)
		return
	}
	name := r.FormValue("name")
	phone := r.FormValue("phone")
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		app.errorHx(w, tAlert, leadWhRowErrs, ErrNotNumber.Error())
		return
	}
	leadWh, err := app.service.UpdateLeadWh(r.Context(), id, name, phone)
	if err != nil {
		if errors.Is(err, service.ErrInternal) {
			app.error(w, err)
			return
		}
		app.errorHx(w, tAlert, leadWhRowErrs, err.Error())
		return
	}
	app.execute(w, tLeadWhRow, "", leadWh)
	return
}

func (app *app) handleLeadWhsEdit(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		app.error(w, err)
		return
	}
	leadWh, err := app.service.GetLeadWh(r.Context(), id)
	if err != nil {
		app.error(w, err)
		return
	}
	app.execute(w, tLeadWhRowEdit, "", leadWh)
	return
}

func (app *app) handleLeadWhsDelete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		app.error(w, err)
		return
	}
	_, err = app.service.DeleteLeadWh(r.Context(), id)
	if err != nil {
		app.error(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}
