package main

import (
	"context"

	"github.com/nurtai325/alaman/internal/config"
	"github.com/nurtai325/alaman/internal/db"
	"github.com/nurtai325/alaman/internal/db/repository"
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
	q.InsertLead(context.Background(), "+77777777777")
}
