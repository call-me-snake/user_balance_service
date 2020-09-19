package storage

import (
	"errors"
	"fmt"

	"github.com/call-me-snake/user_balance_service/internal/model"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

//методы, реализующие интерфейс model.IBalanceInfoStorage

//GetAccountBalance - реализует метод интерфейса IBalanceInfoStorage
func (db *Storage) GetAccountBalance(id int) (*model.BalanceInfo, *model.CustomErr) {
	result := &model.BalanceInfo{}
	query := db.database.First(result, id)
	if query.Error != nil {
		if query.Error == gorm.ErrRecordNotFound {
			err := &model.CustomErr{
				Err:     fmt.Errorf("storage.GetAccount: %v", query.Error),
				ErrCode: model.AccountNotExistsCode,
			}
			return nil, err
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
func (db *Storage) ChangeAccountBalance(id int, delta float64) (isChanged bool, err *model.CustomErr) {
	transaction := db.database.Begin()
	acc := &model.BalanceInfo{}
	query := transaction.First(acc, id)
	if query.Error != nil {
		transaction.Rollback()
		if query.Error == gorm.ErrRecordNotFound {
			err = &model.CustomErr{
				Err:     fmt.Errorf("storage.ChangeAccountBalance: %v", query.Error),
				ErrCode: model.AccountNotExistsCode,
			}
		} else {
			err = &model.CustomErr{
				Err:     fmt.Errorf("storage.ChangeAccountBalance: %v", query.Error),
				ErrCode: model.DefaultErrCode,
			}
		}
		return false, err
	}
	newAccountBalance := acc.Balance + delta
	if newAccountBalance < 0 {
		transaction.Rollback()
		err = &model.CustomErr{
			Err:     errors.New("storage.ChangeAccountBalance: Недостаточно средств на балансе"),
			ErrCode: model.InsufficientFundsCode,
		}
		return false, err
	}
	query = transaction.Model(&model.BalanceInfo{}).Where("account_id = ?", id).Update("balance", newAccountBalance)
	if query.Error != nil {
		transaction.Rollback()
		err = &model.CustomErr{
			Err:     fmt.Errorf("storage.ChangeAccountBalance: %v", query.Error),
			ErrCode: model.DefaultErrCode,
		}
		return false, err
	}
	transaction.Commit()

	return true, nil
}

//TransferSumBetweenAccounts - реализует метод интерфейса IBalanceInfoStorage
func (db *Storage) TransferSumBetweenAccounts(id1, id2 int, delta float64) (isTransferred bool, err *model.CustomErr) {
	transaction := db.database.Begin()
	acc1, acc2 := &model.BalanceInfo{}, &model.BalanceInfo{}
	query := transaction.First(acc1, id1)
	if query.Error != nil {
		transaction.Rollback()
		if query.Error == gorm.ErrRecordNotFound {
			err = &model.CustomErr{
				Err:     fmt.Errorf("storage.TransferSumBetweenAccounts: %v", query.Error),
				ErrCode: model.AccountNotExistsCode,
			}
		} else {
			err = &model.CustomErr{
				Err:     fmt.Errorf("storage.TransferSumBetweenAccounts: %v", query.Error),
				ErrCode: model.DefaultErrCode,
			}
		}
		return false, err
	}
	query = transaction.First(acc2, id2)
	if query.Error != nil {
		transaction.Rollback()
		if query.Error == gorm.ErrRecordNotFound {
			err = &model.CustomErr{
				Err:     fmt.Errorf("storage.TransferSumBetweenAccounts: %v", query.Error),
				ErrCode: model.AccountNotExistsCode,
			}
		} else {
			err = &model.CustomErr{
				Err:     fmt.Errorf("storage.TransferSumBetweenAccounts: %v", query.Error),
				ErrCode: model.DefaultErrCode,
			}
		}
		return false, err
	}

	newAccountBalance1 := acc1.Balance - delta
	newAccountBalance2 := acc2.Balance + delta
	if newAccountBalance1 < 0 || newAccountBalance2 < 0 {
		transaction.Rollback()
		err = &model.CustomErr{
			Err:     errors.New("storage.TransferSumBetweenAccounts: Недостаточно средств на балансе"),
			ErrCode: model.DefaultErrCode,
		}
		return false, err
	}
	query = transaction.Model(&model.BalanceInfo{}).Where("account_id = ?", id1).Update("balance", newAccountBalance1)
	if query.Error != nil {
		transaction.Rollback()
		err = &model.CustomErr{
			Err:     fmt.Errorf("storage.TransferSumBetweenAccounts: %v", query.Error),
			ErrCode: model.DefaultErrCode,
		}
		return false, err
	}
	query = transaction.Model(&model.BalanceInfo{}).Where("account_id = ?", id2).Update("balance", newAccountBalance2)
	if query.Error != nil {
		transaction.Rollback()
		err = &model.CustomErr{
			Err:     fmt.Errorf("storage.TransferSumBetweenAccounts: %v", query.Error),
			ErrCode: model.DefaultErrCode,
		}
		return false, err
	}
	transaction.Commit()
	return true, nil
}
