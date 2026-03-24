package db

import (
	"database/sql"

	"alioth-hrc/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewMySQL(cfg config.MySQLConfig) (*gorm.DB, *sql.DB, error) {
	gdb, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}

	sqlDB, err := gdb.DB()
	if err != nil {
		return nil, nil, err
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)

	return gdb, sqlDB, nil
}
