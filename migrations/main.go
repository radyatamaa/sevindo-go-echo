package main

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

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

	country := model.Country{}
	errcountry := db.AutoMigrate(&country)
	if errcountry != nil {
		migration := model.MigrationHistory{
			DescMigration: "Add Table Country",
			Date:          time.Now(),
		}

		db.Create(&migration)
	}

	branch := model.Branch{}
	errbranch := db.AutoMigrate(&branch)
	if errbranch != nil {
		migration := model.MigrationHistory{
			DescMigration: "Add Table Branch",
			Date:          time.Now(),
		}

		db.Create(&migration)
	}

	userAdmin := model.UserAdmin{}
	erruserAdmin := db.AutoMigrate(&userAdmin).AddForeignKey("branch_id", "branches(id)", "RESTRICT", "RESTRICT")
	if erruserAdmin != nil {
		migration := model.MigrationHistory{
			DescMigration: "Add Table User Admin",
			Date:          time.Now(),
		}
		db.Create(&migration)
	}
	currency := model.Currency{}
	errcurrency := db.AutoMigrate(&currency)
	if errcurrency != nil {
		migration := model.MigrationHistory{
			DescMigration: "Add Table Currency",
			Date:          time.Now(),
		}

		db.Create(&migration)
	}

	language := model.Language{}
	errlanguage := db.AutoMigrate(&language)
	if errlanguage != nil {
		migration := model.MigrationHistory{
			DescMigration: "Add Table language",
			Date:          time.Now(),
		}

		db.Create(&migration)
	}
	province := model.Province{}
	errprovince := db.AutoMigrate(&province).AddForeignKey("country_id", "countries(id)", "RESTRICT", "RESTRICT")
	if errprovince != nil {
		migration := model.MigrationHistory{
			DescMigration: "Add Table Province",
			Date:          time.Now(),
		}

		db.Create(&migration)
	}

	city := model.City{}
	errcity := db.AutoMigrate(&city).AddForeignKey("province_id", "provinces(id)", "RESTRICT", "RESTRICT")
	if errcity != nil {
		migration := model.MigrationHistory{
			DescMigration: "Add Table City",
			Date:          time.Now(),
		}

		db.Create(&migration)
	}
	districts := model.Districts{}
	errdistricts := db.AutoMigrate(&districts).AddForeignKey("city_id", "cities(id)", "RESTRICT", "RESTRICT")
	if errdistricts != nil {
		migration := model.MigrationHistory{
			DescMigration: "Add Table Districts",
			Date:          time.Now(),
		}

		db.Create(&migration)
	}
	accessibility := model.Accessibility{}
	erraccessibility := db.AutoMigrate(&accessibility)
	if erraccessibility != nil {
		migration := model.MigrationHistory{
			DescMigration: "Add Table Accessibility",
			Date:          time.Now(),
		}

		db.Create(&migration)
	}

	amenities := model.Amenities{}
	erramenities := db.AutoMigrate(&amenities)
	if erramenities != nil {
		migration := model.MigrationHistory{
			DescMigration: "Add Table Amenities",
			Date:          time.Now(),
		}

		db.Create(&migration)
	}

	resort := model.Resort{}
	errResort := db.AutoMigrate(&resort).AddForeignKey("districts_id", "districts(id)", "RESTRICT", "RESTRICT")
	if errResort != nil {
		migration := model.MigrationHistory{
			DescMigration: "Add Table Resort",
			Date:          time.Now(),
		}

		db.Create(&migration)
	}

	errResort2 := db.Model(&resort).AddForeignKey("branch_id", "branches(id)", "RESTRICT", "RESTRICT")
	if errResort2 != nil {
		migration := model.MigrationHistory{
			DescMigration: "Add forgen key branch_id Table Resort",
			Date:          time.Now(),
		}

		db.Create(&migration)
	}

	resortRoom := model.ResortRoom{}
	errresortRoom := db.AutoMigrate(&resortRoom).AddForeignKey("resort_id", "resorts(id)", "RESTRICT", "RESTRICT")
	if errresortRoom != nil {
		migration := model.MigrationHistory{
			DescMigration: "Add Table resortRoom",
			Date:          time.Now(),
		}

		db.Create(&migration)
	}

	accessibilityResortRoom := model.AccessibilityResortRoom{}
	erraccessibilityResortRoom := db.AutoMigrate(&accessibilityResortRoom).AddForeignKey("resort_room_id", "resort_rooms(id)", "RESTRICT", "RESTRICT")
	if erraccessibilityResortRoom != nil {
		migration := model.MigrationHistory{
			DescMigration: "Add Table AccessibilityResortRoom",
			Date:          time.Now(),
		}

		db.Create(&migration)
	}

	erraccessibilityResortRoom2 := db.Model(&accessibilityResortRoom).AddForeignKey("accessibility_id", "accessibilities(id)", "RESTRICT", "RESTRICT")
	if erraccessibilityResortRoom2 != nil {
		migration := model.MigrationHistory{
			DescMigration: "Add Foregn key accesibility_id Table AccessibilityResortRoom",
			Date:          time.Now(),
		}

		db.Create(&migration)
	}


	amenitiesResortRoom := model.AmenitiesResortRoom{}
	erramenitiesResortRoom := db.AutoMigrate(&amenitiesResortRoom).AddForeignKey("resort_room_id", "resort_rooms(id)", "RESTRICT", "RESTRICT")
	if erramenitiesResortRoom != nil {
		migration := model.MigrationHistory{
			DescMigration: "Add Table AmenitiesResortRoom",
			Date:          time.Now(),
		}

		db.Create(&migration)
	}

	erramenitiesResortRoom2 := db.Model(&amenitiesResortRoom).AddForeignKey("amenities_id", "amenities(id)", "RESTRICT", "RESTRICT")
	if erramenitiesResortRoom2 != nil {
		migration := model.MigrationHistory{
			DescMigration: "Add Foregn key amenities_id Table AmenitiesResortRoom",
			Date:          time.Now(),
		}

		db.Create(&migration)
	}

	ResortPhoto := model.ResortPhoto{}
	errrResortPhoto := db.AutoMigrate(&ResortPhoto).AddForeignKey("resort_id", "resorts(id)", "RESTRICT", "RESTRICT")
	if errrResortPhoto != nil {
		migration := model.MigrationHistory{
			DescMigration: "Add Table ResortPhoto",
			Date:          time.Now(),
		}

		db.Create(&migration)
	}

	ResortRoomPayment := model.ResortRoomPayment{}
	errrResortRoomPayment := db.AutoMigrate(&ResortRoomPayment).AddForeignKey("resort_room_id", "resort_rooms(id)", "RESTRICT", "RESTRICT")
	if errrResortRoomPayment != nil {
		migration := model.MigrationHistory{
			DescMigration: "Add Table ResortRoomPayment",
			Date:          time.Now(),
		}

		db.Create(&migration)
	}

	ResortRoomPhoto := model.ResortRoomPhoto{}
	errrResortRoomPhoto := db.AutoMigrate(&ResortRoomPhoto).AddForeignKey("resort_room_id", "resort_rooms(id)", "RESTRICT", "RESTRICT")
	if errrResortRoomPhoto != nil {
		migration := model.MigrationHistory{
			DescMigration: "Add Table ResortRoomPhoto",
			Date:          time.Now(),
		}

		db.Create(&migration)
	}

}
