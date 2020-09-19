package model

import "time"

const (
	DefaultErrCode        = 0
	InsufficientFundsCode = 1
	AccountNotExistsCode  = 2
)

//IBalanceInfoStorage - интерфейс для работы с балансом пользователей
type IBalanceInfoStorage interface {
	GetAccountBalance(id int) (*BalanceInfo, *CustomErr)
	//ChangeAccountBalance: баланс меняется по принципу newBalance = curBalance + delta
	ChangeAccountBalance(id int, delta float64) (isChanged bool, err *CustomErr)
	//TransferSumBetweenAccounts: delta может быть как положительной, так и отрицательной
	//баланс аккаунтов меняется по принципу newBalance1 = curBalance1 - delta; newBalance2 = curBalance2 + delta
	TransferSumBetweenAccounts(id1, id2 int, delta float64) (isTransferred bool, err *CustomErr)
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

//ITransactionLogger - интерфейс для логирования успешных транзакций
type ITransactionLogger interface {
	CreateNewLog(log Log) error
	GetUserLogsFiltered(filter map[string]interface{}) (logs []UserLog, err error, wrongInput bool)
}

//Log - структура для хранения логов по транзакциям
type Log struct {
	UserLog
	LogInternalMessage string `gorm:"column:log_internal_message"`
}

//UserLog - структура для хранения пользовательской информации по транзакциям
type UserLog struct {
	AccountId          int       `gorm:"column:account_id"`
	Delta              float64   `gorm:"column:delta"`
	LogUserMessage     string    `gorm:"column:log_user_message"`
	OperationCompleted bool      `gorm:"column:operation_completed"`
	CreatedAt          time.Time `gorm:"column:created_at"`
}

// TableName - declare table name for GORM
func (UserLog) TableName() string {
	return "logs"
}
