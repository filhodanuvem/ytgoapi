package http

import (
	"github.com/gin-gonic/gin"
)

func SetRoutes(g *gin.Engine) {
	g.POST("/posts", PostPosts)
	g.DELETE("/posts/:id", DeletePosts)
	g.GET("/posts/:id", GetPosts)
}
