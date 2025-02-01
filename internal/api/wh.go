package api

import (
	"net/http"

	"github.com/nurtai325/alaman/internal/wh"
)

func (app *app) HandleQr(w http.ResponseWriter, r *http.Request) {
	data, err := wh.GetQr()
	if err != nil {
		app.error(w, err)
		return
	}
	app.execute(w, tQrTempl, "", data)
}
