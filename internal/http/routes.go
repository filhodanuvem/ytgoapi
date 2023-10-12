package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetRoutes(g *gin.Engine) {
	g.POST("/posts", PostPosts)
	g.DELETE("/posts/:id", DeletePosts)
	g.GET("/posts/:id", GetPosts)
	g.GET("/posts", GetAll)
	g.PUT("/posts/:id", Update)

	g.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
}
