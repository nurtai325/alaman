package api

import (
	"fmt"
	"net/http"
	"strconv"

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

func (app *app) handleProductsPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.error(w, err)
		return
	}
	name := r.FormValue("name")
	inStockStr := r.FormValue("in_stock")
	inStock, err := strconv.Atoi(inStockStr)
	if err != nil {
		app.error(w, err)
		return
	}
	priceStr := r.FormValue("price")
	price, err := strconv.Atoi(priceStr)
	if err != nil {
		app.error(w, err)
		return
	}
	stockPriceStr := r.FormValue("stock_price")
	stockPrice, err := strconv.Atoi(stockPriceStr)
	if err != nil {
		app.error(w, err)
		return
	}
	product, err := app.service.InsertProduct(r.Context(), name, inStock, price, stockPrice)
	if err != nil {
		app.error(w, err)
		return
	}
	app.execute(w, tProductRow, "", product)
	return
}

func (app *app) handleProductsPut(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.error(w, err)
		return
	}
	name := r.FormValue("name")
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		app.error(w, err)
		return
	}
	priceStr := r.FormValue("price")
	price, err := strconv.Atoi(priceStr)
	if err != nil {
		app.error(w, err)
		return
	}
	stockPriceStr := r.FormValue("stock_price")
	stockPrice, err := strconv.Atoi(stockPriceStr)
	if err != nil {
		app.error(w, err)
		return
	}
	product, err := app.service.UpdateProduct(r.Context(), name, id, price, stockPrice)
	if err != nil {
		app.error(w, err)
		return
	}
	app.execute(w, tProductRow, "", product)
	return
}

func (app *app) handleProductsEdit(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		app.error(w, err)
		return
	}
	product, err := app.service.GetProduct(r.Context(), id)
	if err != nil {
		app.error(w, err)
		return
	}
	app.execute(w, tProductRowEdit, "", product)
	return
}

func (app *app) handleProductsDelete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		app.error(w, err)
		return
	}
	_, err = app.service.DeleteProduct(r.Context(), id)
	if err != nil {
		app.error(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}

func (app *app) handleProductsAdd(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		app.error(w, err)
		return
	}
	err = r.ParseForm()
	if err != nil {
		app.error(w, err)
		return
	}
	quantityStr := r.FormValue("quantity")
	quantity, err := strconv.Atoi(quantityStr)
	if err != nil {
		w.Header().Add("HX-Retarget", fmt.Sprintf("#product-row-errors-%d", id))
		w.Header().Add("HX-Reswap", "innerHTML")
		w.WriteHeader(http.StatusUnprocessableEntity)
		app.execute(w, tAlert, "", "сан жазыңыз")
		return
	}
	inStock, err := app.service.AddStockProduct(r.Context(), id, quantity)
	if err != nil {
		app.error(w, err)
		return
	}
	app.execute(w, tText, "", inStock)
	return
}
