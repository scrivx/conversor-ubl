package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/scrivx/conversor-ubl/internal/api/routes"
)

func main() {
	r := gin.Default()

	routes.RegisterRoutes(r)

	log.Println("🚀Iniciando servidor en http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}

}
