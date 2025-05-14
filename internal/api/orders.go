package api

import (
	"net/http"
	"strconv"
)

func (app *app) handleOrdersGet(w http.ResponseWriter, r *http.Request) {
	orderIdStr := r.PathValue("id")
	orderId, err := strconv.Atoi(orderIdStr)
	if err != nil {
		app.error(w, err)
		return
	}
	order, err := app.service.GetOrder(r.Context(), orderId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	app.execute(w, tOrders, "/pages/orders", layoutData{
		BarsData: barsData{
			Page:     "orders",
			PageName: "Тапсырыс күйі",
			Pages:    getPage(r),
		},
		User: app.service.GetAuthUser(r),
		Data: order,
	})
}
