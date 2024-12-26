package store

import (
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

type DB struct {
	*gorm.DB
}

func (db *DB) Close() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (db *DB) Update() *gorm.DB {
	return db.Clauses(dbresolver.Write)
}

func (db *DB) View() *gorm.DB {
	return db.Clauses(dbresolver.Read)
}

func (db *DB) Tx(f func(tx *DB) error) error {
	return db.Transaction(func(tx *gorm.DB) error {
		return f(&DB{DB: tx})
	})
}