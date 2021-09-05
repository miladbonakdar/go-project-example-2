package sql

import (
	"github.com/jinzhu/gorm"
	"hotel-engine/core/dbmodel"
	"hotel-engine/infrastructure/logger"
)

func newSqlClient(connectionString string) *gorm.DB {
	db, err := gorm.Open("mssql", connectionString)

	if err != nil {
		logger.WithException(err).
			Fatal("Error creating connection pool")
	}
	logger.Info("DB Connected!")
	db.SetLogger(&logger.GormLogger{})
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(10)
	db.AutoMigrate(&dbmodel.Amenity{}, &dbmodel.Hotel{}, &dbmodel.City{},
		&dbmodel.Place{}, &dbmodel.OrderRoom{}, &dbmodel.Order{}, &dbmodel.AmenityCategory{},
		&dbmodel.Badge{}, &dbmodel.FAQ{})
	return db
}

// InitDatabase returns an implementation of the sql database orm with given connection string.
func InitDatabase(connectionString string) *gorm.DB {
	db := newSqlClient(connectionString)
	return db
}
