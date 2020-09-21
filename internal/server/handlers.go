package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/call-me-snake/user_balance_service/internal/convert"
	"github.com/call-me-snake/user_balance_service/internal/model"
	"github.com/gorilla/mux"
	"github.com/labstack/gommon/log"
	"golang.org/x/exp/errors/fmt"
)

const badRequestMessage = "Некорректные входные данные"
const internalErrorMessage = "Внутренняя ошибка сервера"
const insufficientFundsMessage = "Недостаточно средств на счету"
const nullSumMessage = "Нулевая сумма пополнения"
const defaultCurrency = "RUB"

func aliveHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from balance service"))
}

//accountById - возврат информации об аккаунте
func accountBalanceById(accStorage model.IBalanceInfoStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ids := mux.Vars(r)["id"]
		id, err := strconv.Atoi(ids)
		if err != nil {
			makeErrResponce(fmt.Sprintf(badRequestMessage+": Поле id должно быть числовым целочисленным типом больше 0."), http.StatusBadRequest, w)
			return
		}

		acc, custErr := accStorage.GetAccountBalance(id)
		if custErr != nil {
			makeErrResponce(internalErrorMessage, http.StatusInternalServerError, w)
			log.Print(custErr.Err.Error())
			return
		}

		currency := strings.ToUpper(r.FormValue("currency"))
		respMessage := accountByIdResponse{Id: acc.AccountId, Balance: acc.Balance, Currency: defaultCurrency}
		if currency != "" {
			balanceInCurrency, err := convert.ConvertToCurrency(acc.Balance, currency, &convert.ConvertDataStorerStruct{})
			if err == nil {
				respMessage.Balance = balanceInCurrency
				respMessage.Currency = currency
			} else {
				makeErrResponce("Не удалось предоставить информацию для выбранного курса валюты", http.StatusInternalServerError, w)
				log.Print(err.Error())
				return
			}
		}
		resp, _ := json.Marshal(respMessage)
		w.Header().Set("content-type", "application/json")
		w.Write(resp)
	}
}

//changeAccBalance - выполняет пополнение аккаунта на delta
//Пример тела запроса: {"Id":1,"Delta":-200}
func changeAccountBalance(accStorage model.IBalanceInfoStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		changeRequest := &changeAccBalanceRequest{}
		err := json.NewDecoder(r.Body).Decode(changeRequest)
		if err != nil {
			makeErrResponce(badRequestMessage, http.StatusBadRequest, w)
			return
		}
		if changeRequest.Delta == 0 {
			makeErrResponce(nullSumMessage, http.StatusBadRequest, w)
			return
		}

		successMessage, custErr := accStorage.ChangeAccountBalance(changeRequest.Id, changeRequest.Delta)

		if custErr != nil {
			if custErr.ErrCode == model.InsufficientFundsCode {
				makeErrResponce(insufficientFundsMessage, http.StatusForbidden, w)
			} else {
				makeErrResponce(internalErrorMessage, http.StatusInternalServerError, w)
			}
			log.Print(custErr.Err.Error())
			return
		}

		respMessage := changeAccBalanceResponse{Message: successMessage}
		resp, _ := json.Marshal(respMessage)
		w.Header().Set("content-type", "application/json")
		w.Write(resp)
	}
}

//transferSum - выполняет перевод суммы между аккаунтами
//пример тела запроса: {"Id1":1,"Id2":3,"Delta":-20}
func transferSum(accStorage model.IBalanceInfoStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		transferRequest := &transferSumRequest{}
		err := json.NewDecoder(r.Body).Decode(transferRequest)
		if err != nil {
			makeErrResponce(badRequestMessage, http.StatusBadRequest, w)
			return
		}
		if transferRequest.Id1 == transferRequest.Id2 {
			makeErrResponce(badRequestMessage, http.StatusBadRequest, w)
			return
		}
		if transferRequest.Delta == 0 {
			makeErrResponce(nullSumMessage, http.StatusBadRequest, w)
			return
		}

		transactionMessage, custErr := accStorage.TransferSumBetweenAccounts(transferRequest.Id1, transferRequest.Id2, transferRequest.Delta)

		if custErr != nil {
			if custErr.ErrCode == model.InsufficientFundsCode {
				makeErrResponce(insufficientFundsMessage, http.StatusForbidden, w)
			} else {
				makeErrResponce(internalErrorMessage, http.StatusInternalServerError, w)
			}
			log.Print(custErr.Err.Error())
			return
		}
		respMessage := transferSumResponce{Message: transactionMessage}
		resp, _ := json.Marshal(respMessage)
		w.Header().Set("content-type", "application/json")
		w.Write(resp)
	}
}

//transactionsHistory - выводит историю операций по аккаунту
//пример тела запроса {"Id":3,"SortedBy":"transaction_sum","SortedByDesc":true}
func transactionsHistory(accStorage model.IBalanceInfoStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		operationsInfoRequest := &transactionsHistoryRequest{}
		err := json.NewDecoder(r.Body).Decode(operationsInfoRequest)
		if err != nil {
			makeErrResponce(badRequestMessage, http.StatusBadRequest, w)
			return
		}

		history, custErr := accStorage.GetSortedTransactionsHistory(operationsInfoRequest.Id, operationsInfoRequest.SortedBy, operationsInfoRequest.SortedByDesc)
		if custErr != nil {
			if custErr.ErrCode == model.WrongInputParamsCode {
				makeErrResponce(badRequestMessage, http.StatusBadRequest, w)
			} else {
				makeErrResponce(internalErrorMessage, http.StatusInternalServerError, w)
			}
			log.Print(custErr.Err.Error())
			return
		}

		if len(history) == 0 {
			makeErrResponce("Отсутсвуют записи по выбранным условиям поиска", http.StatusNotFound, w)
			return
		}
		resp, _ := json.Marshal(history)
		w.Header().Set("content-type", "application/json")
		w.Write(resp)
	}
}
