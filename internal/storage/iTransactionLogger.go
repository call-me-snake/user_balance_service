package storage

import (
	"fmt"
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
