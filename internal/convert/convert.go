package convert

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/call-me-snake/user_balance_service/internal/model"
)

var (
	convertDataStorage model.ConvertData
	coursesStorerUrl   = "https://api.exchangeratesapi.io/latest?base=RUB"
)

const rubCurrency = "RUB"
const updateDataInterval time.Duration = time.Hour

//ConvertDataStorer - содержит метод GetConvertData. Нужен для mock, чтобы не вызывать http
type ConvertDataStorer interface {
	GetConvertData() (model.ConvertData, error)
}

//ConvertDataStorerStruct - структура для реализации Updater
type ConvertDataStorerStruct struct{}

//GetConvertData - получает структуру данных, необходимую для конвертации валют
func (c *ConvertDataStorerStruct) GetConvertData() (model.ConvertData, error) {
	t := convertDataStorage.FillingTime
	if time.Since(t) > updateDataInterval {
		req, _ := http.NewRequest(http.MethodGet, coursesStorerUrl, nil)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return convertDataStorage, fmt.Errorf("convert.getConvertData: %v", err)
		}
		r, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return convertDataStorage, fmt.Errorf("convert.getConvertData: %v", err)
		}
		err = json.Unmarshal(r, &convertDataStorage)
		if err != nil {
			return convertDataStorage, fmt.Errorf("convert.getConvertData: %v", err)
		}
		convertDataStorage.FillingTime = time.Now()
	}
	return convertDataStorage, nil
}

//ConvertToCurrency - конвертирует сумму в выбранную валюту
func ConvertToCurrency(balance float64, currency string, storer ConvertDataStorer) (balanceInCurrency float64, err error) {
	if currency == "" {
		return 0, errors.New("convert.ConvertToCurrency: Пустая строка на входе")
	}
	var data model.ConvertData
	if data, err = storer.GetConvertData(); err != nil {
		return 0, fmt.Errorf("convert.ConvertToCurrency: %s", err.Error())
	}
	if data.Base != rubCurrency {
		return 0, fmt.Errorf("convert.ConvertToCurrency: convertDataStorage содержит неверную информацию: %#v", data)
	}
	if course, ok := data.Rates[currency]; ok {
		balanceInCurrency = course * balance
		return balanceInCurrency, nil
	}
	return 0, fmt.Errorf("convert.ConvertToCurrency: convertDataStorage %#v не содержит значения cur: %s", data, currency)
}
