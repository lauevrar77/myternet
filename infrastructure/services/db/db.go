package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Connect(models ...interface{}) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("myternet.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Migrate the schema
	err = db.AutoMigrate(models...)
	if err != nil {
		return nil, err
	}

	return db, nil
}
