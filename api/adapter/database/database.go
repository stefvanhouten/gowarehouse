package database

import (
	"database/sql"
	"time"

	"example/logger" // Our local logging module.
	"github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

// TODO: Add functionality to re-use open connections instead of initializing a new one
// every time we want to access the database.

// Interface for the datbase config.
type DatabaseConfig interface {
	GetUser() string
	GetPassword() string
	GetDatabase() string
}

// Connect to the MYSQL database and return the connection.
func connect(conf DatabaseConfig) *sql.DB {
	// Setup a database CONFIG with User/Passwd sourced from environment variables.
	logger.DefaultLogger.WithFields(logrus.Fields{
		"user":     conf.GetUser(),
		"database": conf.GetDatabase(),
	}).Debug("Setting up connection to database.")

	cfg := mysql.Config{
		User:                 conf.GetUser(),
		Passwd:               conf.GetPassword(),
		Net:                  "tcp",
		Addr:                 "localhost:3306",
		DBName:               conf.GetDatabase(),
		AllowNativePasswords: true,
	}

	db, err := sql.Open(
		"mysql",
		cfg.FormatDSN(),
	)

	if err != nil {
		panic(err.Error())
	}

	err = db.Ping()
	if err != nil {
		logger.DefaultLogger.WithFields(logrus.Fields{
			"user":     conf.GetUser(),
			"database": conf.GetDatabase(),
		}).Error("Error connecting to database:", err.Error())
		panic(err.Error())
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	logger.DefaultLogger.Info("Connected to database.")
	return db
}

func GetConnection(conf DatabaseConfig) *sql.DB {
	return connect(conf)
}
