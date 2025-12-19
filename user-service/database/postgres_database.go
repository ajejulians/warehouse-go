package database

import (
	"fmt"
	"warehouse-go/user-service/configs"
	"warehouse-go/user-service/model"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Postgres struct {
	DB *gorm.DB
}

func ConnectPostgres(cfg configs.Config) (*Postgres, error) {
	connString := fmt.Sprintf(
	"postgres://%s:%s@%s:%s/%s",
	cfg.SqlDB.User,
	cfg.SqlDB.Password,
	cfg.SqlDB.Host,
	cfg.SqlDB.Port,
	cfg.SqlDB.DBname,
)
	db, err := gorm.Open(postgres.Open(connString), &gorm.Config{})
	if err != nil {
		log.Errorf("[Postgres] ConnectionPostgres - 1: %v", err)
		return nil, err
	}

	db.AutoMigrate(&model.User{}, &model.Role{}, &model.UserRole{})
	sqlDB, err := db.DB()
	if err != nil {
		log.Errorf("[Postgres] Connection Postgres - 2: %v", err)
		return nil, err
	}

	SeedRole(db)
	SeedManager(db)
	
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	return &Postgres{DB: db}, nil
}