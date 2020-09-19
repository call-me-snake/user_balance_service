package main

import (
	"fmt"

	"github.com/call-me-snake/user_balance_service/internal/server"
	"github.com/call-me-snake/user_balance_service/internal/storage"
	"github.com/jessevdk/go-flags"
	"github.com/labstack/gommon/log"
)

//envs получает переменные окружения
type envs struct {
	ServerAddress                string `long:"http" env:"SERVER" description:"address of microservice" default:":8000"`
	AccountStorageConn           string `long:"accstconn" env:"ACC_STORAGE" description:"Connection string to account storage database" default:"user=postgres password=example dbname=accounts sslmode=disable port=5432 host=localhost"`
	TransactionLoggerStorageConn string `long:"logstconn" env:"LOGGER_STORAGE" description:"Connection string to transaction logger storage database" default:"user=postgres password=example dbname=accounts sslmode=disable port=5432 host=localhost"`
}

//Config хранит переменные окружения
type config struct {
	ServerAddress                string
	AccountStorageConn           string
	TransactionLoggerStorageConn string
}

//initConfig - получает переменные окружения с помощью envs (в перспективе может осуществлять их проверку)
func initConfig() (config, error) {
	e := envs{}
	c := config{}
	var err error
	parser := flags.NewParser(&e, flags.Default)
	if _, err = parser.Parse(); err != nil {
		return c, fmt.Errorf("Init: %v", err)
	}
	c.ServerAddress = e.ServerAddress
	c.AccountStorageConn = e.AccountStorageConn
	c.TransactionLoggerStorageConn = e.TransactionLoggerStorageConn
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
