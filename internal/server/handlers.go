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
const accountNotFoundMessage = "Аккаунт не найден"
const internalErrorMessage = "Внутренняя ошибка сервера"
const insufficientFundsMessage = "Недостаточно средств на счету"
const defaultCurrency = "RUB"

func aliveHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from balance service"))
}

//accountById - возврат информации об аккаунте
func accountById(accStorage model.IAccountsStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ids := mux.Vars(r)["id"]
		id, err := strconv.Atoi(ids)
		if err != nil {
			makeErrResponce(badRequestMessage, http.StatusBadRequest, w)
			return
		}

		acc, custErr := accStorage.GetAccount(id)
		if custErr != nil {
			if custErr.AccountNotExists {
				makeErrResponce(accountNotFoundMessage, http.StatusNotFound, w)
			} else {
				makeErrResponce(internalErrorMessage, http.StatusInternalServerError, w)
			}
			log.Print(custErr.Err.Error())
			return
		}

		currency := strings.ToUpper(r.FormValue("currency"))
		respMessage := accountByIdResponse{Id: acc.AccountId, Balance: acc.Balance, Currency: defaultCurrency}
		if currency != "" {
			balanceInCurrency, err := convert.ConvertToCurrency(acc.Balance, currency, convert.UpdaterStruct{})
			if err == nil {
				respMessage.Balance = balanceInCurrency
				respMessage.Currency = currency
			} else {
				log.Print(err.Error())
			}
		}
		resp, _ := json.Marshal(respMessage)
		w.Header().Set("content-type", "application/json")
		w.Write(resp)
	}
}

//changeAccBalance - выполняет пополнение аккаунта на delta
//Пример тела запроса: {"Id":1,"Delta":-200}
func changeAccBalance(accStorage model.IAccountsStorage, logger model.ITransactionLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		changeRequest := &changeAccBalanceRequest{}
		err := json.NewDecoder(r.Body).Decode(changeRequest)
		if err != nil {
			makeErrResponce(badRequestMessage, http.StatusBadRequest, w)
			return
		}

		isChanged, custErr := accStorage.ChangeAccountBalance(changeRequest.Id, changeRequest.Delta)
		var successMessage string

		//Логирование
		defer func() {
			operationLog := model.Log{
				AccountId: changeRequest.Id,
				Delta:     changeRequest.Delta,
			}
			if custErr != nil {
				if custErr.AccountNotExists {
					return
				} else {
					operationLog.LogMessage = custErr.Err.Error()
					operationLog.OperationCompleted = false
				}
			} else {
				operationLog.LogMessage = successMessage
				operationLog.OperationCompleted = true
			}
			err = logger.CreateNewLog(operationLog)
			if err != nil {
				log.Printf("changeAccBalance: не удалось логировать %v. Ошибка %v", operationLog, err)
				return
			}
		}()

		if !isChanged && custErr != nil {
			if custErr.AccountNotExists {
				makeErrResponce(accountNotFoundMessage, http.StatusNotFound, w)
			} else if custErr.InsufficientFunds {
				makeErrResponce(insufficientFundsMessage, http.StatusForbidden, w)
			} else {
				makeErrResponce(internalErrorMessage, http.StatusInternalServerError, w)
			}
			log.Print(custErr.Err.Error())
			return
		}

		if changeRequest.Delta >= 0 {
			successMessage = fmt.Sprintf("Аккаунт %d успешно пополнен на сумму %.2f руб.", changeRequest.Id, changeRequest.Delta)
		} else {
			successMessage = fmt.Sprintf("С аккаунта %d успешно снята сумма %.2f руб.", changeRequest.Id, -changeRequest.Delta)
		}
		respMessage := changeAccBalanceResponse{Message: successMessage}
		resp, _ := json.Marshal(respMessage)
		w.Header().Set("content-type", "application/json")
		w.Write(resp)
	}
}

//transferSum - выполняет перевод суммы между аккаунтами
//пример тела запроса: {"Id1":1,"Id2":3,"Delta":-20}
func transferSum(accStorage model.IAccountsStorage, logger model.ITransactionLogger) http.HandlerFunc {
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

		isTransfered, custErr := accStorage.TransferSumBetweenAccounts(transferRequest.Id1, transferRequest.Id2, transferRequest.Delta)
		var successMessage string

		//логирование
		defer func() {
			operationLog1 := model.Log{
				AccountId: transferRequest.Id1,
				Delta:     -transferRequest.Delta,
			}
			operationLog2 := model.Log{
				AccountId: transferRequest.Id2,
				Delta:     transferRequest.Delta,
			}
			if custErr != nil {
				if custErr.AccountNotExists {
					return
				} else {
					operationLog1.LogMessage, operationLog2.LogMessage = custErr.Err.Error(), custErr.Err.Error()
					operationLog1.OperationCompleted, operationLog2.OperationCompleted = false, false
				}
			} else {
				operationLog1.LogMessage, operationLog2.LogMessage = successMessage, successMessage
				operationLog1.OperationCompleted, operationLog2.OperationCompleted = true, true
			}
			err = logger.CreateNewLog(operationLog1)
			if err != nil {
				log.Printf("transferSum: не удалось логировать %v. Ошибка %v", operationLog1, err)
				return
			}
			err = logger.CreateNewLog(operationLog2)
			if err != nil {
				log.Printf("transferSum: не удалось логировать %v. Ошибка %v", operationLog1, err)
				return
			}
		}()

		if !isTransfered && custErr != nil {
			if custErr.AccountNotExists {
				makeErrResponce(accountNotFoundMessage, http.StatusNotFound, w)
			} else if custErr.InsufficientFunds {
				makeErrResponce(insufficientFundsMessage, http.StatusForbidden, w)
			} else {
				makeErrResponce(internalErrorMessage, http.StatusInternalServerError, w)
			}
			log.Print(custErr.Err.Error())
			return
		}

		if transferRequest.Delta >= 0 {
			successMessage = fmt.Sprintf("Перевод на сумму %.2f с аккаунта %d на аккаунт %d выполнен успешно.", transferRequest.Delta, transferRequest.Id1, transferRequest.Id2)
		} else {
			successMessage = fmt.Sprintf("Перевод на сумму %.2f с аккаунта %d на аккаунт %d выполнен успешно.", -transferRequest.Delta, transferRequest.Id2, transferRequest.Id1)
		}
		respMessage := transferSumResponce{Message: successMessage}
		resp, _ := json.Marshal(respMessage)
		w.Header().Set("content-type", "application/json")
		w.Write(resp)
	}
}
