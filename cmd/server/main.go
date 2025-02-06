package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/nurtai325/alaman/internal/api"
	"github.com/nurtai325/alaman/internal/auth"
	"github.com/nurtai325/alaman/internal/config"
	"github.com/nurtai325/alaman/internal/db"
	"github.com/nurtai325/alaman/internal/db/repository"
	"github.com/nurtai325/alaman/internal/service"
	_ "github.com/nurtai325/alaman/internal/timezone"
	"github.com/nurtai325/alaman/internal/wh"
)

func openLog(name string, lFlags int) (*log.Logger, error) {
	f, err := os.OpenFile(fmt.Sprintf("./logs/%s.log", name), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0640)
	if err != nil {
		return nil, err
	}
	logger := log.New(f, "", log.LstdFlags|log.Lmicroseconds|lFlags)
	return logger, nil
}

func parseTemplPages(templ *template.Template, pages ...string) *template.Template {
	for _, page := range pages {
		glob := fmt.Sprintf("./views/pages/%s/*.html", page)
		templ = template.Must(templ.ParseGlob(glob))
	}
	return templ
}

// TODO: report, lead sort
func main() {
	infoLog, err := openLog("info", log.Lshortfile)
	accessLog, err := openLog("access", 0)
	errLog, err := openLog("error", 0)
	log.SetOutput(errLog.Writer())

	templates := template.Must(template.ParseGlob("./views/*.html"))
	templates = parseTemplPages(templates, "users", "dashboard", "products", "leads", "reports", "leadwhs")
	conf, err := config.New()
	if err != nil {
		panic(err)
	}
	newDB, err := db.New(conf)
	if err != nil {
		panic(err)
	}
	newSqlDb, err := db.NewSql(conf)
	err = wh.InitContainer(newSqlDb)
	if err != nil {
		panic(err)
	}
	newService := service.New(repository.New(newDB))
	leadWhs, err := newService.GetLeadWhs(context.Background(), 0, api.PagesLimit)
	if err != nil {
		panic(err)
	}
	for _, leadWh := range leadWhs {
		err := wh.Connect(leadWh.Jid, wh.HandleLeadEvents)
		if err != nil {
			panic(err)
		}
	}
	go service.ListenNewLeads(newService)
	app := api.New(http.NewServeMux(), templates, newService, infoLog, accessLog, errLog)
	go auth.Cleanup()

	err = app.Run(conf)
	if err != nil {
		panic(err)
	}
}
