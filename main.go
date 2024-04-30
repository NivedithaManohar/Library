package main

import (
	"github.com/Library/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	routes.SetUpRouter(router)
	router.Run(":8989")
}
