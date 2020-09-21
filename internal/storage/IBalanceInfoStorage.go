package storage

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/call-me-snake/user_balance_service/internal/model"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const positiveBalanceConstraint = "positive_balance"

//методы, реализующие интерфейс model.IBalanceInfoStorage

//GetAccountBalance - реализует метод интерфейса IBalanceInfoStorage
func (db *storage) GetAccountBalance(id int) (*model.BalanceInfo, *model.CustomErr) {
	result := &model.BalanceInfo{}
	query := db.database.First(result, id)
	if query.Error != nil {
		if query.Error == gorm.ErrRecordNotFound {
			return &model.BalanceInfo{AccountId: id, Balance: 0}, nil
		}
		err := &model.CustomErr{
			Err:     fmt.Errorf("storage.GetAccount: %v", query.Error),
			ErrCode: model.DefaultErrCode,
		}
		return nil, err
	}
	return result, nil
}

//ChangeAccountBalance - реализует метод интерфейса IBalanceInfoStorage
func (db *storage) ChangeAccountBalance(id int, delta float64) (successMessage string, err *model.CustomErr) {
	//начало транзакции
	transaction := db.database.Begin()
	acc := &model.BalanceInfo{AccountId: id}
	//попытка изменения баланса
	err = updateOrCreateBalanceInfo(transaction, id, delta)
	if err != nil {
		transaction.Rollback()
		return "", err
	}

	//получение измененной суммы
	query := transaction.First(acc, id)
	if query.Error != nil {
		transaction.Rollback()
		err = &model.CustomErr{
			Err:     fmt.Errorf("storage.ChangeAccountBalance: %v", query.Error),
			ErrCode: model.DefaultErrCode,
		}
		return "", err
	}

	record := &model.TransactionRecord{
		AccountId:        id,
		Delta:            delta,
		RemainingBalance: acc.Balance,
		CreatedAt:        time.Now(),
	}

	if delta > 0 {
		record.TransactionMessage = fmt.Sprintf("Аккаунт %d успешно пополнен на сумму %.2f руб.", id, delta)
	} else {
		record.TransactionMessage = fmt.Sprintf("С аккаунта %d успешно снята сумма %.2f руб.", id, -delta)
	}
	//сохранение изменения баланса
	query = transaction.Create(record)

	if query.Error != nil {
		transaction.Rollback()
		err = &model.CustomErr{
			Err:     fmt.Errorf("storage.ChangeAccountBalance: %v", query.Error),
			ErrCode: model.DefaultErrCode,
		}
		return "", err
	}
	//конец транзакции
	transaction.Commit()

	return record.TransactionMessage, nil
}

//TransferSumBetweenAccounts - реализует метод интерфейса IBalanceInfoStorage
func (db *storage) TransferSumBetweenAccounts(id1, id2 int, delta float64) (successMessage string, err *model.CustomErr) {
	//начало транзакции
	transaction := db.database.Begin()
	acc1, acc2 := &model.BalanceInfo{AccountId: id1}, &model.BalanceInfo{AccountId: id2}
	//попытка передачи суммы
	err = updateOrCreateBalanceInfo(transaction, id1, -delta)
	if err != nil {
		transaction.Rollback()
		return "", err
	}

	err = updateOrCreateBalanceInfo(transaction, id2, delta)
	if err != nil {
		transaction.Rollback()
		return "", err
	}

	//получение изменений
	query := transaction.First(acc1, id1)
	if query.Error != nil {
		transaction.Rollback()
		err = &model.CustomErr{
			Err:     fmt.Errorf("storage.TransferSumBetweenAccounts: %v", query.Error),
			ErrCode: model.DefaultErrCode,
		}
		return "", err
	}
	query = transaction.First(acc2, id2)
	if query.Error != nil {
		transaction.Rollback()
		err = &model.CustomErr{
			Err:     fmt.Errorf("storage.TransferSumBetweenAccounts: %v", query.Error),
			ErrCode: model.DefaultErrCode,
		}
		return "", err
	}

	record1 := &model.TransactionRecord{
		AccountId:        id1,
		Delta:            -delta,
		RemainingBalance: acc1.Balance,
		CreatedAt:        time.Now(),
	}

	record2 := &model.TransactionRecord{
		AccountId:        id2,
		Delta:            delta,
		RemainingBalance: acc2.Balance,
		CreatedAt:        time.Now(),
	}

	var transactionMessage string
	if delta > 0 {
		transactionMessage = fmt.Sprintf("Перевод на сумму %.2f руб. с аккаунта %d на аккаунт %d выполнен успешно.", delta, id1, id2)
	} else {
		transactionMessage = fmt.Sprintf("Перевод на сумму %.2f руб. с аккаунта %d на аккаунт %d выполнен успешно.", -delta, id2, id1)
	}
	record1.TransactionMessage, record2.TransactionMessage = transactionMessage, transactionMessage

	//сохранение в истории
	query = transaction.Create(record1)
	if query.Error != nil {
		transaction.Rollback()
		err = &model.CustomErr{
			Err:     fmt.Errorf("storage.ChangeAccountBalance: %v", query.Error),
			ErrCode: model.DefaultErrCode,
		}
		return "", err
	}

	query = transaction.Create(record2)
	if query.Error != nil {
		transaction.Rollback()
		err = &model.CustomErr{
			Err:     fmt.Errorf("storage.ChangeAccountBalance: %v", query.Error),
			ErrCode: model.DefaultErrCode,
		}
		return "", err
	}
	transaction.Commit()
	return transactionMessage, nil
}

