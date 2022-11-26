package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"

	"example/database" // Our local database module.
	"example/logger"   // Our local logging module.

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

type product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

var products = []product{
	{ID: 1, Name: "Product 1", Price: 10.00},
	{ID: 2, Name: "Product 2", Price: 20.00},
	{ID: 3, Name: "Product 3", Price: 30.00},
	{ID: 4, Name: "Product 4", Price: 40.00},
}

// Load required environment variables and put them in the config struct.
func loadEnv() *config {
	logger.DefaultLogger.Info("Loading environment variables into config.")

	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		panic(err.Error())
	}

	return &cfg
}

func main() {
	logger.DefaultLogger.Info("Starting server")
	cfg := loadEnv()
	env := &Env{
		cfg: cfg,
		db:  database.GetConnection(cfg),
	}
	router := gin.Default()

	router.GET("/products", env.getProducts)
	router.GET("/products/:id", env.getProductByID)
	router.POST("/products", env.postProducts)
	router.Run("localhost:8080")
}

func (env *Env) getProducts(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, products)
}

func (env *Env) postProducts(c *gin.Context) {
	var newProduct product

	if err := c.BindJSON(&newProduct); err != nil {
		return
	}

	products = append(products, newProduct)
	c.IndentedJSON(http.StatusCreated, newProduct)
}

func (env *Env) getProductByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	for _, p := range products {
		if p.ID == int(id) {
			c.IndentedJSON(http.StatusOK, p)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "product not found"})
}
