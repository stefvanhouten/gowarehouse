package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"

	"example/database" // Our local database module.
	"example/logger"   // Our local logging module.
	"example/models"   // Our local models module.

	"github.com/caarlos0/env/v6"
)

type Env struct {
	db  *sql.DB
	cfg *config
}

type config struct {
	User        string `env:"DBUSER"`
	Password    string `env:"DBPASSWORD"`
	Database    string `env:"DATABASE"`
	Environment string `env:"ENVIRONMENT" envDefault:"DEV"`
	LogDir      string `env:"LOGDIR"`
}

func (c config) GetUser() string {
	return c.User
}

func (c config) GetPassword() string {
	return c.Password
}

func (c config) GetDatabase() string {
	return c.Database
}

func (c config) GetEnvironment() string {
	return c.Environment
}

func (c config) GetLogDir() string {
	return c.LogDir
}

// Load required environment variables and put them in the config struct.
func loadEnv() *config {
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		panic(err.Error())
	}

	return &cfg
}

func main() {
	cfg := loadEnv()
	logger.Setup(cfg)
	logger.DefaultLogger.Info("Starting server")

	env := &Env{
		cfg: cfg,
		db:  database.GetConnection(cfg),
	}
	router := gin.Default()

	router.GET("/products", env.getProducts)
	router.GET("/products/:id", env.getProductByID)
	// router.POST("/products", env.postProducts)
	router.Run("localhost:8080")
}

func (env *Env) getProductByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	product, err := products.ProductByID(env.db, id)

	if err != nil {
		logger.DefaultLogger.Debug("Error getting product:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if product == nil {
		c.IndentedJSON(
			http.StatusNotFound,
			gin.H{"error": fmt.Sprintf("No product found for ID: '%d'", id)})
		return
	}

	c.IndentedJSON(http.StatusOK, product)
}

func (env *Env) getProducts(c *gin.Context) {
	products, err := products.AllProducts(env.db)

	if err != nil {
		logger.DefaultLogger.Error("Error getting products:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, products)
}

// func (env *Env) postProducts(c *gin.Context) {
// 	var newProduct product

// 	if err := c.BindJSON(&newProduct); err != nil {
// 		return
// 	}

// 	products = append(products, newProduct)
// 	c.IndentedJSON(http.StatusCreated, newProduct)
// }
