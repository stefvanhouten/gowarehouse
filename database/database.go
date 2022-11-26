package main

import (
	"database/sql"
	"time"

	"example/logger" // Our local logging module.
	"github.com/caarlos0/env/v6"
	"github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

type config struct {
	User     string `env:"DBUSER"`
	Password string `env:"DBPASSWORD"`
	Database string `env:"DATABASE"`
}

// Load required environment variables and put them in the config struct.
func loadEnv() config {
	logger.DefaultLogger.Info("Loading environment variables.")

	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		panic(err.Error())
	}

	return cfg
}

// Connect to the MYSQL database and return the connection.
func connect(conf *config) *sql.DB {
	// Setup a database CONFIG with User/Passwd sourced from environment variables.
	cfg := mysql.Config{
		User:                 conf.User,
		Passwd:               conf.Password,
		Net:                  "tcp",
		Addr:                 "localhost:3306",
		DBName:               conf.Database,
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
			"user":     conf.User,
			"database": conf.Database,
		}).Error("Error connecting to database:", err.Error())
		panic(err.Error())
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return db
}

func main() {
	conf := loadEnv()
	connect(&conf)

	// Log the user and database we are attempting to connect to.
	logger.DefaultLogger.WithFields(logrus.Fields{
		"user":     conf.User,
		"database": conf.Database,
	}).Debug("Setting up connection to database.")

	logger.DefaultLogger.Info("Connected to database.")
}
