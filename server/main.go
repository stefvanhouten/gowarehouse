package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

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

func main() {
	router := gin.Default()
	router.GET("/products", getProducts)
	router.GET("/products/:id", getProductByID)
	router.POST("/products", postProducts)
	router.Run("localhost:8080")
}

func getProducts(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, products)
}

func postProducts(c *gin.Context) {
	var newProduct product

	if err := c.BindJSON(&newProduct); err != nil {
		return
	}

	products = append(products, newProduct)
	c.IndentedJSON(http.StatusCreated, newProduct)
}

func getProductByID(c *gin.Context) {
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
