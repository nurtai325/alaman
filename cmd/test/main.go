package main

import (
	"context"
	"fmt"

	"github.com/nurtai325/alaman/internal/config"
	"github.com/nurtai325/alaman/internal/db"
	"github.com/nurtai325/alaman/internal/db/repository"
	"github.com/nurtai325/alaman/internal/service"
)

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
	res, err := pool.Query(context.Background(), "select phone, count(*) from leads group by phone having count(*) > 1;")
	if err != nil {
		panic(err)
	}
	defer res.Close()
	var duplicates []service.Lead
	for res.Next() {
		var phone string
		var count int
		err := res.Scan(&phone, &count)
		if err != nil {
			panic(err)
		}
		duplicates = append(duplicates, service.Lead{
			Phone: phone,
		})
	}
	leads, err := q.GetNewLeads(context.Background(), repository.GetNewLeadsParams{
		Offset: 0,
		Limit:  100000000,
	})
	for _, lead := range leads {
		for _, duplicate := range duplicates {
			if lead.Phone == duplicate.Phone {
				fmt.Println("duplicate: ", lead.Phone)
			}
		}
	}
}
