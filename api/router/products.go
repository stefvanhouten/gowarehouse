package router

import (
	"fmt"
	"net/http"
	"strconv"

	"example/api/models"
	"example/logger"

	"github.com/gin-gonic/gin"
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

// func (env *Env) postProducts(c *gin.Context) {
// 	var newProduct product

// 	if err := c.BindJSON(&newProduct); err != nil {
// 		return
// 	}

// 	products = append(products, newProduct)
// 	c.IndentedJSON(http.StatusCreated, newProduct)
// }
