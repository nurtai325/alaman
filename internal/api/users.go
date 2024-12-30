package api

import (
	"net/http"
	"strconv"

	"github.com/nurtai325/alaman/internal/auth"
)

var userColumns = []string{"Аты", "Номер", "Статус", "Енгізілді"}

func (app *app) handleUsersGet(w http.ResponseWriter, r *http.Request) {
	pageQ := r.URL.Query().Get("page")
	page := 1
	if pageQ != "" {
		converted, err := strconv.Atoi(pageQ)
		if err != nil {
			app.execute(w, tLayout, layoutData{
				Page:  pUsers,
				User:  auth.GetUser(r),
				Pages: pages,
				Error: ErrPageNotFound.Error(),
				TableData: tableData{
					Error: ErrPageNotFound.Error(),
				},
			})
			return
		}
		page = converted
	}
	minId := pageOffset*(page-1) + 1
	maxId := minId + pageOffset - 1
	users, err := app.service.GetUsers(r.Context(), minId, maxId)
	if err != nil {
		app.error(w, err)
		return
	}
	rows := make([]row, 0, len(users))
	for _, user := range users {
		rows = append(rows, row{
			Id: int(user.ID),
			Cells: []cell{
				{inputCell, user.Name},
				{inputCell, user.Phone},
				{inputCell, user.Status.String},
				{inputCell, user.CreatedAt.Time.String()},
			},
		})
	}
	app.execute(w, tLayout, layoutData{
		Page:  pUsers,
		User:  auth.GetUser(r),
		Pages: pages,
		TableData: tableData{
			Resource: pUsers.Slug,
			Columns:  userColumns,
			Rows:     rows,
			Page:     page,
		},
	})
}
