package api

import (
	"errors"
	"net/http"

	"github.com/nurtai325/alaman/internal/auth"
	"github.com/nurtai325/alaman/internal/service"
)

func (app *app) handleLoginGet(w http.ResponseWriter, _ *http.Request) {
	app.execute(w, tLayout, layoutData{
		Page: pLogin,
		Error: "",
	})
}

func (app *app) handleLoginPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.error(w, err)
		return
	}
	phone := r.FormValue("phone")
	password := r.FormValue("password")
	sessionCookie, err := app.service.Login(r.Context(), phone, password)
	if err != nil {
		if !errors.Is(err, service.ErrInvalidLoginInfo) {
			app.error(w, err)
			return
		}
		app.execute(w, tLayout, layoutData{
			Page: pLogin,
			Error: err.Error(),
		})
		return
	}
	http.SetCookie(w, sessionCookie)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (app *app) handleLogout(w http.ResponseWriter, r *http.Request) {
	emptyCookie := auth.DeleteSession(r)
	http.SetCookie(w, emptyCookie)
	http.Redirect(w, r, "/login", http.StatusFound)
}
