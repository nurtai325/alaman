package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/nurtai325/alaman/internal/auth"
	"github.com/nurtai325/alaman/internal/service"
)

type userContent struct {
	Rows []service.User
}

func (app *app) handleUsersGet(w http.ResponseWriter, r *http.Request) {
	users, err := app.service.GetUsers(r.Context(), 0, pagesLimit)
	if err != nil {
		app.error(w, err)
		return
	}
	app.execute(w, tUsers, "/pages/users", layoutData{
		BarsData: barsData{
			Page:     "users",
			PageName: "Қызметкерлер",
			Pages:    getPage(r),
		},
		User: app.service.GetAuthUser(r),
		Data: userContent{
			Rows: users,
		},
	})
}

func (app *app) handleUsersPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.error(w, err)
		return
	}
	phone := r.FormValue("phone")
	imagePath, jid, err := app.service.GetLeadWhQr(phone)
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
	password := r.FormValue("password")
	passwordCheck := r.FormValue("passwordCheck")
	role, err := auth.ToRole(r.FormValue("role"))
	if err != nil {
		app.error(w, err)
		return
	}
	user, err := app.service.InsertUser(r.Context(), name, phone, password, passwordCheck, jid, role)
	if err != nil {
		if errors.Is(err, service.ErrInternal) {
			app.error(w, err)
			return
		}
		app.errorHx(w, tText, "#user-modal-errors", err.Error())
		return
	}
	w.Header().Add("HX-Trigger", closeModalEvent)
	app.execute(w, tUserRow, "", user)
	return
}

func (app *app) handleUsersPut(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.error(w, err)
		return
	}
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		app.error(w, err)
		return
	}
	name := r.FormValue("name")
	phone := r.FormValue("phone")
	role, err := auth.ToRole(r.FormValue("role"))
	if err != nil {
		app.error(w, err)
		return
	}
	user, err := app.service.UpdateUser(r.Context(), id, name, phone, role)
	if err != nil {
		if errors.Is(err, service.ErrInternal) {
			app.error(w, err)
			return
		}
		app.errorHx(w, tAlert, "#user-row-errors", err.Error())
		return
	}
	app.execute(w, tUserRow, "", user)
}

func (app *app) handleUsersEdit(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		app.error(w, err)
		return
	}
	user, err := app.service.GetUser(r.Context(), id)
	if err != nil {
		app.error(w, err)
		return
	}
	app.execute(w, tUserRowEdit, "", user)
	return
}

func (app *app) handleUsersDelete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		app.error(w, err)
		return
	}
	_, err = app.service.DeleteUser(r.Context(), id)
	if err != nil {
		app.error(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}
