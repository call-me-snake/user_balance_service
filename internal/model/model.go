package model

import "time"

const (
	DefaultErrCode        = 0
	InsufficientFundsCode = 1
	WrongInputParamsCode  = 2

	//Строковые константы используются в качестве возможных значений поля sortedBy в методе IBalanceInfoStorage.GetSortedTransactionsHistory
	TransactionSum  = "transaction_sum"
	TransactionTime = "transaction_time"
)

//IBalanceInfoStorage - интерфейс для работы с балансом пользователей
type IBalanceInfoStorage interface {
	//GetAccountBalance - получение баланса аккаунта
	GetAccountBalance(id int) (*BalanceInfo, *CustomErr)
	//ChangeAccountBalance: баланс меняется по принципу newBalance = curBalance + delta
	ChangeAccountBalance(id int, delta float64) (successMessage string, err *CustomErr)
	//TransferSumBetweenAccounts: delta может быть как положительной, так и отрицательной
	//баланс аккаунтов меняется по принципу newBalance1 = curBalance1 - delta; newBalance2 = curBalance2 + delta
	TransferSumBetweenAccounts(id1, id2 int, delta float64) (successMessage string, err *CustomErr)
	//GetSortedTransactionsHistory - получение отсортированной истории переводов для пользователя
	GetSortedTransactionsHistory(id int, sortedBy string, sortedByDesc bool) (history []TransactionRecord, err *CustomErr)
}

//BalanceInfo - структура для хранения информации по балансу пользователя
type BalanceInfo struct {
	AccountId int     `gorm:"primary_key;column:account_id"`
	Balance   float64 `gorm:"column:balance"`
}

// TableName - declare table name for GORM
func (BalanceInfo) TableName() string {
	return "accounts"
}

//CustomErr - кастомный тип ошибки, возвращаемый методами интерфейса IBalanceInfoStorage.
//Содержит переменную ErrCode, указывающую на тип ошибки
type CustomErr struct {
	Err     error
	ErrCode int
}

//TransactionRecord - структура для сохранения успешного изменения баланса в истории
type TransactionRecord struct {
	AccountId          int       `gorm:"column:account_id"`
	Delta              float64   `gorm:"column:delta"`
	RemainingBalance   float64   `gorm:"column:remaining_balance"`
	TransactionMessage string    `gorm:"column:transaction_message"`
	CreatedAt          time.Time `gorm:"column:created_at"`
}

// TableName - declare table name for GORM
func (TransactionRecord) TableName() string {
	return "transactions_history"
}

//Config хранит переменные окружения
type Config struct {
	ServerAddress      string
	AccountStorageConn string
}

//ConvertData - структура для хранения коэффициэнтов конвертирования
type ConvertData struct {
	FillingTime time.Time
	Base        string             `json:"base"`
	Date        string             `json:"date"`
	Rates       map[string]float64 `json:"rates"`
}
