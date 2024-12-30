package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/nurtai325/alaman/internal/auth"
)

type templateName string

type pageMeta struct {
	Name string
	Slug string
}

const (
	tLayout templateName = "layout"

	pageOffset = 15
)

var (
	ErrPageNotFound = errors.New("Бет табылмады")

	// pDashboard = pageMeta{
	// 	Slug: "dashboard",
	// 	Name: "Басты бет",
	// }
	pLeads = pageMeta{
		Slug: "leads",
		Name: "Лидтер",
	}
	pSales = pageMeta{
		Slug: "sales",
		Name: "Сатулар",
	}
	pRepors = pageMeta{
		Slug: "reports",
		Name: "Есеп",
	}
	pUsers = pageMeta{
		Slug: "users",
		Name: "Қызметкерлер",
	}
	pLogin = pageMeta{
		Slug: "login",
		Name: "Логин",
	}
	pages = []pageMeta{pLeads, pSales, pRepors, pUsers}
)

type layoutData struct {
	Page      pageMeta
	Pages     []pageMeta
	User      auth.User
	TableData tableData
	Error     string
}

type tableData struct {
	Resource string
	Columns  []string
	Rows     []row
	Page     int
	Error    string
}

type row struct {
	Id    int
	Cells []cell
}

type cellType int

const (
	inputCell cellType = iota
	selectCell
	dateCell
)

type cell struct {
	Type    cellType
	Content string
}

func (app *app) execute(w http.ResponseWriter, name templateName, data any) {
	err := app.templates.ExecuteTemplate(w, string(name), data)
	if err != nil {
		err = fmt.Errorf("executing template: %s data: %v: %w", name, data, err)
		app.error(w, err)
		return
	}
}

func (app *app) error(w http.ResponseWriter, err error) {
	app.errLog.Println(err)
	http.Error(w, "Сервер қатесі.", http.StatusInternalServerError)
	return
}
