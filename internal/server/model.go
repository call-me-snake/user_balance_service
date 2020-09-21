package server

import (
	"encoding/json"
	"net/http"
)

type accountByIdResponse struct {
	Id       int     `json:"Id"`
	Balance  float64 `json:"Balance"`
	Currency string  `json:"Currency"`
}

type changeAccBalanceRequest struct {
	Id    int     `json:"Id"`
	Delta float64 `json:"Delta"`
}

type changeAccBalanceResponse struct {
	Message string `json:"Message"`
}

type transferSumRequest struct {
	Id1   int     `json:"Id1"`
	Id2   int     `json:"Id2"`
	Delta float64 `json:"Delta"`
}

type transferSumResponce struct {
	Message string `json:"Message"`
}

type transactionsHistoryRequest struct {
	Id           int    `json:"Id"`
	SortedBy     string `json:"SortedBy,omitempty"`
	SortedByDesc bool   `json:"SortedByDesc,omitempty"`
}

type errorResponce struct {
	Message string `json:"Message"`
	ErrCode int    `json:"ErrCode"`
}

func makeErrResponce(userMessage string, errCode int, w http.ResponseWriter) {
	res, _ := json.Marshal(errorResponce{Message: userMessage, ErrCode: errCode})
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(errCode)
	w.Write(res)
}
