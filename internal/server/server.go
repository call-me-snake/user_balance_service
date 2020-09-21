package server

import (
	"fmt"
	"net/http"

	"github.com/call-me-snake/user_balance_service/internal/model"
	"github.com/gorilla/mux"
)

type Connector struct {
	router  *mux.Router
	address string
}

//New - Конструктор *Connector
func New(addr string) *Connector {
	c := &Connector{}
	c.router = mux.NewRouter()
	c.address = addr
	return c
}

func (c *Connector) executeHandlers(accStorage model.IBalanceInfoStorage) {
	c.router.HandleFunc("/alive", aliveHandler).Methods("GET")
	c.router.HandleFunc("/account/balance/info/{id:[0-9]+}", accountBalanceById(accStorage)).Methods("GET")
	c.router.HandleFunc("/account/balance/change", changeAccountBalance(accStorage)).Methods("POST")
	c.router.HandleFunc("/account/balance/transfer", transferSum(accStorage)).Methods("POST")
	c.router.HandleFunc("/account/balance/history", transactionsHistory(accStorage)).Methods("POST")
}

//Start запуск http сервера
func (c *Connector) Start(accStorage model.IBalanceInfoStorage) error {
	c.executeHandlers(accStorage)
	err := http.ListenAndServe(c.address, c.router)
	return fmt.Errorf("server.Start: %v", err)
}
