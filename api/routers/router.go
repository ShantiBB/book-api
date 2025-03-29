package routers

import "github.com/gin-gonic/gin"

func BookRoutes(router *gin.Engine, handler *Handler) {
	router.POST("/books", handler.CreateBook)
	router.GET("/books", handler.GetAllBooks)
	router.GET("/books/:id", handler.GetBookByID)
	router.PATCH("/books/:id", handler.UpdateBookByID)
	router.DELETE("/books/:id", handler.DeleteBookByID)
}
