package pkg

import (
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/guillaumejacquart/go-http-scheduler/pkg/domain"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/spf13/viper"
)

var db *gorm.DB

func initDb() {
	dbType := viper.GetString("db.type")
	dbPort := viper.GetInt("db.port")
	dbHost := viper.GetString("db.host")
	dbUsername := viper.GetString("db.username")
	dbPassword := viper.GetString("db.password")
	dbName := viper.GetString("db.name")
	dbPath := viper.GetString("db.path")

	var connectionString string
	switch dbType {
	case "mysql":
		connectionString = fmt.Sprintf(
			"%v:%v@tcp(%v:%v)/%v?parseTime=true",
			dbUsername,
			dbPassword,
			dbHost,
			dbPort,
			dbName)
	case "postgres":
		connectionString = fmt.Sprintf(
			"host=%v port=%v user=%v dbname=%v sslmode=disable password=%v",
			dbHost,
			dbPort,
			dbUsername,
			dbName,
			dbPassword)
	case "sqlite3":
		connectionString = dbPath
	}

	log.Println("Connection to", connectionString)
	dbInit, err := gorm.Open(dbType, connectionString)
	if err != nil {
		panic(err)
	}

	log.Print("Connected !")
	db = dbInit

	// Migrate the schema
	log.Println("Migrating schema for ORM ...")
	db.AutoMigrate(&domain.App{}, &domain.History{}, &domain.Header{})
	log.Println("Schema migrated !")
}

func getAllApps() ([]domain.App, error) {
	var apps []domain.App
	err := db.Preload("Headers").Find(&apps).Error

	for _, a := range apps {
		_, exists := appsJob[a.ID]
		a.JobStatus = exists
	}

	return apps, err
}

func getApp(id uint) (domain.App, error) {
	app := domain.App{}
	err := db.Preload("Headers").First(&app, id).Error

	_, exists := appsJob[app.ID]
	app.JobStatus = exists

	return app, err
}

func insertApp(app *domain.App) error {
	return db.Create(app).Error
}

func insertHistory(history domain.History) error {
	return db.Create(&history).Error
}

func getAppHistory(appID uint) ([]domain.History, error) {
	histories := []domain.History{}
	err := db.Order("date desc").Limit(5).Where("app_id = ?", appID).Find(&histories).Error

	return histories, err
}

func updateApp(id uint, app domain.App) error {
	existingApp, err := getApp(id)
	if err != nil {
		return err
	}

	existingApp.URL = app.URL
	existingApp.Method = app.Method
	existingApp.Body = app.Body
	existingApp.CronExpression = app.CronExpression

	db.Model(existingApp).Association("Headers").Clear()
	existingApp.Headers = app.Headers

	return db.Save(&existingApp).Error
}

func updateAppStatus(id uint, status string) error {
	existingApp, err := getApp(id)
	if err != nil {
		return err
	}

	existingApp.Status = status
	if status == "up" {
		now := time.Now()
		existingApp.LastUpDate = &now
	}

	return db.Save(&existingApp).Error
}

func deleteApp(id uint) error {
	app, err := getApp(id)
	if err != nil {
		return err
	}

	return db.Delete(&app).Error
}
