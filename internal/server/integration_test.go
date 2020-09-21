// +build integration
package server

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/call-me-snake/user_balance_service/internal/model"
	"github.com/call-me-snake/user_balance_service/internal/storage"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"golang.org/x/exp/errors/fmt"
)

var (
	envs                  = []string{"POSTGRES_PASSWORD=example", "POSTGRES_DB=accounts", "POSTGRES_USER=postgres"}
	imageName             = "postgr_balance_storage_img"
	testContainerName     = "test_balance_storage"
	hostPort              = "5433"
	hostIp                = "0.0.0.0"
	testStorageConnString = "user=postgres password=example dbname=accounts sslmode=disable port=5433 host=localhost"
)

type balanceIntegrationTestSuite struct {
	suite.Suite
	Cli  *client.Client
	Resp container.ContainerCreateCreatedBody
	Db   model.IBalanceInfoStorage
}

//TestRun runs all tests
func TestRun(t *testing.T) {
	suite.Run(t, new(balanceIntegrationTestSuite))
}

//SetupSuite implements interface SetupAllSuite. method, which will run before the tests in the suite are run.
func (mySuite *balanceIntegrationTestSuite) SetupSuite() {
	ctx := context.Background()
	hostConfig := &container.HostConfig{
		AutoRemove:   true,
		PortBindings: nat.PortMap{"5432/tcp": []nat.PortBinding{{HostIP: hostIp, HostPort: hostPort}}},
	}
	var err error
	mySuite.Cli, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	mySuite.Resp, err = mySuite.Cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
		Tty:   false,
		Env:   envs,
	},
		hostConfig, nil, testContainerName)
	if err != nil {
		log.Fatal(err)
	}
	if err = mySuite.Cli.ContainerStart(ctx, mySuite.Resp.ID, types.ContainerStartOptions{}); err != nil {
		log.Fatal(err)
	}

	mySuite.Db, err = waitDbConnection(testStorageConnString, 10*time.Second)
	if err != nil {
		log.Fatal(err)
	}

}

//TearDownSuite implements interface TearDownAllSuite. method, which will run after all the tests in the suite have been run.
func (mySuite *balanceIntegrationTestSuite) TearDownSuite() {
	err := mySuite.Cli.ContainerStop(context.Background(), mySuite.Resp.ID, nil)
	if err != nil {
		log.Fatal(err)
	}
}

//TestAccountBalanceByIdSuccessful - тест успешного вызова ручки accountBalanceById
func (mySuite *balanceIntegrationTestSuite) TestAccountBalanceByIdSuccessful() {
	if mySuite.Db != nil {
		var accId = 1
		var balance = 500.0
		_, custErr := mySuite.Db.ChangeAccountBalance(accId, balance)
		assert.Nil(mySuite.T(), custErr)

		router := mux.NewRouter()
		router.HandleFunc("/account/balance/info/{id:[0-9]+}", accountBalanceById(mySuite.Db)).Methods("GET")
		req, err := http.NewRequest("GET", fmt.Sprintf("/account/balance/info/%d", accId), nil)
		if err != nil {
			log.Fatal(err)
		}
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		assert.Equal(mySuite.T(), http.StatusOK, rr.Code)

		testRespMessage := accountByIdResponse{Id: accId, Balance: balance, Currency: defaultCurrency}
		res, _ := json.Marshal(testRespMessage)
		assert.Equal(mySuite.T(), res, rr.Body.Bytes())
	}
}

