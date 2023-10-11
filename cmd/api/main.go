package main

import (
	"fmt"

	"github.com/filhodanuvem/ytgoapi/internal/configs"
	"github.com/filhodanuvem/ytgoapi/internal/database"
	"github.com/filhodanuvem/ytgoapi/internal/http"
	"github.com/gin-gonic/gin"
)

func main() {
	connectionString := configs.ReadConfig()

	fmt.Println(connectionString)
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
