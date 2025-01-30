package main

import (
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
)

func openLog(name string, lFlags int) (*log.Logger, error) {
	f, err := os.OpenFile(fmt.Sprintf("./logs/%s.log", name), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0640)
	if err != nil {
		return nil, err
	}
	logger := log.New(f, "", log.LstdFlags|log.Lmicroseconds|lFlags)
	return logger, nil
}

func parsePages(templ *template.Template, pages ...string) *template.Template {
	for _, page := range pages {
		glob := fmt.Sprintf("./views/pages/%s/*.html", page)
		templ = template.Must(templ.ParseGlob(glob))
	}
	return templ
}

func main() {
	infoLog, err := openLog("info", log.Lshortfile)
	errLog, err := openLog("error", 0)
	accessLog, err := openLog("access", 0)

	templates := template.Must(template.ParseGlob("./views/*.html"))
	templates = parsePages(templates, "users", "dashboard", "products", "leads")
	conf, err := config.New()
	if err != nil {
		panic(err)
	}
	newDB, err := db.New(conf)
	if err != nil {
		panic(err)
	}
	newService := service.New(repository.New(newDB))
	app := api.New(http.NewServeMux(), templates, newService, infoLog, accessLog, errLog)
	go auth.Cleanup()

	err = app.Run(conf)
	if err != nil {
		panic(err)
	}
}
