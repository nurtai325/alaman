package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/nurtai325/alaman/internal/auth"
	"github.com/nurtai325/alaman/internal/service"
)

var (
	ErrChooseUser       = errors.New("қызметкерді таңдаңыз")
	ErrInvalidCartItems = errors.New("invalid cart items format")
)

const (
	leadCellErrors = "#lead-cell-errors"
)

type leadsContent struct {
	New           []service.Lead
	Assigned      []service.Lead
	InDelivery    []service.Lead
	Completed     []service.Lead
	Users         []service.User
	Products      []service.Product
	Role          string
	NewLeadsCount int
}

func (app *app) handleLeadsGet(w http.ResponseWriter, r *http.Request) {
	var newLeads []service.Lead
	var assignedLeads []service.Lead
	var inDeliveryLeads []service.Lead
	var completedLeads []service.Lead
	user := auth.GetUser(r)
	if user.Role == auth.AdminRole || user.Role == auth.RopRole {
		leads, err := app.service.GetNewLeads(r.Context())
		if err != nil {
			app.error(w, err)
			return
		}
		newLeads = leads
	}
	if user.Role == auth.AdminRole || user.Role == auth.RopRole {
		leads, err := app.service.GetAssignedLeads(r.Context())
		if err != nil {
			app.error(w, err)
			return
		}
		assignedLeads = leads
	} else {
		leads, err := app.service.GetAssignedLeadsUser(r.Context(), user.Id)
		if err != nil {
			app.error(w, err)
			return
		}
		assignedLeads = leads
	}
	if user.Role == auth.AdminRole || user.Role == auth.LogistRole || user.Role == auth.RopRole {
		leads, err := app.service.GetInDeliveryLeads(r.Context())
		if err != nil {
			app.error(w, err)
			return
		}
		inDeliveryLeads = leads
	} else {
		leads, err := app.service.GetInDeliveryLeadsUser(r.Context(), user.Id)
		if err != nil {
			app.error(w, err)
			return
		}
		inDeliveryLeads = leads
	}
	if user.Role == auth.AdminRole || user.Role == auth.LogistRole || user.Role == auth.RopRole {
		leads, err := app.service.GetCompletedLeads(r.Context())
		if err != nil {
			app.error(w, err)
			return
		}
		completedLeads = leads
	} else {
		leads, err := app.service.GetCompletedLeadsUser(r.Context(), user.Id)
		if err != nil {
			app.error(w, err)
			return
		}
		completedLeads = leads
	}
	users, err := app.service.GetUsers(r.Context(), 0, pagesLimit)
	if err != nil {
		app.error(w, err)
		return
	}
	products, err := app.service.GetProducts(r.Context(), 0, pagesLimit)
	if err != nil {
		app.error(w, err)
		return
	}
	app.execute(w, tLeads, "/pages/leads", layoutData{
		BarsData: barsData{
			Page:     "leads",
			PageName: "Лидтер",
			Pages:    getPage(r),
		},
		User: app.service.GetAuthUser(r),
		Data: leadsContent{
			New:           newLeads,
			Assigned:      assignedLeads,
			InDelivery:    inDeliveryLeads,
			Completed:     completedLeads,
			Users:         users,
			Products:      products,
			Role:          string(user.Role),
			NewLeadsCount: len(newLeads),
		},
	})
}

func (app *app) handleLeadsNewGet(w http.ResponseWriter, r *http.Request) {
	newLeads, err := app.service.GetNewLeads(r.Context())
	if err != nil {
		app.error(w, err)
		return
	}
	resp, err := json.Marshal(newLeads)
	if err != nil {
		app.error(w, err)
	}
	_, err = w.Write(resp)
	if err != nil {
		app.error(w, err)
	}
}

func (app *app) handleLeadsAssignedGet(w http.ResponseWriter, r *http.Request) {
	assignedLeads, err := app.service.GetAssignedLeads(r.Context())
	if err != nil {
		app.error(w, err)
		return
	}
	resp, err := json.Marshal(assignedLeads)
	if err != nil {
		app.error(w, err)
	}
	_, err = w.Write(resp)
	if err != nil {
		app.error(w, err)
	}
}

func (app *app) handleLeadsNew(w http.ResponseWriter, _ *http.Request) {
	app.execute(w, tLeadsNewForm, "", nil)
}

func (app *app) handleLeadsPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.error(w, err)
		return
	}
	phone := r.FormValue("phone")
	lead, err := app.service.InsertLead(r.Context(), phone)
	if err != nil {
		if errors.Is(err, service.ErrInternal) {
			app.error(w, err)
			return
		}
		app.errorHx(w, tAlert, leadCellErrors, ErrChooseUser.Error())
		return
	}
	if r.Header.Get(acceptHeader) == jsonContentType {
		resp, err := json.Marshal(lead)
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
	users, err := app.service.GetUsers(r.Context(), 0, pagesLimit)
	if err != nil {
		app.error(w, err)
		return
	}
	app.execute(w, tLeadsNewCell, "", struct {
		Id    int
		Phone string
		Users []service.User
	}{
		Id:    lead.Id,
		Phone: lead.Phone,
		Users: users,
	})
}

