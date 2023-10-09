package main

import (
	"github.com/filhodanuvem/ytgoapi/internal/database"

	"github.com/filhodanuvem/ytgoapi/internal/http"
	"github.com/gin-gonic/gin"
)

func main() {
	connectionString := "postgresql://posts:p0stgr3s@localhost:5432/posts"
	conn, err := database.NewConnection(connectionString)
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	g := gin.Default()
	http.Configure()
	http.SetRoutes(g)
	g.Run(":3000")
}
