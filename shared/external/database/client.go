package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/cchristian77/wallet-service/util/config"
	"github.com/cchristian77/wallet-service/util/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// openDB opens a database connection using pgx driver
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

// OpenGormDB initializes a gorm.DB instance from sql.DB
func OpenGormDB(sqlDB *sql.DB) (*gorm.DB, error) {
	return gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{
		PrepareStmt: true,
	})
}

// ConnectToDB establishes a connection to the PostgreSQL database using the configured DSN and retries on failure.
func ConnectToDB() *sql.DB {
	var counts int

	// get dsn from env.json
	dsn := config.Env().DSN()

	// get dsn from docker
	//dsn := os.Getenv("DSN")
	for {
		connection, err := openDB(dsn)
		if err != nil {
			logger.L().Warn(fmt.Sprintf("Postgres not yet ready: %v", err))
			counts++
		} else {
			logger.L().Info("Connected to Postgres")
			return connection
		}

		if counts > 10 {
			logger.L().Error(fmt.Sprintf("Failed to connect to Postgres: %v", err))
			return nil
		}

		logger.L().Info("backing off for two seconds...")
		time.Sleep(2 * time.Second)
	}
}
