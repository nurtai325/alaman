package api

import (
	"fmt"
	"net/http"

	"github.com/nurtai325/alaman/internal/service"
)

type templateName string

const (
	tLayout    templateName = "layout.html"
	tLogin     templateName = "login.html"
	tDashboard templateName = "dashboard.html"
	tCharts    templateName = "charts"
	tAlert     templateName = "alert"
	tText      templateName = "text"
	tEmpty     templateName = "empty"

	tUsers       templateName = "users.html"
	tUserRow     templateName = "user-row"
	tUserRowEdit templateName = "user-row-edit"

	tOrders templateName = "orders.html"

	tProducts       templateName = "products.html"
	tProductRow     templateName = "product-row"
	tProductRowEdit templateName = "product-row-edit"

	tLeads                templateName = "leads.html"
	tLeadsNewForm         templateName = "new-lead-form"
	tLeadsNewCell         templateName = "lead-cell-new"
	tLeadsAssignedCell    templateName = "lead-cell-assigned"
	tLeadsInDeliveryCell  templateName = "lead-cell-in-delivery"
	tLeadsCompletedCell   templateName = "lead-cell-completed"
	tLeadsNewCells        templateName = "lead-cells-new"
	tLeadsAssignedCells   templateName = "lead-cells-assigned"
	tLeadsInDeliveryCells templateName = "lead-cells-in-delivery"
	tLeadsCompletedCells  templateName = "lead-cells-completed"
	tLeadsProduct         templateName = "leads-product"

	tReports        templateName = "reports.html"
	tReportsRow     templateName = "reports-row"
	tReportsRowEdit templateName = "reports-row-edit"

	tLeadWhs       templateName = "leadwhs.html"
	tLeadWhRow     templateName = "leadwh-row"
	tLeadWhRowEdit templateName = "leadwh-row-edit"

	tChats    templateName = "chats.html"
	tMessages templateName = "messages"

	tQrTempl templateName = "qr"
	tQrImage templateName = "qr-image"

	pagesLimit = 1000

	openModalEvent  = "openModal"
	closeModalEvent = "closeModal"

	jsonContentType   = "application/json"
	contentTypeHeader = "Content-Type"
	acceptHeader      = "Accept"
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
	adminPages = []string{"dashboard", "leads", "chats", "products", "reports", "leadwhs", "users"}
	ropPages   = []string{"dashboard", "leads", "chats"}
	normPages  = []string{"leads"}
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

func (app *app) errorHx(w http.ResponseWriter, template templateName, elemId, msg string) {
	w.Header().Add("HX-Retarget", elemId)
	w.Header().Add("HX-Reswap", "innerHTML")
	w.WriteHeader(http.StatusUnprocessableEntity)
	app.execute(w, template, "", msg)
}

func redirect(w http.ResponseWriter, location string) {
	w.Header().Add("HX-Redirect", location)
	w.WriteHeader(http.StatusOK)
}
