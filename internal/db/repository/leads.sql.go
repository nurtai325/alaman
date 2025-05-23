// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: leads.sql

package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const assignLead = `-- name: AssignLead :one
UPDATE leads
SET user_id = $2
WHERE id = $1
RETURNING id, name, address, phone, completed, user_id, sale_id, created_at, sold_at, first_photo, second_photo
`

type AssignLeadParams struct {
	ID     int32
	UserID pgtype.Int4
}

func (q *Queries) AssignLead(ctx context.Context, arg AssignLeadParams) (Lead, error) {
	row := q.db.QueryRow(ctx, assignLead, arg.ID, arg.UserID)
	var i Lead
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Address,
		&i.Phone,
		&i.Completed,
		&i.UserID,
		&i.SaleID,
		&i.CreatedAt,
		&i.SoldAt,
		&i.FirstPhoto,
		&i.SecondPhoto,
	)
	return i, err
}

const completeLead = `-- name: CompleteLead :one
UPDATE leads
SET completed = true, first_photo = $2, second_photo = $3
WHERE id = $1
RETURNING id, name, address, phone, completed, user_id, sale_id, created_at, sold_at, first_photo, second_photo
`

type CompleteLeadParams struct {
	ID          int32
	FirstPhoto  string
	SecondPhoto string
}

func (q *Queries) CompleteLead(ctx context.Context, arg CompleteLeadParams) (Lead, error) {
	row := q.db.QueryRow(ctx, completeLead, arg.ID, arg.FirstPhoto, arg.SecondPhoto)
	var i Lead
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Address,
		&i.Phone,
		&i.Completed,
		&i.UserID,
		&i.SaleID,
		&i.CreatedAt,
		&i.SoldAt,
		&i.FirstPhoto,
		&i.SecondPhoto,
	)
	return i, err
}

const getAssignedLeads = `-- name: GetAssignedLeads :many
SELECT l.id, l.name, l.address, l.phone, l.completed, l.user_id, l.sale_id, l.created_at, l.sold_at, l.first_photo, l.second_photo, u.name AS user_name FROM leads AS l
INNER JOIN users u ON l.user_id = u.id
WHERE user_id IS NOT NULL AND sale_id IS NULL
ORDER BY created_at DESC
LIMIT $2
OFFSET $1
`

type GetAssignedLeadsParams struct {
	Offset int64
	Limit  int64
}

type GetAssignedLeadsRow struct {
	ID          int32
	Name        pgtype.Text
	Address     pgtype.Text
	Phone       string
	Completed   bool
	UserID      pgtype.Int4
	SaleID      pgtype.Int4
	CreatedAt   pgtype.Timestamptz
	SoldAt      pgtype.Timestamptz
	FirstPhoto  string
	SecondPhoto string
	UserName    string
}

