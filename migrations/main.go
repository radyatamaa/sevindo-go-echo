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

	accessibilityResortRoom := model.AccessibilityResort{}
	erraccessibilityResortRoom := db.AutoMigrate(&accessibilityResortRoom).AddForeignKey("resort_id", "resorts(id)", "RESTRICT", "RESTRICT")
	if erraccessibilityResortRoom != nil {
		migration := model.MigrationHistory{
			DescMigration: "Add Table AccessibilityResort",
			Date:          time.Now(),
		}

		db.Create(&migration)
	}

	articleblog := model.ArticleBlog{}
	errarticleblog := db.AutoMigrate(&articleblog)
	if errarticleblog != nil {
		migration := model.MigrationHistory{
			DescMigration: "Add Table Article Blog",
			Date:          time.Now(),
		}

		db.Create(&migration)
	}

	articlecategory := model.ArticleCategory{}
	errarticlecategory := db.AutoMigrate(&articlecategory)
	if errarticlecategory != nil {
		migration := model.MigrationHistory{
			DescMigration: "Add Table Article Category",
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

	amenitiesResortRoom := model.AmenitiesResort{}
	erramenitiesResortRoom := db.AutoMigrate(&amenitiesResortRoom).AddForeignKey("resort_id", "resorts(id)", "RESTRICT", "RESTRICT")
	if erramenitiesResortRoom != nil {
		migration := model.MigrationHistory{
			DescMigration: "Add Table AmenitiesResort",
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

	bank := model.Bank{}
	errbank := db.AutoMigrate(&bank)
	if errbank != nil {
		migration := model.MigrationHistory{
			DescMigration: "Add Table Bank",
			Date:          time.Now(),
		}

		db.Create(&migration)
	}

	booking := model.Booking{}
	errBooking := db.AutoMigrate(&booking).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
	if errBooking != nil {
		migration := model.MigrationHistory{
			DescMigration: "Add Table Booking",
			Date:          time.Now(),
		}

		db.Create(&migration)
	}

	errBooking1 := db.Model(&booking).AddForeignKey("resort_id", "resorts(id)", "RESTRICT", "RESTRICT")
	if errBooking1 != nil {
		migration := model.MigrationHistory{
			DescMigration: "add foregn key resort_id in table booking",
			Date:          time.Now(),
		}

		db.Create(&migration)
	}

	errBooking2 := db.Model(&booking).AddForeignKey("resort_room_id", "resort_rooms(id)", "RESTRICT", "RESTRICT")
	if errBooking2 != nil {
		migration := model.MigrationHistory{
			DescMigration: "add foregn key resort_room_id in table booking",
			Date:          time.Now(),
		}

		db.Create(&migration)
	}

	promo := model.Promo{}
	errpromo := db.AutoMigrate(&promo)
	if errpromo != nil {
		migration := model.MigrationHistory{
			DescMigration: "Add Table Promo",
			Date:          time.Now(),
		}

		db.Create(&migration)
	}

	paymentMethod := model.PaymentMethod{}
	errpaymentMethod := db.AutoMigrate(&paymentMethod)
	if errpaymentMethod != nil {
		migration := model.MigrationHistory{
			DescMigration: "Add Table PaymentMethod",
			Date:          time.Now(),
		}

		db.Create(&migration)
	}

	GalleryExperience := model.GalleryExperience{}
	errGalleryExperience := db.AutoMigrate(&GalleryExperience)
	if errGalleryExperience != nil {
		migration := model.MigrationHistory{
			DescMigration: "Add Table GalleryExperience",
			Date:          time.Now(),
		}

		db.Create(&migration)
	}

	transaction := model.Transaction{}
	errTransaction := db.AutoMigrate(&transaction).AddForeignKey("booking_id", "bookings(id)", "RESTRICT", "RESTRICT")
	if errTransaction != nil {
		migration := model.MigrationHistory{
			DescMigration: "Add Table Transaction",
			Date:          time.Now(),
		}

		db.Create(&migration)
	}

	errTransaction1 := db.Model(&transaction).AddForeignKey("promo_id", "promos(id)", "RESTRICT", "RESTRICT")
	if errTransaction1 != nil {
		migration := model.MigrationHistory{
			DescMigration: "add foregn key promo_id in table Transaction",
			Date:          time.Now(),
		}

		db.Create(&migration)
	}

	errTransaction2 := db.Model(&transaction).AddForeignKey("payment_method_id", "payment_methods(id)", "RESTRICT", "RESTRICT")
	if errTransaction2 != nil {
		migration := model.MigrationHistory{
			DescMigration: "add foregn key payment_method_id in table Transaction",
			Date:          time.Now(),
		}

		db.Create(&migration)
	}

	errTransaction3 := db.Model(&transaction).AddForeignKey("resort_room_payment", "resort_room_payments(id)", "RESTRICT", "RESTRICT")
	if errTransaction3 != nil {
		migration := model.MigrationHistory{
			DescMigration: "add foregn key resort_room_payment in table Transaction",
			Date:          time.Now(),
		}

		db.Create(&migration)
	}

	review := model.Review{}
	errreview := db.AutoMigrate(&review).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
	if errreview != nil {
		migration := model.MigrationHistory{
			DescMigration: "Add Table Review",
			Date:          time.Now(),
		}

		db.Create(&migration)
	}

	errreview1 := db.Model(&review).AddForeignKey("transaction_id", "transactions(id)", "RESTRICT", "RESTRICT")
	if errreview1 != nil {
		migration := model.MigrationHistory{
			DescMigration: "add foregn key transaction_id in table Review",
			Date:          time.Now(),
		}

		db.Create(&migration)
	}

}
