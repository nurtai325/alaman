package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/nurtai325/alaman/internal/service"
)

var (
	ErrNotNumber = errors.New("сан жазыңыз")
)

const (
	productModalErrs = "#product-modal-errors"
	productRowErrs   = "#product-row-errors"
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
	if r.Header.Get(acceptHeader) == jsonContentType {
		resp, err := json.Marshal(products)
		if err != nil {
			app.error(w, err)
			return
		}
		w.Header().Add(contentTypeHeader, jsonContentType)
		_, err = w.Write(resp)
		if err != nil {
			app.error(w, err)
			return
		}
		return
	}
	app.execute(w, tProducts, "/pages/products", layoutData{
		BarsData: barsData{
			Page:     "products",
			PageName: "Өнімдер",
			Pages:    getPage(r),
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
	saleCountStr := r.FormValue("sale_count")
	saleCount, err := strconv.Atoi(saleCountStr)
	if err != nil {
		app.error(w, err)
		return
	}
	inStockStr := r.FormValue("in_stock")
	inStock, err := strconv.Atoi(inStockStr)
	if err != nil {
		app.error(w, err)
		return
	}
	priceStr := r.FormValue("price")
	price, err := strconv.Atoi(priceStr)
	if err != nil {
		if errors.Is(err, service.ErrInternal) {
			app.error(w, err)
			return
		}
		app.errorHx(w, tText, productModalErrs, ErrNotNumber.Error())
		return
	}
	stockPriceStr := r.FormValue("stock_price")
	stockPrice, err := strconv.Atoi(stockPriceStr)
	if err != nil {
		if errors.Is(err, service.ErrInternal) {
			app.error(w, err)
			return
		}
		app.errorHx(w, tText, productModalErrs, ErrNotNumber.Error())
		return
	}
	product, err := app.service.InsertProduct(r.Context(), name, inStock, price, stockPrice, saleCount)
	if err != nil {
		if errors.Is(err, service.ErrInternal) {
			app.error(w, err)
			return
		}
		app.errorHx(w, tText, productModalErrs, err.Error())
		return
	}
	w.Header().Add("HX-Trigger", closeModalEvent)
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
		app.errorHx(w, tAlert, productRowErrs, ErrNotNumber.Error())
		return
	}
	saleCountStr := r.FormValue("sale_count")
	saleCount, err := strconv.Atoi(saleCountStr)
	if err != nil {
		app.errorHx(w, tAlert, productRowErrs, ErrNotNumber.Error())
		return
	}
	priceStr := r.FormValue("price")
	price, err := strconv.Atoi(priceStr)
	if err != nil {
		app.errorHx(w, tAlert, productRowErrs, ErrNotNumber.Error())
		return
	}
	stockPriceStr := r.FormValue("stock_price")
	stockPrice, err := strconv.Atoi(stockPriceStr)
	if err != nil {
		app.errorHx(w, tAlert, productRowErrs, ErrNotNumber.Error())
		return
	}
	product, err := app.service.UpdateProduct(r.Context(), name, id, price, stockPrice, saleCount)
	if err != nil {
		if errors.Is(err, service.ErrInternal) {
			app.error(w, err)
			return
		}
		app.errorHx(w, tAlert, productRowErrs, err.Error())
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
	id, quantity, err := app.validStockParams(w, r)
	if err != nil {
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

func (app *app) handleProductsRemove(w http.ResponseWriter, r *http.Request) {
	id, quantity, err := app.validStockParams(w, r)
	if err != nil {
		return
	}
	inStock, err := app.service.RemoveStockProduct(r.Context(), id, quantity)
	if err != nil {
		app.error(w, err)
		return
	}
	app.execute(w, tText, "", inStock)
	return
}

func (app *app) validStockParams(w http.ResponseWriter, r *http.Request) (int, int, error) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		app.error(w, err)
		return 0, 0, ErrNotNumber
	}
	err = r.ParseForm()
	if err != nil {
		app.error(w, err)
		return 0, 0, ErrNotNumber
	}
	quantityStr := r.FormValue("quantity")
	quantity, err := strconv.Atoi(quantityStr)
	if err != nil {
		app.errorHx(w, tAlert, productRowErrs, ErrNotNumber.Error())
		return 0, 0, ErrNotNumber
	}
	return id, quantity, nil
}