func (q *Queries) GetAssignedLeads(ctx context.Context, arg GetAssignedLeadsParams) ([]GetAssignedLeadsRow, error) {
	rows, err := q.db.Query(ctx, getAssignedLeads, arg.Offset, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAssignedLeadsRow
	for rows.Next() {
		var i GetAssignedLeadsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Address,
			&i.Phone,
			&i.Completed,
			&i.UserID,
			&i.SaleID,
			&i.CreatedAt,
			&i.SoldAt,
			&i.FirstPhoto,
			&i.SecondPhoto,
			&i.UserName,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAssignedLeadsByUser = `-- name: GetAssignedLeadsByUser :many
SELECT l.id, l.name, l.address, l.phone, l.completed, l.user_id, l.sale_id, l.created_at, l.sold_at, l.first_photo, l.second_photo, u.name AS user_name FROM leads AS l
INNER JOIN users u ON l.user_id = u.id
WHERE user_id IS NOT NULL AND sale_id IS NULL AND user_id = $1
ORDER BY created_at DESC
LIMIT $3
OFFSET $2
`

type GetAssignedLeadsByUserParams struct {
	UserID pgtype.Int4
	Offset int64
	Limit  int64
}

type GetAssignedLeadsByUserRow struct {
	ID          int32
	Name        pgtype.Text
	Address     pgtype.Text
	Phone       string
	Completed   bool
	UserID      pgtype.Int4
	SaleID      pgtype.Int4
	CreatedAt   pgtype.Timestamptz
	SoldAt      pgtype.Timestamptz
	FirstPhoto  string
	SecondPhoto string
	UserName    string
}

func (q *Queries) GetAssignedLeadsByUser(ctx context.Context, arg GetAssignedLeadsByUserParams) ([]GetAssignedLeadsByUserRow, error) {
	rows, err := q.db.Query(ctx, getAssignedLeadsByUser, arg.UserID, arg.Offset, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAssignedLeadsByUserRow
	for rows.Next() {
		var i GetAssignedLeadsByUserRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Address,
			&i.Phone,
			&i.Completed,
			&i.UserID,
			&i.SaleID,
			&i.CreatedAt,
			&i.SoldAt,
			&i.FirstPhoto,
			&i.SecondPhoto,
			&i.UserName,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAssignedLeadsSearch = `-- name: GetAssignedLeadsSearch :many
SELECT l.id, l.name, l.address, l.phone, l.completed, l.user_id, l.sale_id, l.created_at, l.sold_at, l.first_photo, l.second_photo, u.name AS user_name FROM leads AS l
INNER JOIN users u ON l.user_id = u.id
WHERE user_id IS NOT NULL AND sale_id IS NULL AND l.phone LIKE $1
ORDER BY created_at DESC
LIMIT 9
`

type GetAssignedLeadsSearchRow struct {
	ID          int32
	Name        pgtype.Text
	Address     pgtype.Text
	Phone       string
	Completed   bool
	UserID      pgtype.Int4
	SaleID      pgtype.Int4
	CreatedAt   pgtype.Timestamptz
	SoldAt      pgtype.Timestamptz
	FirstPhoto  string
	SecondPhoto string
	UserName    string
}

func (q *Queries) GetAssignedLeadsSearch(ctx context.Context, phone string) ([]GetAssignedLeadsSearchRow, error) {
	rows, err := q.db.Query(ctx, getAssignedLeadsSearch, phone)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAssignedLeadsSearchRow
	for rows.Next() {
		var i GetAssignedLeadsSearchRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Address,
			&i.Phone,
			&i.Completed,
			&i.UserID,
			&i.SaleID,
			&i.CreatedAt,
			&i.SoldAt,
			&i.FirstPhoto,
			&i.SecondPhoto,
			&i.UserName,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getCompletedLeads = `-- name: GetCompletedLeads :many
SELECT l.id, l.name, l.address, l.phone, l.completed, l.user_id, l.sale_id, l.created_at, l.sold_at, l.first_photo, l.second_photo, u.name AS user_name, s.full_sum, s.delivery_type, s.payment_at FROM leads AS l
INNER JOIN users u ON l.user_id = u.id
INNER JOIN sales s ON l.sale_id = s.id
WHERE user_id IS NOT NULL AND sale_id IS NOT NULL AND completed = true
ORDER BY sold_at DESC
LIMIT $2
OFFSET $1
`

type GetCompletedLeadsParams struct {
	Offset int64
	Limit  int64
}

type GetCompletedLeadsRow struct {
	ID           int32
	Name         pgtype.Text
	Address      pgtype.Text
	Phone        string
	Completed    bool
	UserID       pgtype.Int4
	SaleID       pgtype.Int4
	CreatedAt    pgtype.Timestamptz
	SoldAt       pgtype.Timestamptz
	FirstPhoto   string
	SecondPhoto  string
	UserName     string
	FullSum      float32
	DeliveryType pgtype.Text
	PaymentAt    pgtype.Timestamptz
}

func (q *Queries) GetCompletedLeads(ctx context.Context, arg GetCompletedLeadsParams) ([]GetCompletedLeadsRow, error) {
	rows, err := q.db.Query(ctx, getCompletedLeads, arg.Offset, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetCompletedLeadsRow
	for rows.Next() {
		var i GetCompletedLeadsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Address,
			&i.Phone,
			&i.Completed,
			&i.UserID,
			&i.SaleID,
			&i.CreatedAt,
			&i.SoldAt,
			&i.FirstPhoto,
			&i.SecondPhoto,
			&i.UserName,
			&i.FullSum,
			&i.DeliveryType,
			&i.PaymentAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getCompletedLeadsByUser = `-- name: GetCompletedLeadsByUser :many
SELECT l.id, l.name, l.address, l.phone, l.completed, l.user_id, l.sale_id, l.created_at, l.sold_at, l.first_photo, l.second_photo, u.name AS user_name, s.full_sum, s.delivery_type, s.payment_at FROM leads AS l
INNER JOIN users u ON l.user_id = u.id
INNER JOIN sales s ON l.sale_id = s.id
WHERE user_id IS NOT NULL AND sale_id IS NOT NULL AND completed = true AND user_id = $1
ORDER BY sold_at DESC
LIMIT $3
OFFSET $2
`

type GetCompletedLeadsByUserParams struct {
	UserID pgtype.Int4
	Offset int64
	Limit  int64
}

type GetCompletedLeadsByUserRow struct {
	ID           int32
	Name         pgtype.Text
	Address      pgtype.Text
	Phone        string
	Completed    bool
	UserID       pgtype.Int4
	SaleID       pgtype.Int4
	CreatedAt    pgtype.Timestamptz
	SoldAt       pgtype.Timestamptz
	FirstPhoto   string
	SecondPhoto  string
	UserName     string
	FullSum      float32
	DeliveryType pgtype.Text
	PaymentAt    pgtype.Timestamptz
}

func (q *Queries) GetCompletedLeadsByUser(ctx context.Context, arg GetCompletedLeadsByUserParams) ([]GetCompletedLeadsByUserRow, error) {
	rows, err := q.db.Query(ctx, getCompletedLeadsByUser, arg.UserID, arg.Offset, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetCompletedLeadsByUserRow
	for rows.Next() {
		var i GetCompletedLeadsByUserRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Address,
			&i.Phone,
			&i.Completed,
			&i.UserID,
			&i.SaleID,
			&i.CreatedAt,
			&i.SoldAt,
			&i.FirstPhoto,
			&i.SecondPhoto,
			&i.UserName,
			&i.FullSum,
			&i.DeliveryType,
			&i.PaymentAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getCompletedLeadsSearch = `-- name: GetCompletedLeadsSearch :many
SELECT l.id, l.name, l.address, l.phone, l.completed, l.user_id, l.sale_id, l.created_at, l.sold_at, l.first_photo, l.second_photo, u.name AS user_name, s.full_sum, s.delivery_type, s.payment_at FROM leads AS l
INNER JOIN users u ON l.user_id = u.id
INNER JOIN sales s ON l.sale_id = s.id
WHERE user_id IS NOT NULL AND sale_id IS NOT NULL AND completed = true AND l.phone LIKE $1
ORDER BY sold_at DESC
LIMIT 9
`

type GetCompletedLeadsSearchRow struct {
	ID           int32
	Name         pgtype.Text
	Address      pgtype.Text
	Phone        string
	Completed    bool
	UserID       pgtype.Int4
	SaleID       pgtype.Int4
	CreatedAt    pgtype.Timestamptz
	SoldAt       pgtype.Timestamptz
	FirstPhoto   string
	SecondPhoto  string
	UserName     string
	FullSum      float32
	DeliveryType pgtype.Text
	PaymentAt    pgtype.Timestamptz
}

func (q *Queries) GetCompletedLeadsSearch(ctx context.Context, phone string) ([]GetCompletedLeadsSearchRow, error) {
	rows, err := q.db.Query(ctx, getCompletedLeadsSearch, phone)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetCompletedLeadsSearchRow
	for rows.Next() {
		var i GetCompletedLeadsSearchRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Address,
			&i.Phone,
			&i.Completed,
			&i.UserID,
			&i.SaleID,
			&i.CreatedAt,
			&i.SoldAt,
			&i.FirstPhoto,
			&i.SecondPhoto,
			&i.UserName,
			&i.FullSum,
			&i.DeliveryType,
			&i.PaymentAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getFullLead = `-- name: GetFullLead :one
SELECT l.id, l.name, l.address, l.phone, l.completed, l.user_id, l.sale_id, l.created_at, l.sold_at, l.first_photo, l.second_photo, u.name AS user_name, u.phone AS user_phone, s.full_sum, s.delivery_cost, s.loan_cost, s.delivery_type, s.payment_at, s.type AS sale_type
FROM leads AS l
INNER JOIN users u ON l.user_id = u.id
INNER JOIN sales s ON l.sale_id = s.id
WHERE l.id = $1
LIMIT 1
`

type GetFullLeadRow struct {
	ID           int32
	Name         pgtype.Text
	Address      pgtype.Text
	Phone        string
	Completed    bool
	UserID       pgtype.Int4
	SaleID       pgtype.Int4
	CreatedAt    pgtype.Timestamptz
	SoldAt       pgtype.Timestamptz
	FirstPhoto   string
	SecondPhoto  string
	UserName     string
	UserPhone    string
	FullSum      float32
	DeliveryCost float32
	LoanCost     float32
	DeliveryType pgtype.Text
	PaymentAt    pgtype.Timestamptz
	SaleType     string
}

func (q *Queries) GetFullLead(ctx context.Context, id int32) (GetFullLeadRow, error) {
	row := q.db.QueryRow(ctx, getFullLead, id)
	var i GetFullLeadRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Address,
		&i.Phone,
		&i.Completed,
		&i.UserID,
		&i.SaleID,
		&i.CreatedAt,
		&i.SoldAt,
		&i.FirstPhoto,
		&i.SecondPhoto,
		&i.UserName,
		&i.UserPhone,
		&i.FullSum,
		&i.DeliveryCost,
		&i.LoanCost,
		&i.DeliveryType,
		&i.PaymentAt,
		&i.SaleType,
	)
	return i, err
}

const getInDeliveryLeads = `-- name: GetInDeliveryLeads :many
SELECT l.id, l.name, l.address, l.phone, l.completed, l.user_id, l.sale_id, l.created_at, l.sold_at, l.first_photo, l.second_photo, u.name AS user_name, s.full_sum, s.delivery_type, s.payment_at FROM leads AS l
INNER JOIN users u ON l.user_id = u.id
INNER JOIN sales s ON l.sale_id = s.id
WHERE user_id IS NOT NULL AND sale_id IS NOT NULL AND completed = false
ORDER BY sold_at DESC
LIMIT $2
OFFSET $1
`

type GetInDeliveryLeadsParams struct {
	Offset int64
	Limit  int64
}

type GetInDeliveryLeadsRow struct {
	ID           int32
	Name         pgtype.Text
	Address      pgtype.Text
	Phone        string
	Completed    bool
	UserID       pgtype.Int4
	SaleID       pgtype.Int4
	CreatedAt    pgtype.Timestamptz
	SoldAt       pgtype.Timestamptz
	FirstPhoto   string
	SecondPhoto  string
	UserName     string
	FullSum      float32
	DeliveryType pgtype.Text
	PaymentAt    pgtype.Timestamptz
}

func (q *Queries) GetInDeliveryLeads(ctx context.Context, arg GetInDeliveryLeadsParams) ([]GetInDeliveryLeadsRow, error) {
	rows, err := q.db.Query(ctx, getInDeliveryLeads, arg.Offset, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetInDeliveryLeadsRow
	for rows.Next() {
		var i GetInDeliveryLeadsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Address,
			&i.Phone,
			&i.Completed,
			&i.UserID,
			&i.SaleID,
			&i.CreatedAt,
			&i.SoldAt,
			&i.FirstPhoto,
			&i.SecondPhoto,
			&i.UserName,
			&i.FullSum,
			&i.DeliveryType,
			&i.PaymentAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getInDeliveryLeadsByUser = `-- name: GetInDeliveryLeadsByUser :many
SELECT l.id, l.name, l.address, l.phone, l.completed, l.user_id, l.sale_id, l.created_at, l.sold_at, l.first_photo, l.second_photo, u.name AS user_name, s.full_sum, s.delivery_type, s.payment_at FROM leads AS l
INNER JOIN users u ON l.user_id = u.id
INNER JOIN sales s ON l.sale_id = s.id
WHERE user_id IS NOT NULL AND sale_id IS NOT NULL AND completed = false AND user_id = $1
ORDER BY sold_at DESC
LIMIT $3
OFFSET $2
`

type GetInDeliveryLeadsByUserParams struct {
	UserID pgtype.Int4
	Offset int64
	Limit  int64
}

type GetInDeliveryLeadsByUserRow struct {
	ID           int32
	Name         pgtype.Text
	Address      pgtype.Text
	Phone        string
	Completed    bool
	UserID       pgtype.Int4
	SaleID       pgtype.Int4
	CreatedAt    pgtype.Timestamptz
	SoldAt       pgtype.Timestamptz
	FirstPhoto   string
	SecondPhoto  string
	UserName     string
	FullSum      float32
	DeliveryType pgtype.Text
	PaymentAt    pgtype.Timestamptz
}

func (q *Queries) GetInDeliveryLeadsByUser(ctx context.Context, arg GetInDeliveryLeadsByUserParams) ([]GetInDeliveryLeadsByUserRow, error) {
	rows, err := q.db.Query(ctx, getInDeliveryLeadsByUser, arg.UserID, arg.Offset, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetInDeliveryLeadsByUserRow
	for rows.Next() {
		var i GetInDeliveryLeadsByUserRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Address,
			&i.Phone,
			&i.Completed,
			&i.UserID,
			&i.SaleID,
			&i.CreatedAt,
			&i.SoldAt,
			&i.FirstPhoto,
			&i.SecondPhoto,
			&i.UserName,
			&i.FullSum,
			&i.DeliveryType,
			&i.PaymentAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getInDeliveryLeadsSearch = `-- name: GetInDeliveryLeadsSearch :many
SELECT l.id, l.name, l.address, l.phone, l.completed, l.user_id, l.sale_id, l.created_at, l.sold_at, l.first_photo, l.second_photo, u.name AS user_name, s.full_sum, s.delivery_type, s.payment_at FROM leads AS l
INNER JOIN users u ON l.user_id = u.id
INNER JOIN sales s ON l.sale_id = s.id
WHERE user_id IS NOT NULL AND sale_id IS NOT NULL AND completed = false AND l.phone LIKE $1
ORDER BY sold_at DESC
LIMIT 9
`

type GetInDeliveryLeadsSearchRow struct {
	ID           int32
	Name         pgtype.Text
	Address      pgtype.Text
	Phone        string
	Completed    bool
	UserID       pgtype.Int4
	SaleID       pgtype.Int4
	CreatedAt    pgtype.Timestamptz
	SoldAt       pgtype.Timestamptz
	FirstPhoto   string
	SecondPhoto  string
	UserName     string
	FullSum      float32
	DeliveryType pgtype.Text
	PaymentAt    pgtype.Timestamptz
}

func (q *Queries) GetInDeliveryLeadsSearch(ctx context.Context, phone string) ([]GetInDeliveryLeadsSearchRow, error) {
	rows, err := q.db.Query(ctx, getInDeliveryLeadsSearch, phone)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetInDeliveryLeadsSearchRow
	for rows.Next() {
		var i GetInDeliveryLeadsSearchRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Address,
			&i.Phone,
			&i.Completed,
			&i.UserID,
			&i.SaleID,
			&i.CreatedAt,
			&i.SoldAt,
			&i.FirstPhoto,
			&i.SecondPhoto,
			&i.UserName,
			&i.FullSum,
			&i.DeliveryType,
			&i.PaymentAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getLead = `-- name: GetLead :one
SELECT id, name, address, phone, completed, user_id, sale_id, created_at, sold_at, first_photo, second_photo FROM leads 
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetLead(ctx context.Context, id int32) (Lead, error) {
	row := q.db.QueryRow(ctx, getLead, id)
	var i Lead
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Address,
		&i.Phone,
		&i.Completed,
		&i.UserID,
		&i.SaleID,
		&i.CreatedAt,
		&i.SoldAt,
		&i.FirstPhoto,
		&i.SecondPhoto,
	)
	return i, err
}

const getLeadByPhone = `-- name: GetLeadByPhone :one
SELECT id, name, address, phone, completed, user_id, sale_id, created_at, sold_at, first_photo, second_photo FROM leads 
WHERE phone = $1
ORDER BY created_at DESC
LIMIT 1
`

func (q *Queries) GetLeadByPhone(ctx context.Context, phone string) (Lead, error) {
	row := q.db.QueryRow(ctx, getLeadByPhone, phone)
	var i Lead
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Address,
		&i.Phone,
		&i.Completed,
		&i.UserID,
		&i.SaleID,
		&i.CreatedAt,
		&i.SoldAt,
		&i.FirstPhoto,
		&i.SecondPhoto,
	)
	return i, err
}

const getNewLeads = `-- name: GetNewLeads :many
SELECT id, name, address, phone, completed, user_id, sale_id, created_at, sold_at, first_photo, second_photo FROM leads AS l
WHERE user_id IS NULL
ORDER BY created_at DESC
LIMIT $2
OFFSET $1
`

type GetNewLeadsParams struct {
	Offset int64
	Limit  int64
}

func (q *Queries) GetNewLeads(ctx context.Context, arg GetNewLeadsParams) ([]Lead, error) {
	rows, err := q.db.Query(ctx, getNewLeads, arg.Offset, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Lead
	for rows.Next() {
		var i Lead
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Address,
			&i.Phone,
			&i.Completed,
			&i.UserID,
			&i.SaleID,
			&i.CreatedAt,
			&i.SoldAt,
			&i.FirstPhoto,
			&i.SecondPhoto,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getNewLeadsCount = `-- name: GetNewLeadsCount :one
SELECT COUNT(*) FROM leads AS l
WHERE user_id IS NULL
`

func (q *Queries) GetNewLeadsCount(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, getNewLeadsCount)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getNewLeadsSearch = `-- name: GetNewLeadsSearch :many
SELECT id, name, address, phone, completed, user_id, sale_id, created_at, sold_at, first_photo, second_photo FROM leads AS l
WHERE user_id IS NULL AND phone LIKE $1
ORDER BY created_at DESC
LIMIT 9
`

func (q *Queries) GetNewLeadsSearch(ctx context.Context, phone string) ([]Lead, error) {
	rows, err := q.db.Query(ctx, getNewLeadsSearch, phone)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Lead
	for rows.Next() {
		var i Lead
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Address,
			&i.Phone,
			&i.Completed,
			&i.UserID,
			&i.SaleID,
			&i.CreatedAt,
			&i.SoldAt,
			&i.FirstPhoto,
			&i.SecondPhoto,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getSaleItems = `-- name: GetSaleItems :many
SELECT id, price, product_name, sale_count, quantity, sale_id, product_id, created_at FROM sale_items AS s
WHERE s.sale_id = $1
`

func (q *Queries) GetSaleItems(ctx context.Context, saleID int32) ([]SaleItem, error) {
	rows, err := q.db.Query(ctx, getSaleItems, saleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SaleItem
	for rows.Next() {
		var i SaleItem
		if err := rows.Scan(
			&i.ID,
			&i.Price,
			&i.ProductName,
			&i.SaleCount,
			&i.Quantity,
			&i.SaleID,
			&i.ProductID,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const insertLead = `-- name: InsertLead :one
INSERT INTO leads(phone)
VALUES ($1)
RETURNING id, name, address, phone, completed, user_id, sale_id, created_at, sold_at, first_photo, second_photo
`

func (q *Queries) InsertLead(ctx context.Context, phone string) (Lead, error) {
	row := q.db.QueryRow(ctx, insertLead, phone)
	var i Lead
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Address,
		&i.Phone,
		&i.Completed,
		&i.UserID,
		&i.SaleID,
		&i.CreatedAt,
		&i.SoldAt,
		&i.FirstPhoto,
		&i.SecondPhoto,
	)
	return i, err
}

const insertSale = `-- name: InsertSale :one
INSERT INTO sales(type, full_sum, delivery_cost, loan_cost, items_sum, delivery_type, payment_at)
VALUES($1, $2, $3, $4, $5, $6, $7)
RETURNING id, type, delivery_type, payment_at, full_sum, delivery_cost, loan_cost, items_sum, created_at
`

type InsertSaleParams struct {
	Type         string
	FullSum      float32
	DeliveryCost float32
	LoanCost     float32
	ItemsSum     float32
	DeliveryType pgtype.Text
	PaymentAt    pgtype.Timestamptz
}

func (q *Queries) InsertSale(ctx context.Context, arg InsertSaleParams) (Sale, error) {
	row := q.db.QueryRow(ctx, insertSale,
		arg.Type,
		arg.FullSum,
		arg.DeliveryCost,
		arg.LoanCost,
		arg.ItemsSum,
		arg.DeliveryType,
		arg.PaymentAt,
	)
	var i Sale
	err := row.Scan(
		&i.ID,
		&i.Type,
		&i.DeliveryType,
		&i.PaymentAt,
		&i.FullSum,
		&i.DeliveryCost,
		&i.LoanCost,
		&i.ItemsSum,
		&i.CreatedAt,
	)
	return i, err
}

const insertSaleItem = `-- name: InsertSaleItem :one
INSERT INTO sale_items(price, product_name, sale_id, quantity, product_id, sale_count)
VALUES($1, $2, $3, $4, $5, $6)
RETURNING id, price, product_name, sale_count, quantity, sale_id, product_id, created_at
`

type InsertSaleItemParams struct {
	Price       float32
	ProductName string
	SaleID      int32
	Quantity    int32
	ProductID   int32
	SaleCount   int32
}

func (q *Queries) InsertSaleItem(ctx context.Context, arg InsertSaleItemParams) (SaleItem, error) {
	row := q.db.QueryRow(ctx, insertSaleItem,
		arg.Price,
		arg.ProductName,
		arg.SaleID,
		arg.Quantity,
		arg.ProductID,
		arg.SaleCount,
	)
	var i SaleItem
	err := row.Scan(
		&i.ID,
		&i.Price,
		&i.ProductName,
		&i.SaleCount,
		&i.Quantity,
		&i.SaleID,
		&i.ProductID,
		&i.CreatedAt,
	)
	return i, err
}

const sellLead = `-- name: SellLead :one
UPDAte leads
SET sale_id = $2, sold_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, name, address, phone, completed, user_id, sale_id, created_at, sold_at, first_photo, second_photo
`

type SellLeadParams struct {
	ID     int32
	SaleID pgtype.Int4
}

func (q *Queries) SellLead(ctx context.Context, arg SellLeadParams) (Lead, error) {
	row := q.db.QueryRow(ctx, sellLead, arg.ID, arg.SaleID)
	var i Lead
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Address,
		&i.Phone,
		&i.Completed,
		&i.UserID,
		&i.SaleID,
		&i.CreatedAt,
		&i.SoldAt,
		&i.FirstPhoto,
		&i.SecondPhoto,
	)
	return i, err
}

const setLeadInfo = `-- name: SetLeadInfo :one
UPDATE leads
SET name = $2, address = $3
WHERE id = $1
RETURNING id, name, address, phone, completed, user_id, sale_id, created_at, sold_at, first_photo, second_photo
`

type SetLeadInfoParams struct {
	ID      int32
	Name    pgtype.Text
	Address pgtype.Text
}

func (q *Queries) SetLeadInfo(ctx context.Context, arg SetLeadInfoParams) (Lead, error) {
	row := q.db.QueryRow(ctx, setLeadInfo, arg.ID, arg.Name, arg.Address)
	var i Lead
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Address,
		&i.Phone,
		&i.Completed,
		&i.UserID,
		&i.SaleID,
		&i.CreatedAt,
		&i.SoldAt,
		&i.FirstPhoto,
		&i.SecondPhoto,
	)
	return i, err
}
