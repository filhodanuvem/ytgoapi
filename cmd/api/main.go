package main

import (
	"github.com/filhodanuvem/ytgoapi/config"
	"github.com/filhodanuvem/ytgoapi/internal/database"

	"github.com/filhodanuvem/ytgoapi/internal/http"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, _ := config.ReadConfig()
	conn, err := database.NewConnection(cfg.GetPostgresConnectionString())
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	g := gin.Default()
	http.Configure()
	http.SetRoutes(g)
	g.Run(":3000")
}
