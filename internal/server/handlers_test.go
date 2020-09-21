package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/call-me-snake/user_balance_service/internal/model"
	mock_model "github.com/call-me-snake/user_balance_service/internal/model/mock"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/errors/fmt"
)

var (
	testId1                         = 1
	testId2                         = 2
	testBalance1                    = 1000.0
	testBalance2                    = 910.0
	testDelta1                      = 15.0
	testDelta2                      = 90.0
	testMessage                     = "Сообщение"
	testBalanceInfo1                = model.BalanceInfo{AccountId: testId1, Balance: testBalance1}
	testRespMessage1                = accountByIdResponse{Id: testId1, Balance: testBalance1, Currency: defaultCurrency}
	testErr1                        = model.CustomErr{Err: errors.New("Ошибка"), ErrCode: model.DefaultErrCode}
	testErrRespMessage1             = errorResponce{Message: internalErrorMessage, ErrCode: http.StatusInternalServerError}
	testChangeAccountBalanceRequest = changeAccBalanceRequest{Id: testId1, Delta: testDelta1}
	testTransferSumRequest          = transferSumRequest{Id1: testId1, Id2: testId2, Delta: testDelta1}
	testTransactionsHistoryRequest  = transactionsHistoryRequest{Id: testId1}
	testHistory                     = []model.TransactionRecord{
		{
			AccountId:          testId1,
			Delta:              testDelta1,
			RemainingBalance:   testBalance1,
			TransactionMessage: testMessage,
		},
		{
			AccountId:          testId2,
			Delta:              testDelta2,
			RemainingBalance:   testBalance2,
			TransactionMessage: testMessage,
		},
	}
)

//TestAliveHandler - тест aliveHandler
func TestAliveHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/ready", nil)
	if err != nil {
		log.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(aliveHandler)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, []byte("Hello from balance service"), rr.Body.Bytes())
}

//TestAccountBalanceById - тест положительного ответа ручки accountBalanceById
func TestAccountBalanceById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockdb := mock_model.NewMockIBalanceInfoStorage(ctrl)
	mockdb.EXPECT().GetAccountBalance(testId1).Return(&testBalanceInfo1, nil)

	//делаю с помощью mux.NewRouter() из-за mux.Vars
	router := mux.NewRouter()
	router.HandleFunc("/account/balance/info/{id:[0-9]+}", accountBalanceById(mockdb)).Methods("GET")
	req, err := http.NewRequest("GET", fmt.Sprintf("/account/balance/info/%d", testId1), nil)
	if err != nil {
		log.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	res, _ := json.Marshal(testRespMessage1)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, res, rr.Body.Bytes())
}

//TestAccountBalanceByIdErrResponce - тест возврата ручкой ошибки из-за внутренней ошибки бд
func TestAccountBalanceByIdErrResponce(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockdb := mock_model.NewMockIBalanceInfoStorage(ctrl)
	mockdb.EXPECT().GetAccountBalance(testId1).Return(nil, &testErr1)

	//делаю с помощью mux.NewRouter() из-за mux.Vars
	router := mux.NewRouter()
	router.HandleFunc("/account/balance/info/{id:[0-9]+}", accountBalanceById(mockdb)).Methods("GET")
	req, err := http.NewRequest("GET", fmt.Sprintf("/account/balance/info/%d", testId1), nil)
	if err != nil {
		log.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	res, _ := json.Marshal(testErrRespMessage1)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, res, rr.Body.Bytes())
}

//TestChangeAccountBalance - тест успешной смены баланса
func TestChangeAccountBalance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	message := fmt.Sprintf("Аккаунт %d успешно пополнен на сумму %.2f руб.", testId1, testDelta1)
	mockdb := mock_model.NewMockIBalanceInfoStorage(ctrl)
	mockdb.EXPECT().ChangeAccountBalance(testId1, testDelta1).Return(message, nil)

	requestBody, _ := json.Marshal(testChangeAccountBalanceRequest)
	res, _ := json.Marshal(changeAccBalanceResponse{Message: message})
	req, err := http.NewRequest("POST", "/account/balance/change", bytes.NewReader(requestBody))
	if err != nil {
		log.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(changeAccountBalance(mockdb))
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, res, rr.Body.Bytes())
}

//TestTransferSum - тест успешной передачи суммы
func TestTransferSum(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	message := fmt.Sprintf("Перевод на сумму %.2f руб. с аккаунта %d на аккаунт %d выполнен успешно.", testDelta1, testId1, testId2)
	mockdb := mock_model.NewMockIBalanceInfoStorage(ctrl)
	mockdb.EXPECT().TransferSumBetweenAccounts(testId1, testId2, testDelta1).Return(message, nil)

	requestBody, _ := json.Marshal(testTransferSumRequest)
	res, _ := json.Marshal(transferSumResponce{Message: message})
	req, err := http.NewRequest("POST", "/account/balance/transfer", bytes.NewReader(requestBody))
	if err != nil {
		log.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(transferSum(mockdb))
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, res, rr.Body.Bytes())
}

//TestTransactionsHistory - тест успешного вывода истории транзакций
func TestTransactionsHistory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockdb := mock_model.NewMockIBalanceInfoStorage(ctrl)
	mockdb.EXPECT().GetSortedTransactionsHistory(testId1, "", false).Return(testHistory, nil)

	requestBody, _ := json.Marshal(testTransactionsHistoryRequest)
	res, _ := json.Marshal(testHistory)

	req, err := http.NewRequest("POST", "/account/balance/history", bytes.NewReader(requestBody))
	if err != nil {
		log.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(transactionsHistory(mockdb))
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, res, rr.Body.Bytes())

}
