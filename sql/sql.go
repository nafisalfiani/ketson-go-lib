package sql

import (
	"context"
	"fmt"

	"github.com/nafisalfiani/ketson-go-lib/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	Url      string
	Port     string
	Name     string
	User     string
	Password string
}

func Init(cfg Config, log log.Interface) (*gorm.DB, error) {
	log.Info(context.Background(), "connecting to postgresql database...")
	dsn := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable", cfg.Url, cfg.Port, cfg.User, cfg.Password, cfg.Name)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:         logger.Default.LogMode(logger.Info),
		TranslateError: true,
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}
