package api

import (
	"errors"
	"net/http"

	"github.com/nurtai325/alaman/internal/auth"
	"github.com/nurtai325/alaman/internal/service"
)

func (app *app) handleLogged(w http.ResponseWriter, r *http.Request) {
	if auth.IsLogged(r) {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}
}

func (app *app) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		app.execute(w, tLogin, "", layoutData{
			Data: "",
		})
		return
	}
	err := r.ParseForm()
	if err != nil {
		app.error(w, err)
	}
	phone := r.FormValue("phone")
	password := r.FormValue("password")
	sessionCookie, err := app.service.Login(r.Context(), phone, password)
	if err != nil {
		if r.Header.Get("Device") == "app" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return
		}
		if errors.Is(err, service.ErrInternal) {
			app.error(w, err)
			return
		}
		w.Write([]byte(err.Error()))
		return
	}
	http.SetCookie(w, sessionCookie)
	if r.Header.Get("Device") == "app" {
		w.WriteHeader(http.StatusOK)
		return
	}
	redirect(w, "/")
	return
}

func (app *app) handleLogout(w http.ResponseWriter, r *http.Request) {
	emptyCookie := auth.DeleteSession(r)
	http.SetCookie(w, emptyCookie)
	redirect(w, "/login")
}
