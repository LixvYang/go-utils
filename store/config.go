package store

import (
	"errors"
	"fmt"

	"github.com/lixvyang/go-utils/db"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

func Connect(dialect, dsn string, gormCfg *gorm.Config) (*DB, error) {
	d, err := dialector(dialect, dsn)
	if err != nil {
		return nil, err
	}

	if gormCfg == nil {
		gormCfg = &gorm.Config{}
	}

	db, err := gorm.Open(d, gormCfg)
	if err != nil {
		return nil, err
	}

	return &DB{
		DB: db,
	}, nil
}

func Open(cfg db.Config, gormCfg *gorm.Config) (*DB, error) {
	dsn, err := cfg.DSN()
	if err != nil {
		return nil, err
	}

	db, err := Connect(cfg.Dialect, dsn, gormCfg)
	if err != nil {
		return nil, err
	}

	if cfg.ReadHost != "" && cfg.ReadHost != cfg.Host {
		readHostDSN, err := cfg.ReadHostDSN()
		if err != nil {
			return nil, err
		}
		rd, err := dialector(cfg.Dialect, readHostDSN)
		if err != nil {
			return nil, err
		}

		db.Use(dbresolver.Register(dbresolver.Config{
			Sources:  []gorm.Dialector{db.Dialector},
			Replicas: []gorm.Dialector{rd},
			Policy:   dbresolver.RandomPolicy{},
		}))
	}

	return db, nil
}

func MustOpen(cfg db.Config, gormCfg *gorm.Config) *DB {
	db, err := Open(cfg, gormCfg)
	if err != nil {
		panic(fmt.Errorf("open db failed: %w", err))
	}

	return db
}

func dialector(dialect string, dsn string) (gorm.Dialector, error) {
	var d gorm.Dialector
	switch dialect {
	case "sqlite", "sqlite3":
		d = sqlite.Open(dsn)
	case "mysql":
		d = mysql.Open(dsn)
	case "postgres":
		d = postgres.Open(dsn)
	default:
		return nil, fmt.Errorf("invalid dialect: %s", dialect)
	}

	return d, nil
}

func IsErrorNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}
