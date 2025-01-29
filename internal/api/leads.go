package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/nurtai325/alaman/internal/service"
)

var (
	ErrChooseUser = errors.New("қызметкерді таңдаңыз")
)

const (
	leadCellErrors = "#lead-cell-errors"
)

type leadsContent struct {
	New        []service.Lead
	Assigned   []service.Lead
	InDelivery []service.Lead
	Completed  []service.Lead
	Users      []service.User
}

func (app *app) handleLeadsGet(w http.ResponseWriter, r *http.Request) {
	newLeads, err := app.service.GetNewLeads(r.Context())
	if err != nil {
		app.error(w, err)
		return
	}
	assignedLeads, err := app.service.GetAssignedLeads(r.Context())
	if err != nil {
		app.error(w, err)
		return
	}
	inDeliveryLeads, err := app.service.GetInDeliveryLeads(r.Context())
	if err != nil {
		app.error(w, err)
		return
	}
	completedLeads, err := app.service.GetCompletedLeads(r.Context())
	if err != nil {
		app.error(w, err)
		return
	}
	users, err := app.service.GetUsers(r.Context(), 0, pagesLimit)
	if err != nil {
		app.error(w, err)
		return
	}
	app.execute(w, tLeads, "/pages/leads", layoutData{
		BarsData: barsData{
			Page:     "leads",
			PageName: "Лидтер",
			Pages:    adminPages,
		},
		User: app.service.GetAuthUser(r),
		Data: leadsContent{
			New:        newLeads,
			Assigned:   assignedLeads,
			InDelivery: inDeliveryLeads,
			Completed:  completedLeads,
			Users:      users,
		},
	})
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
	w.Header().Add("HX-Trigger", fmt.Sprintf("lead-cell-%d", lead.Id))
	app.execute(w, tLeadsAssignedCell, "", lead)
}
