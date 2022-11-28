package router

import (
	"database/sql"

	"example/api/adapter/database"
	"example/api/types"

	"github.com/gin-gonic/gin"
)

type Env struct {
	db  *sql.DB
	cfg *types.Config
}

// Set the Content-Type of all responses to JSON.
func JSONMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Next()
	}
}

func New(cfg *types.Config) *gin.Engine {
	env := &Env{
		db:  database.GetConnection(cfg),
		cfg: cfg,
	}

	router := gin.Default()
	router.Use(JSONMiddleware())

	router.GET("/products", env.getProducts)
	router.GET("/products/:id", env.getProductByID)
	router.POST("/products", env.postProducts)

	return router
}
