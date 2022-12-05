package db

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/kevin-luvian/moonlay-test/model"
	"github.com/kevin-luvian/moonlay-test/pkg/settings"
)

func NewPSQL() *gorm.DB {
	setting := settings.DatabaseSetting
	dbUri := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", setting.Host, setting.Port, setting.User, setting.Name, setting.Password)

	fmt.Println(dbUri)
	db, err := gorm.Open("postgres", dbUri)
	if err != nil {
		fmt.Println("storage err: ", err)
	}

	db.DB().SetMaxIdleConns(3)
	db.LogMode(true)
	return db
}

func New() *gorm.DB {
	switch settings.DatabaseSetting.Type {
	case "postgres":
		return NewPSQL()
	}

	log.Fatal("invalid db config")
	return nil
}

func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(
		&model.List{},
	)
}
