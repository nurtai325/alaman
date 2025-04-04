package main

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/nurtai325/alaman/internal/config"
	"github.com/nurtai325/alaman/internal/db"
	"github.com/nurtai325/alaman/internal/db/repository"
)

func Some1() {
	conf, err := config.New()
	if err != nil {
		panic(err)
	}
	pool, err := db.New(conf)
	if err != nil {
		panic(err)
	}
	defer pool.Close()
	q := repository.New(pool)
	needed := []string{
		"+77478778628",
		"+77018638855",
		"+77024635948",
		"+77782148869",
		"+77018338948",
		"+77018008079",
		"+77753536757",
	}
	leads, err := q.GetNewLeads(context.Background(), repository.GetNewLeadsParams{
		Offset: 0,
		Limit:  100000,
	})
	if err != nil {
		panic(err)
	}
	j := 0
Loop:
	for _, lead := range leads {
		for _, newLead := range needed {
			if lead.Phone == newLead {
				continue Loop
			}
		}
		_, err := q.AssignLead(context.Background(), repository.AssignLeadParams{
			ID:     lead.ID,
			UserID: pgtype.Int4{Valid: true, Int32: 7},
		})
		if err != nil {
			panic(err)
		}
		j++
		fmt.Printf("%d. would assign: %s\n", j, lead.Phone)
	}
}

func main() {
	conf, err := config.New()
	if err != nil {
		panic(err)
	}
	pool, err := db.New(conf)
	if err != nil {
		panic(err)
	}
	defer pool.Close()
	q := repository.New(pool)
	rows, err := q.GetAssignedLeads(context.Background(), repository.GetAssignedLeadsParams{
		Offset: 0,
		Limit:  100000,
	})
	if err != nil {
		panic(err)
	}
	rowsLen := len(rows)
	fmt.Printf("got %d leads from the database\n", rowsLen)
	falseFound := 0
	i := rowsLen - 1 - 400
	for sent := 0; sent < 500; {
		fmt.Println(rows[i].Phone)
		if rows[i].Phone != "+77777777777" {
			sent++
			i--
		} else {
			falseFound++
		}
	}
	fmt.Printf("false found count: %d\n", falseFound)
}

type Lead struct {
	Id                 int          `json:"id"`
	Name               string       `json:"name"`
	UserName           string       `json:"user_name"`
	Address            string       `json:"address"`
	Phone              string       `json:"phone"`
	Completed          bool         `json:"completed"`
	SaleType           saleType     `json:"sale_type"`
	FullPrice          float32      `json:"full_price"`
	DeliveryCost       float32      `json:"delivery_cost"`
	LoanCost           float32      `json:"loan_cost"`
	Items              []SaleItem   `json:"items"`
	UserId             int          `json:"user_id"`
	SaleId             int          `json:"sale_id"`
	DeliveryType       deliveryType `json:"delivery_type"`
	DeliveryTypeName   string       `json:"delivery_type_name"`
	PaymentAt          time.Time    `json:"payment_at"`
	PaymentAtFormatted string       `json:"payment_at_formatted"`
	CreatedAt          time.Time    `json:"created_at"`
	Page               int          `json:"page"`
}

type SaleItem struct {
	Id          int     `json:"id"`
	ProductName string  `json:"product_name"`
	ProductId   int     `json:"product_id"`
	Price       float32 `json:"price"`
	Quantity    int     `json:"quantity"`
	SaleCount   int     `json:"sale_count"`
}

