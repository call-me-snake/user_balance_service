package storage

import (
	"fmt"
	"log"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

//sleepDurationInSec - время пинга функции checkConnection в секундах
const sleepDurationInSec = 5

//Storage ...
type Storage struct {
	database *gorm.DB
	address  string
}

//New возвращает объект интерфейса IAccountsStorage (storage)
func New(adress string) (*Storage, error) {
	var err error
	db := &Storage{}
	db.address = adress
	db.database, err = gorm.Open("postgres", adress)
	if err != nil {
		return nil, fmt.Errorf("storage.New: %v", err)
	}

	err = db.ping()
	if err != nil {
		return nil, fmt.Errorf("storage.New: %s", err.Error())
	}
	db.checkConnection()

	return db, nil
}

//ping (internal)
func (db *Storage) ping() error {
	//db.database.LogMode(true)
	result := struct {
		Result int
	}{}

	err := db.database.Raw("select 1+1 as result").Scan(&result).Error
	if err != nil {
		return fmt.Errorf("storage.ping: %v", err)
	}
	if result.Result != 2 {
		return fmt.Errorf("storage.ping: incorrect result!=2 (%d)", result.Result)
	}
	return nil
}

//checkConnection (internal)
func (db *Storage) checkConnection() {
	go func() {
		for {
			err := db.ping()
			if err != nil {

				log.Printf("storage.checkConnection: no connection: %s", err.Error())
				tempDb, err := gorm.Open("postgres", db.address)

				if err != nil {
					log.Printf("storage.checkConnection: could not establish connection: %v", err)
				} else {
					db.database = tempDb
				}
			}
			time.Sleep(sleepDurationInSec * time.Second)
		}
	}()
}