//TestChangeAccountBalanceFail - тест ошибки изменения баланса из-за недостатка средств
func (mySuite *balanceIntegrationTestSuite) TestChangeAccountBalanceFail() {
	if mySuite.Db != nil {
		var accId = 2
		var delta = -1.0

		requestBody, _ := json.Marshal(changeAccBalanceRequest{Id: accId, Delta: delta})
		req, err := http.NewRequest("POST", "/account/balance/change", bytes.NewReader(requestBody))
		if err != nil {
			log.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(changeAccountBalance(mySuite.Db))
		handler.ServeHTTP(rr, req)
		assert.Equal(mySuite.T(), http.StatusForbidden, rr.Code)
	}
}

//TestChangeAccountBalanceSuccess - тест успешного пополнения баланса
func (mySuite *balanceIntegrationTestSuite) TestChangeAccountBalanceSuccess() {
	if mySuite.Db != nil {
		var accId = 3
		var delta = 1.0

		requestBody, _ := json.Marshal(changeAccBalanceRequest{Id: accId, Delta: delta})
		req, err := http.NewRequest("POST", "/account/balance/change", bytes.NewReader(requestBody))
		if err != nil {
			log.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(changeAccountBalance(mySuite.Db))
		handler.ServeHTTP(rr, req)
		assert.Equal(mySuite.T(), http.StatusOK, rr.Code)
	}
}

//TestTransferSumSuccess - тест успешной передачи суммы
func (mySuite *balanceIntegrationTestSuite) TestTransferSumSuccess() {
	if mySuite.Db != nil {
		var accId1 = 4
		var accId2 = 5
		var delta = 500.0

		_, custErr := mySuite.Db.ChangeAccountBalance(accId1, delta)
		assert.Nil(mySuite.T(), custErr)

		requestBody, _ := json.Marshal(transferSumRequest{Id1: accId1, Id2: accId2, Delta: delta})
		req, err := http.NewRequest("POST", "/account/balance/transfer", bytes.NewReader(requestBody))
		if err != nil {
			log.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(transferSum(mySuite.Db))
		handler.ServeHTTP(rr, req)
		assert.Equal(mySuite.T(), http.StatusOK, rr.Code)
	}
}

//TestTransferSumFail - тест ошибки передачи суммы из-за недостатка средств
func (mySuite *balanceIntegrationTestSuite) TestTransferSumFail() {
	if mySuite.Db != nil {
		var accId1 = 6
		var accId2 = 7
		var delta = 500.0

		requestBody, _ := json.Marshal(transferSumRequest{Id1: accId1, Id2: accId2, Delta: delta})
		req, err := http.NewRequest("POST", "/account/balance/transfer", bytes.NewReader(requestBody))
		if err != nil {
			log.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(transferSum(mySuite.Db))
		handler.ServeHTTP(rr, req)
		assert.Equal(mySuite.T(), http.StatusForbidden, rr.Code)
	}
}

//TestTransactionsHistorySuccess - тест успешного вывода истории транзакций
func (mySuite *balanceIntegrationTestSuite) TestTransactionsHistorySuccess() {
	if mySuite.Db != nil {
		var accId = 8
		var delta = 500.0

		_, custErr := mySuite.Db.ChangeAccountBalance(accId, delta)
		assert.Nil(mySuite.T(), custErr)

		requestBody, _ := json.Marshal(transactionsHistoryRequest{Id: accId})
		req, err := http.NewRequest("POST", "/account/balance/history", bytes.NewReader(requestBody))
		if err != nil {
			log.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(transactionsHistory(mySuite.Db))
		handler.ServeHTTP(rr, req)
		assert.Equal(mySuite.T(), http.StatusOK, rr.Code)
	}
}

//TestTransactionsHistoryFailNotFound - тест пустого вывода истории транзакций
func (mySuite *balanceIntegrationTestSuite) TestTransactionsHistoryFailNotFound() {
	if mySuite.Db != nil {
		var accId = 9

		requestBody, _ := json.Marshal(transactionsHistoryRequest{Id: accId})
		req, err := http.NewRequest("POST", "/account/balance/history", bytes.NewReader(requestBody))
		if err != nil {
			log.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(transactionsHistory(mySuite.Db))
		handler.ServeHTTP(rr, req)
		assert.Equal(mySuite.T(), http.StatusNotFound, rr.Code)
	}
}

//TestTransactionsHistoryFailWrongInput - тест ответа ручки TransactionsHistory при неверных входных параметрах
func (mySuite *balanceIntegrationTestSuite) TestTransactionsHistoryFailWrongInput() {
	if mySuite.Db != nil {
		var accId = 10

		requestBody, _ := json.Marshal(transactionsHistoryRequest{Id: accId, SortedBy: "wrong!"})
		req, err := http.NewRequest("POST", "/account/balance/history", bytes.NewReader(requestBody))
		if err != nil {
			log.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(transactionsHistory(mySuite.Db))
		handler.ServeHTTP(rr, req)
		assert.Equal(mySuite.T(), http.StatusBadRequest, rr.Code)
	}
}

func waitDbConnection(connString string, maxWait time.Duration) (db model.IBalanceInfoStorage, err error) {
	done := time.Now().Add(maxWait)
	for time.Now().Before(done) {
		db, err = storage.New(testStorageConnString)
		if err == nil {
			return db, nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return nil, err
}
