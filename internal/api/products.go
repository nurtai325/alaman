package api

import (
	"net/http"

	"github.com/nurtai325/alaman/internal/service"
)

type productContent struct {
	Rows []service.Product
}

func (app *app) handleProductsGet(w http.ResponseWriter, r *http.Request) {
	products, err := app.service.GetProducts(r.Context(), 0, pagesLimit)
	if err != nil {
		app.error(w, err)
		return
	}
	app.execute(w, tProducts, "/pages/products", layoutData{
		BarsData: barsData{
			Page:     "products",
			PageName: "Өнімдер",
			Pages:    adminPages,
		},
		User: app.service.GetAuthUser(r),
		Data: productContent{
			Rows: products,
		},
	})
}