//GetSortedTransactionsHistory - реализует метод интерфейса IBalanceInfoStorage
func (db *storage) GetSortedTransactionsHistory(id int, sortedBy string, sortedByDesc bool) (history []model.TransactionRecord, err *model.CustomErr) {
	query := db.database.Where("account_id = ?", id)

	if sortedBy != "" {
		var sortBy string
		switch sortedBy {
		case model.TransactionSum:
			sortBy = "delta"
		case model.TransactionTime:
			sortBy = "created_at"
		default:
			err = &model.CustomErr{
				Err:     fmt.Errorf("storage.GetSortedTransactionsHistory: некорректный входной параметр sortedBy: %s . sortedBy должен быть равен пустой строке, либо строковой константе в internal/model", sortedBy),
				ErrCode: model.WrongInputParamsCode,
			}
			return nil, err
		}

		if sortedByDesc {
			sortBy += " desc"
		}

		query = query.Order(sortBy)
	}

	query = query.Find(&history)
	if query.Error != nil {
		err = &model.CustomErr{
			Err:     fmt.Errorf("storage.GetSortedTransactionsHistory: %v", query.Error),
			ErrCode: model.DefaultErrCode,
		}
		return nil, err
	}
	if query.RowsAffected == 0 {
		return nil, nil
	}
	return history, nil
}

func updateOrCreateBalanceInfo(transaction *gorm.DB, id int, delta float64) (err *model.CustomErr) {
	query := transaction.Model(model.BalanceInfo{AccountId: id}).UpdateColumn("balance", gorm.Expr("balance + ?", delta))
	if query.Error == nil && query.RowsAffected == 0 {
		if delta > 0 {
			//Попытка создания новой записи в случае отсутствия ее в таблице
			query = transaction.Create(&model.BalanceInfo{AccountId: id, Balance: delta})
			if query.Error != nil {
				err = &model.CustomErr{
					Err:     fmt.Errorf("storage.ChangeAccountBalance: %v", query.Error),
					ErrCode: model.DefaultErrCode,
				}
				return err
			}
		} else {
			err = &model.CustomErr{
				Err:     errors.New("storage.ChangeAccountBalance: Попытка отрицательного пополнения аккаунта с нулевым балансом"),
				ErrCode: model.InsufficientFundsCode,
			}
			return err
		}
	} else if query.Error != nil {
		if strings.Contains(query.Error.Error(), positiveBalanceConstraint) {
			err = &model.CustomErr{
				Err:     fmt.Errorf("storage.ChangeAccountBalance: %v", query.Error),
				ErrCode: model.InsufficientFundsCode,
			}
		} else {
			err = &model.CustomErr{
				Err:     fmt.Errorf("storage.ChangeAccountBalance: %v", query.Error),
				ErrCode: model.DefaultErrCode,
			}
		}
		return err
	}
	return nil
}
