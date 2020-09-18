package main

import (
	"github.com/call-me-snake/user_balance_service/internal/model"
	"github.com/call-me-snake/user_balance_service/internal/server"
	"github.com/call-me-snake/user_balance_service/internal/storage"
	"github.com/labstack/gommon/log"
)

func main() {
	log.Print("Started")
	//Устанавливаем значения переменных окружения
	config, err := model.Init()
	if err != nil {
		log.Print(err.Error())
		return
	}
	//Программа позволяет получать данные из двух разных мест. Но в данном случае нам не нужно это усложнение
	accSt, err := storage.New(config.AccountStorageConn)
	//logSt,err := storage.New(config.TransactionLoggerStorageConn)
	if err != nil {
		log.Error(err)
	}
	s := server.New(config.ServerAddress)
	err = s.Start(accSt, accSt)
	log.Print(err.Error())
}
