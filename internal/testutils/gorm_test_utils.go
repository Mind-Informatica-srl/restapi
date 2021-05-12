package testutils

import (
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// SetupTestForGorm provide a fake connection pool and the mock to verify the test
func SetupTestForGorm() (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})

	connectionPool, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return connectionPool, mock
}
