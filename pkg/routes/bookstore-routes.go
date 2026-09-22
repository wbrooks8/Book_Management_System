package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wbrook8/book_management_system/pkg/controllers"
)

func RegisterRoutes(server *gin.Engine, handler *controllers.Handler) {
	server.POST("/book/", handler.CreateBook)
	server.GET("/book/", handler.GetBooks)
	server.GET("/book/:bookId", handler.GetBookByID)
	server.PUT("/book/:bookId", handler.UpdateBook)
	server.DELETE("/book/:bookId", handler.DeleteBook)
}
