package api

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/nurtai325/alaman/internal/auth"
	"github.com/nurtai325/alaman/internal/config"
	"github.com/nurtai325/alaman/internal/service"
)

type app struct {
	mux       *http.ServeMux
	templates *template.Template
	service   *service.Service
	infoLog   *log.Logger
	errLog    *log.Logger
	accessLog *log.Logger
}

func New(mux *http.ServeMux, templates *template.Template, service *service.Service, infoLog, accessLog, errLog *log.Logger) *app {
	return &app{
		mux:       mux,
		templates: templates,
		service:   service,
		infoLog:   infoLog,
		errLog:    errLog,
		accessLog: accessLog,
	}
}

func (app *app) register(pattern string, handler http.HandlerFunc, secured bool) {
	app.mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				app.errLog.Printf("unhandled panic: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()
		app.accessLog.Printf("%s %v %s", r.Method, r.URL, r.Header.Get("User-Agent"))
		if secured && !auth.IsLogged(r) {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		handler(w, r)
	})
}

func (app *app) Run(conf config.Config) error {
	app.registerHandlers()
	app.infoLog.Printf("started http server on port: %s", conf.PORT)
	return http.ListenAndServeTLS(fmt.Sprintf(":%s", conf.PORT), "cert/certificate.crt", "cert/private.key", app.mux)
}
