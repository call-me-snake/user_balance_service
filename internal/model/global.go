package model

import (
	"fmt"

	"github.com/jessevdk/go-flags"
)

//envs - приватная структура, получает переменные окружения
type envs struct {
	ServerAddress                string `long:"http" env:"SERVER" description:"address of microservice" default:":8000"`
	AccountStorageConn           string `long:"accstconn" env:"ACC_STORAGE" description:"Connection string to account storage database" default:"user=postgres password=example dbname=accounts sslmode=disable port=5432 host=localhost"`
	TransactionLoggerStorageConn string `long:"logstconn" env:"LOGGER_STORAGE" description:"Connection string to transaction logger storage database" default:"user=postgres password=example dbname=accounts sslmode=disable port=5432 host=localhost"`
}

//Config - публичная структура, хранит переменные окружения
type Config struct {
	ServerAddress                string
	AccountStorageConn           string
	TransactionLoggerStorageConn string
}

//Init - получает переменные окружения с помощью приватной структуры envs и проверяет их
func Init() (Config, error) {
	e := envs{}
	c := Config{}
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
