package storage

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/call-me-snake/user_balance_service/internal/model"
)

//методы, реализующие интерфейс model.ITransactionLogger

//CreateNewLog - реализует метод интерфейса ITransactionLogger
func (db *Storage) CreateNewLog(log model.Log) error {
	log.CreatedAt = time.Now()
	query := db.database.Create(&log)
	if query.Error != nil {
		return fmt.Errorf("storage.CreateNewLog: %v", query.Error)
	}
	return nil
}

//GetUserLogsFiltered - реализует метод интерфейса ITransactionLogger
func (db *Storage) GetUserLogsFiltered(filter map[string]interface{}) (filteredLogs []model.UserLog, err error, wrongInput bool) {
	query := db.database
	var accId float64
	var operationCompleted bool
	var orderDate, orderDelta string
	aId, ok := filter["user_id"]
	if !ok {
		return nil, errors.New("storage.GetUserLogsFiltered: неверный формат входных данных: отсутствует ключ user_id"), true
	}
	//при отправлении числа в структуре json без явно выраженного типа, оно приходит в формате float64
	if accId, ok = aId.(float64); ok {
		query = query.Where("account_id = ?", int(accId))
	} else {
		fmt.Println(reflect.TypeOf(aId))
		return nil, errors.New("storage.GetUserLogsFiltered: неверный формат входных данных: ключ user_id содержит неверный тип значения"), true
	}

	if oCompleted, ok := filter["operation_completed"]; ok {
		if operationCompleted, ok = oCompleted.(bool); ok {
			query = query.Where("operation_completed = ?", operationCompleted)
		} else {
			return nil, errors.New("storage.GetUserLogsFiltered: неверный формат входных данных: ключ operation_completed содержит неверный тип значения"), true
		}
	}

	if oDate, ok := filter["order_date"]; ok {
		if orderDate, ok = oDate.(string); ok {
			switch orderDate {
			case "desc":
				query = query.Order("created_at desc")
			case "asc":
				query = query.Order("created_at")
			default:
				return nil, errors.New("storage.GetUserLogsFiltered: неверный формат входных данных: ключ order_date должен содержать только значения desc или asc"), true
			}
		}
	}

	if oSum, ok := filter["order_sum"]; ok {
		if orderDelta, ok = oSum.(string); ok {
			switch orderDelta {
			case "desc":
				query = query.Order("delta desc")
			case "asc":
				query = query.Order("delta asc")
			default:
				return nil, errors.New("storage.GetUserLogsFiltered: неверный формат входных данных: ключ order_sum должен содержать только значения desc или asc"), true
			}
		}
	}
	query = query.Find(&filteredLogs)
	if query.Error != nil {
		return nil, fmt.Errorf("storage.GetUserLogsFiltered: %v", query.Error), false
	}
	if query.RowsAffected == 0 {
		return nil, nil, false
	}
	return filteredLogs, nil, false
}
