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

func (c *Connector) executeHandlers(accStorage model.IAccountsStorage, logger model.ITransactionLogger) {
	c.router.HandleFunc("/alive", aliveHandler).Methods("GET")
	c.router.HandleFunc("/account/info/{id:[0-9]+}", accountById(accStorage)).Methods("GET")
	c.router.HandleFunc("/account/change-balance", changeAccBalance(accStorage, logger)).Methods("PUT")
	c.router.HandleFunc("/account/transfer", transferSum(accStorage, logger)).Methods("PUT")
}

//Start запуск http сервера
func (c *Connector) Start(accStorage model.IAccountsStorage, logger model.ITransactionLogger) error {
	c.executeHandlers(accStorage, logger)
	err := http.ListenAndServe(c.address, c.router)
	return fmt.Errorf("server.Start: %v", err)
}
