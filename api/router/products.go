package router

import (
	"fmt"
	"net/http"
	"strconv"

	"example/api/models"
	"example/logger"

	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v9"
)

func (env *Env) getProducts(c *gin.Context) {
	products, err := products.AllProducts(env.db)

	if err != nil {
		logger.DefaultLogger.Error("Error getting products:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, products)
}

func (env *Env) getProductByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	product, err := products.ProductByID(env.db, int64(id))

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

func (env *Env) postProducts(c *gin.Context) {
	var newProduct products.Product

	// Load POST data into Product object.
	if err := c.BindJSON(&newProduct); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create a validator and validate the Product object.
	v := validator.New()
	if err := v.Struct(newProduct); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	// Insert the product into the database.
	res, err := products.CreateProduct(env.db, &newProduct)
	if err != nil {
		logger.DefaultLogger.Error("Error creating product:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Load the newly created product from the database
	prod, err := products.ProductByID(env.db, res)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.IndentedJSON(http.StatusCreated, prod)
}
