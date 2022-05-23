package db

import (
	"fmt"

	"github.com/nightlord189/ulp/internal/config"
	"github.com/nightlord189/ulp/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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
	dbLogger := logger.Discard
	if cfg.DB.Log {
		dbLogger = logger.Default.LogMode(logger.Info)
	}
	db, err := gorm.Open(postgres.Open(address), &gorm.Config{
		Logger: dbLogger,
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
	err := db.AutoMigrate(
		model.TaskDB{},
		model.UserDB{},
		model.AttemptDB{},
	)
	if err != nil {
		return err
	}
	db.Exec("ALTER TABLE tasks ADD CONSTRAINT tasks_fk FOREIGN KEY (creator_id) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE")
	db.Exec("ALTER TABLE attempts ADD CONSTRAINT tasks_fk FOREIGN KEY (task_id) REFERENCES tasks (id) ON DELETE SET NULL ON UPDATE CASCADE")
	db.Exec("ALTER TABLE attempts ADD CONSTRAINT users_fk FOREIGN KEY (creator_id) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE")
	return nil
}
