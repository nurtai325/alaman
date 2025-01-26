package api

import (
	"fmt"
	"net/http"

	"github.com/nurtai325/alaman/internal/service"
)

// send whatsapp message to the worker to attach him to the lead

type templateName string

const (
	tLayout    templateName = "layout.html"
	tLogin     templateName = "login.html"
	tDashboard templateName = "dashboard.html"
	tAlert     templateName = "alert"
	tText      templateName = "text"
	tEmpty     templateName = "empty"

	tUsers       templateName = "users.html"
	tUserRow     templateName = "user-row"
	tUserRowEdit templateName = "user-row-edit"

	tProducts       templateName = "products.html"
	tProductRow     templateName = "product-row"
	tProductRowEdit templateName = "product-row-edit"

	pagesLimit = 1000

	openModalEvent = "openModal"
)

type layoutData struct {
	BarsData barsData
	User     service.User
	Data     any
}

type barsData struct {
	Page     string
	PageName string
	Pages    []string
}

var (
	adminPages = []string{"dashboard", "leads", "products", "users"}
)

func (app *app) execute(w http.ResponseWriter, name templateName, dir string, data any) {
	t, err := app.templates.Clone()
	if err != nil {
		app.error(w, err)
		return
	}
	if dir != "" {
		path := fmt.Sprintf("./views%s/%s", dir, name)
		t, err = t.ParseFiles(path)
		if err != nil {
			app.error(w, err)
			return
		}
	}
	err = t.ExecuteTemplate(w, string(name), data)
	if err != nil {
		app.error(w, err)
	}
}

func (app *app) error(w http.ResponseWriter, err error) {
	app.errLog.Println(err)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("Сервер қатесі"))
}

func redirect(w http.ResponseWriter, location string) {
	w.Header().Add("HX-Redirect", location)
	w.WriteHeader(http.StatusOK)
}
