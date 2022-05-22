package db

import (
	"fmt"

	"github.com/nightlord189/ulp/internal/config"
	"github.com/nightlord189/ulp/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	gormLogger "gorm.io/gorm/logger"
)

var ErrRecordNotFound = gorm.ErrRecordNotFound

// Manager - слой БД. Предоставляет простые атомарные методы работы с таблицами GORM
type Manager struct {
	Config *config.Config
	DB     *gorm.DB
}

// InitDb - конструктор DBManager
func InitDb(cfg *config.Config) (*Manager, error) {
	dbInstance, err := connectDb(cfg)
	dbManager := Manager{
		DB:     dbInstance,
		Config: cfg,
	}
	return &dbManager, err
}

// ConnectDb - возвращает инстанс gorm, подключенный к postgres
func connectDb(cfg *config.Config) (*gorm.DB, error) {
	address := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable password=%s", cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Name, cfg.DB.Password)
	db, err := gorm.Open(postgres.Open(address), &gorm.Config{
		Logger: gormLogger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic("Failed to connect database: " + err.Error())
	} else {
		db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
		if cfg.DB.Migrate {
			err = autoMigrate(db)
			if err != nil {
				err = fmt.Errorf("migrate error: %w", err)
			}
		}
	}
	return db, err
}

// AutoMigrate - применить миграции
func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		model.TaskDB{},
		model.UserDB{},
		model.AttemptDB{},
	)
}
