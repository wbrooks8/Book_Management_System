package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/wbrook8/book_management_system/pkg/config"
	"github.com/wbrook8/book_management_system/pkg/controllers"
	"github.com/wbrook8/book_management_system/pkg/models"
	"github.com/wbrook8/book_management_system/pkg/routes"
)

func main() {
	db, err := config.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer config.Close(db)

	books := models.NewRepository(db)
	if err := books.Migrate(); err != nil {
		log.Fatal(err)
	}

	server := gin.Default()
	routes.RegisterRoutes(server, controllers.NewHandler(books))

	if err := server.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