func (app *app) handleLeadsAssign(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.error(w, err)
		return
	}
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		app.errorHx(w, tAlert, leadCellErrors, ErrChooseUser.Error())
		return
	}
	userIdStr := r.FormValue("user_id")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		app.errorHx(w, tAlert, leadCellErrors, ErrChooseUser.Error())
		return
	}
	lead, err := app.service.AssignLead(r.Context(), id, userId)
	if err != nil {
		app.error(w, err)
	}
	user, err := app.service.GetUser(r.Context(), lead.UserId)
	if err != nil {
		app.error(w, err)
	}
	lead.UserName = user.Name
	if r.Header.Get(acceptHeader) == jsonContentType {
		resp, err := json.Marshal(lead)
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
	w.Header().Add("HX-Trigger", fmt.Sprintf("lead-cell-%d", lead.Id))
	app.execute(w, tLeadsAssignedCell, "", lead)
}

func (app *app) handleLeadsComplete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		app.error(w, fmt.Errorf("%w: %s", service.ErrInvalidId, idStr))
		return
	}
	err = app.service.CompleteLead(r.Context(), id)
	if err != nil {
		app.error(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (app *app) handleLeadsSell(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.error(w, err)
		return
	}
	idStr := r.FormValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		app.error(w, fmt.Errorf("%w: %s", service.ErrInvalidId, idStr))
		return
	}
	name := r.FormValue("name")
	address := r.FormValue("address")
	saletype := r.FormValue("saletype")
	deliveryType := r.FormValue("delivery-type")
	paymentAtStr := r.FormValue("payment-at")
	paymentAt, err := time.Parse("2006-01-02T15:04", paymentAtStr)
	if err != nil {
		app.error(w, fmt.Errorf("error parsing payment at: %s: %w", paymentAtStr, err))
		return
	}
	deliveryCostStr := r.FormValue("deliverycost")
	deliveryCost, err := strconv.ParseFloat(deliveryCostStr, 32)
	if err != nil {
		app.errorHx(w, tText, "#product-modal-errors", "жеткізу құны сан болуы тиіс")
		return
	}
	loanCostStr := r.FormValue("loancost")
	loanCost, err := strconv.ParseFloat(loanCostStr, 32)
	if err != nil {
		app.errorHx(w, tText, "#leads-modal-errors", "несие құны сан болуы тиіс")
		return
	}
	fullSumStr := r.FormValue("fullsum")
	fullSum, err := strconv.ParseFloat(fullSumStr, 32)
	if err != nil {
		app.error(w, fmt.Errorf("full sum isn't number: %s", fullSumStr))
		return
	}
	itemsSumStr := r.FormValue("itemsum")
	itemsSum, err := strconv.ParseFloat(itemsSumStr, 32)
	if err != nil {
		app.error(w, fmt.Errorf("items sum isn't number: %s", fullSumStr))
		return
	}
	items, err := app.parseCartItems(r.Context(), r.FormValue("items"))
	if err != nil {
		app.error(w, err)
		return
	}
	lead, err := app.service.SellLead(r.Context(), service.SellLeadParams{
		Id:           id,
		Name:         name,
		Address:      address,
		Type:         saletype,
		DeliveryCost: float32(deliveryCost),
		LoanCost:     float32(loanCost),
		FullSum:      float32(fullSum),
		ItemsSum:     float32(itemsSum),
		Items:        items,
		DeliveryType: deliveryType,
		PaymentAt:    paymentAt,
	})
	if err != nil {
		app.error(w, err)
		return
	}
	w.Header().Add("HX-Trigger", "closeModal")
	w.Header().Add("HX-Trigger", fmt.Sprintf("lead-cell-%d", lead.Id))
	app.execute(w, tLeadsInDeliveryCell, "", lead)
}

func (app *app) parseCartItems(ctx context.Context, itemsStr string) ([]service.SaleItem, error) {
	itemsSplit := strings.Split(itemsStr, ";")
	items := make([]service.SaleItem, 0, len(itemsSplit))
	if len(itemsSplit) == 0 {
		return nil, fmt.Errorf("%w: %s", ErrInvalidCartItems, itemsStr)
	}
	for _, itemStr := range itemsSplit {
		if itemStr == "" {
			continue
		}
		temp := strings.Split(itemStr, ",")
		if len(temp) != 2 {
			return nil, fmt.Errorf("%w: %s", ErrInvalidCartItems, itemsStr)
		}
		productIdStr, quantityStr := temp[0], temp[1]
		productId, err := strconv.Atoi(productIdStr)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", ErrInvalidCartItems, itemsStr)
		}
		quantity, err := strconv.Atoi(quantityStr)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", ErrInvalidCartItems, itemsStr)
		}
		product, err := app.service.GetProduct(ctx, productId)
		if err != nil {
			return nil, err
		}
		items = append(items, service.SaleItem{
			Id:          productId,
			Price:       float32(product.Price) * float32(quantity),
			ProductName: product.Name,
			Quantity:    quantity,
			ProductId:   productId,
			SaleCount:   product.SaleCount * quantity,
		})
	}
	return items, nil
}

func (app *app) handleLeadsProduct(w http.ResponseWriter, r *http.Request) {
	productIdStr := r.FormValue("select-product")
	productId, err := strconv.Atoi(productIdStr)
	if err != nil {
		app.error(w, fmt.Errorf("%w: %s", service.ErrInvalidId, productIdStr))
		return
	}
	p, err := app.service.GetProduct(r.Context(), productId)
	app.execute(w, tLeadsProduct, "", p)
}
