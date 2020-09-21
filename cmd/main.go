package main

import (
	"fmt"

	"github.com/call-me-snake/user_balance_service/internal/model"
	"github.com/call-me-snake/user_balance_service/internal/server"
	"github.com/call-me-snake/user_balance_service/internal/storage"
	"github.com/jessevdk/go-flags"
	"github.com/labstack/gommon/log"
)

//envs получает переменные окружения
type envs struct {
	ServerAddress      string `long:"http" env:"SERVER" description:"address of microservice" default:":8000"`
	AccountStorageConn string `long:"accstconn" env:"ACC_STORAGE" description:"Connection string to account storage database" default:"user=postgres password=example dbname=accounts sslmode=disable port=5432 host=localhost"`
}

//initConfig - получает переменные окружения с помощью envs (в перспективе может осуществлять их проверку)
func initConfig() (model.Config, error) {
	e := envs{}
	c := model.Config{}
	var err error
	parser := flags.NewParser(&e, flags.Default)
	if _, err = parser.Parse(); err != nil {
		return c, fmt.Errorf("Init: %v", err)
	}
	c.ServerAddress = e.ServerAddress
	c.AccountStorageConn = e.AccountStorageConn
	return c, nil
}

func main() {
	log.Print("Started")
	//Устанавливаем значения переменных окружения
	config, err := initConfig()
	if err != nil {
		log.Print(err.Error())
		return
	}
	//Подключаемся к бд
	accSt, err := storage.New(config.AccountStorageConn)
	if err != nil {
		log.Print(err.Error())
		return
	}
	//Разворачиваем сервер
	s := server.New(config.ServerAddress)
	err = s.Start(accSt)
	log.Print(err.Error())
}
