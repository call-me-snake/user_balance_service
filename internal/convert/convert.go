package convert

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	convertDataStorage ConvertData
	coursesStorerUrl   = "https://api.exchangeratesapi.io/latest?base=RUB"
)

const rubCurrency = "RUB"

//VerifyOrUpdateConvertData - получает структуру данных, необходимую для конвертации валют
func (UpdaterStruct) VerifyOrUpdateConvertData() (ConvertData, error) {
	t := convertDataStorage.fillingTime
	if time.Since(t) > time.Hour {
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
		convertDataStorage.fillingTime = time.Now()
	}
	return convertDataStorage, nil
}

//ConvertData - структура для хранения коэффициэнтов конвертирования
type ConvertData struct {
	fillingTime time.Time
	Base        string             `json:"base"`
	Date        string             `json:"date"`
	Rates       map[string]float64 `json:"rates"`
}

//Updater - содержит метод VerifyOrUpdateConvertData. Нужен для mock, чтобы не вызывать http
type Updater interface {
	VerifyOrUpdateConvertData() (ConvertData, error)
}
type UpdaterStruct struct{}

//ConvertToCurrency - конвертирует сумму в выбранную валюту
func ConvertToCurrency(balance float64, currency string, updater Updater) (balanceInCurrency float64, err error) {
	if currency == "" {
		return 0, errors.New("convert.ConvertToCurrency: Пустая строка на входе")
	}
	var data ConvertData
	if data, err = updater.VerifyOrUpdateConvertData(); err != nil {
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
