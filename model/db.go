package model

import (
	"fmt"
	"log"
	"os"

	"gorm.io/gorm"
	"gorm.io/driver/postgres"
	"github.com/joho/godotenv"
)


// GetDB retun Db connetion
func GetDB() (*gorm.DB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dbUser := os.Getenv("db_user")
	dbPassword := os.Getenv("db_pass")
	dbName := os.Getenv("db_name")
	dbHost := os.Getenv("db_host")
	//dbDriver := os.Getenv("db_driver")
	dsn := "host=" + dbHost + " user=" + dbUser + " password=" + dbPassword + " dbname=" + dbName +" port=5432 sslmode=disable TimeZone=Asia/Tehran"
	fmt.Println(dsn)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	//db, err := gorm.Open(postgres.New(postgres.Config{
  //DSN: dsn,
  //PreferSimpleProtocol: true, // disables implicit prepared statement usage
//}), &gorm.Config{})
	//db, err := gorm.Open(dbDriver, dbUser+":"+dbPassword+"@/"+dbName+"?charset=utf8mb4&parseTime=True")
	//dbString := "postgres://"+dbUser+":"+dbPassword+"@localhost/schema?sslmode=disable"

	//db, err := gorm.Open(dbDriver, dbString)

	if err != nil {
		return nil, err
	}
	return db, nil
}
