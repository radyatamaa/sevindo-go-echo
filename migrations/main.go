package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"

	model "github.com/models"
)

func main() {
	//dev
	db, err := gorm.Open("mysql", "adminbkni@bkni-ri:Standar123.@(bkni-ri.mysql.database.azure.com)/sevindo_dev?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println(err)
	}

	migration := model.MigrationHistory{}
	errmigration := db.AutoMigrate(&migration)
	if errmigration != nil {
		migration := model.MigrationHistory{
			DescMigration: "Add Table Migration",
			Date:          time.Now(),
		}

		db.Create(&migration)
	}


	user := model.User{}
	erruser := db.AutoMigrate(&user)
	if erruser != nil {
		migration := model.MigrationHistory{
			DescMigration: "Add Table User",
			Date:          time.Now(),
		}

		db.Create(&migration)
	}

}
