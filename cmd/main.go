package main

import (
	"github.com/call-me-snake/user_balance_service/internal/server"
	"github.com/call-me-snake/user_balance_service/internal/storage"
	"github.com/labstack/gommon/log"
)

func main() {
	s1, err := storage.New("user=postgres password=example dbname=accounts sslmode=disable port=5432 host=localhost")
	if err != nil {
		log.Error(err)
	}
	s := server.New(":8000")
	s.Start(s1, s1)

}
