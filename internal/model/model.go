package model

import "time"

//IAccountsStorage - интерфейс для работы с балансом пользователей
type IAccountsStorage interface {
	GetAccount(id int) (*Account, *CustomErr)
	//ChangeAccountBalance: баланс меняется по принципу newBalance = curBalance + delta
	ChangeAccountBalance(id int, delta float64) (isChanged bool, err *CustomErr)
	//TransferSumBetweenAccounts: delta может быть как положительной, так и отрицательной
	//баланс аккаунтов меняется по принципу newBalance1 = curBalance1 - delta; newBalance2 = curBalance2 + delta
	TransferSumBetweenAccounts(id1, id2 int, delta float64) (isTransferred bool, err *CustomErr)
}

//Account - структура для хранения информации по балансу пользователя
type Account struct {
	AccountId int     `gorm:"primary_key;column:account_id"`
	Balance   float64 `gorm:"column:balance"`
}

//CustomErr - кастомный тип ошибки, возвращаемый методами интерфейса IAccountsStorage.
//Содержит флаг InsufficientFunds, указывающий на недостаток средств на счету для проведения операции
//Содержит флаг AccountNotExists, поднимаемый при ошибках базы вида "record not found"
type CustomErr struct {
	Err               error
	InsufficientFunds bool
	AccountNotExists  bool
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