func Some() {
	conf, err := config.New()
	if err != nil {
		panic(err)
	}
	pool, err := db.New(conf)
	if err != nil {
		panic(err)
	}
	defer pool.Close()
	q := repository.New(pool)
	fullLead, err := q.GetFullLead(context.Background(), 5002)
	if err != nil {
		panic(err)
	}
	items, err := q.GetSaleItems(context.Background(), fullLead.SaleID.Int32)
	sItems := make([]SaleItem, 0, len(items))
	itemsStr := ""
	itemSum := 0
	for _, item := range items {
		sItems = append(sItems, SaleItem{
			Id:          int(item.ID),
			ProductName: item.ProductName,
			Quantity:    int(item.Quantity),
			ProductId:   int(item.ProductID),
		})
		itemsStr += fmt.Sprintf("%d %s\n", item.Quantity, item.ProductName)
		itemSum += int(item.Price)
	}
	lead := Lead{
		Id:                 int(fullLead.ID),
		Phone:              fullLead.Phone,
		Address:            fullLead.Address.String,
		Name:               fullLead.Name.String,
		UserName:           fullLead.UserName,
		FullPrice:          fullLead.FullSum,
		DeliveryCost:       fullLead.DeliveryCost,
		LoanCost:           fullLead.LoanCost,
		SaleType:           saleType(fullLead.SaleType),
		Items:              sItems,
		PaymentAt:          fullLead.PaymentAt.Time,
		PaymentAtFormatted: fullLead.PaymentAt.Time.Format("2006/01/02 15:04"),
		DeliveryType:       deliveryType(fullLead.DeliveryType.String),
		DeliveryTypeName:   getDeliveryTypeName(deliveryType(fullLead.DeliveryType.String)),
	}
	lastLine := fmt.Sprintf("%s: %d", getSaleTypeName(lead.SaleType), int(itemSum))
	if lead.DeliveryType != noDelivery {
		lastLine += fmt.Sprintf(" + %d", int(lead.DeliveryCost))
	}
	if lead.SaleType == kaspiLoan {
		lastLine += fmt.Sprintf(" + %d", int(lead.LoanCost))
	}
	lastLine += fmt.Sprintf(" = %d", int(fullLead.FullSum))
	msg := fmt.Sprintf(`Кеңесші маман: %s
Аты: %s
Номер: %s
Адрес: %s
Төлем түрі: %s
Жеткізу түрі: %s
Төлем уақыты: %s
%s
%s
`, lead.UserName, lead.Name, lead.Phone, lead.Address, getSaleTypeName(lead.SaleType), getDeliveryTypeName(lead.DeliveryType), lead.PaymentAt.Format("2006/01/02 15:04"), itemsStr, lastLine)
	fmt.Println(msg)
}

type saleType string
type deliveryType string

const (
	kaspiLoan     saleType = "kaspi-loan"
	cash          saleType = "cash"
	kaspiRed      saleType = "red"
	kaspiTransfer saleType = "kaspi-transfer"
	kaspiQr       saleType = "kaspi-qr"

	noDelivery deliveryType = "no"
	mail       deliveryType = "mail"
	train      deliveryType = "train"
	taxi       deliveryType = "taxi"
)

func validSaleType(saletypeStr string) bool {
	saletype := saleType(saletypeStr)
	switch saletype {
	case kaspiLoan:
		return true
	case cash:
		return true
	case kaspiRed:
		return true
	case kaspiTransfer:
		return true
	case kaspiQr:
		return true
	default:
		return false
	}
}

func validDeliveryType(deliveryTypeStr string) bool {
	deliveryType := deliveryType(deliveryTypeStr)
	switch deliveryType {
	case noDelivery:
		return true
	case mail:
		return true
	case train:
		return true
	case taxi:
		return true
	default:
		return false
	}
}

func getDeliveryTypeName(dType deliveryType) string {
	switch dType {
	case noDelivery:
		return "жоқ"
	case mail:
		return "почта"
	case train:
		return "пойыз"
	case taxi:
		return "такси"
	default:
		return ""
	}
}

func getSaleTypeName(val saleType) string {
	switch val {
	case kaspiLoan:
		return "бөліп төлеу"
	case cash:
		return "қолма-қол"
	case kaspiRed:
		return "kaspi red"
	case kaspiTransfer:
		return "kaspi аударым"
	case kaspiQr:
		return "kaspi qr"
	}
	return ""
}
