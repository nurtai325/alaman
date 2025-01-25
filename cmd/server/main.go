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

func main() {
	infoLog, err := openLog("info", log.Lshortfile)
	errLog, err := openLog("error", 0)
	accessLog, err := openLog("access", 0)

	templates := template.Must(template.ParseGlob("./views/*.html"))
	templates = template.Must(templates.ParseGlob("./views/pages/users/*.html"))
	templates = template.Must(templates.ParseGlob("./views/pages/dashboard/*.html"))
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
