package main

import (
	"database/sql"
	"os"
	"time"

	"example/logger"
	"github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

func main() {
	user := os.Getenv("DBUSER")
	database := "gowarehouse"

	logger.DefaultLogger.WithFields(logrus.Fields{
		"user":     user,
		"database": database,
	}).Debug("Setting up connection to database.")

	// Setup a database CONFIG with User/Passwd sourced from environment variables.
	cfg := mysql.Config{
		User:                 user,
		Passwd:               os.Getenv("DBPASSWORD"),
		Net:                  "tcp",
		Addr:                 "localhost:3306",
		DBName:               database,
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
			"user":     user,
			"database": database,
		}).Error("Error connecting to database:", err.Error())
		panic(err.Error())
	}

	//
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	logger.DefaultLogger.Info("Connected to database.")
}
