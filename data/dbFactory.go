package data

import (
	"learn-viper/config"

	"github.com/golang/glog"
	"github.com/jinzhu/gorm"
	// importing all possible database dialect
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// DBFactory struct
type DBFactory struct {
	config config.DatabaseConfiguration
}

// NewDbFactory instantiate new DB Factory object
func NewDbFactory(cfg config.DatabaseConfiguration) *DBFactory {
	return &DBFactory{config: cfg}
}

// DBConnection get open database connection
func (f *DBFactory) DBConnection() (*gorm.DB, error) {
	db, err := gorm.Open(f.config.DbType, f.config.ConnectionUri)
	if err != nil {
		glog.Fatalf("Failed to connect to database: %s", err)
		return nil, err
	}

	return db, nil
}
