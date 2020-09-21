package convert

import (
	"errors"
	"testing"

	mock_convert "github.com/call-me-snake/user_balance_service/internal/convert/mock"
	"github.com/call-me-snake/user_balance_service/internal/model"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var (
	testConvertData1 = model.ConvertData{
		Base: rubCurrency,
		Rates: map[string]float64{
			dollarCur: 0.025,
		},
	}
	dollarCur             = "USD"
	rubBalance            = 40.0
	expectedDollarBalance = 1.0
)

//TestConvertToCurrency - тест успешной конвертации
func TestConvertToCurrencySuccessful(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockStorer := mock_convert.NewMockConvertDataStorer(ctrl)
	mockStorer.EXPECT().GetConvertData().Return(testConvertData1, nil)
	dollarBalance, err := ConvertToCurrency(rubBalance, dollarCur, mockStorer)
	assert.Nil(t, err)
	assert.Equal(t, expectedDollarBalance, dollarBalance)
}

//TestConvertToCurrencyFail1 - тест ошибки из-за ненулевой ошибки на выходе GetConvertData
func TestConvertToCurrencyFail1(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockStorer := mock_convert.NewMockConvertDataStorer(ctrl)
	mockStorer.EXPECT().GetConvertData().Return(model.ConvertData{}, errors.New("Ошибка"))
	_, err := ConvertToCurrency(1, "ничего не значащая строка", mockStorer)
	assert.Error(t, err)
}

//TestConvertToCurrencyFail2 - тест ошибки из-за поля Base в model.ConvertData
func TestConvertToCurrencyFail2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockStorer := mock_convert.NewMockConvertDataStorer(ctrl)
	mockStorer.EXPECT().GetConvertData().Return(model.ConvertData{Base: "Не рубль"}, nil)
	_, err := ConvertToCurrency(1, "ничего не значащая строка", mockStorer)
	assert.Error(t, err)
}

//TestConvertToCurrencyEmptyCur - тест ошибки из-за пустой переменной currency на входе
func TestConvertToCurrencyEmptyCur(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockStorer := mock_convert.NewMockConvertDataStorer(ctrl)
	_, err := ConvertToCurrency(1, "", mockStorer)
	assert.Error(t, err)
}
