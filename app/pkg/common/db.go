package common

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	DB  *gorm.DB
	Err error
)

func init() {
	fp := Data + "database.db"
	param := "?_journal_mode=WAL"

	DB, Err = gorm.Open(sqlite.Open(fp+param), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{NoLowerCase: true},
	})
	if Err != nil {
		fmt.Println("error opening database: ", Err.Error())
	}
}
