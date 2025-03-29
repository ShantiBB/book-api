package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	_ "modernc.org/sqlite"

	"book/api/routers"
	"book/internal/config"
	"book/internal/lib/sl"
	bookQuery "book/internal/storage/book"
	"book/internal/storage/sqlite"
)

func main() {
	cfg := config.MustLoad()
	log := sl.SetupLogger(cfg.Env)

	log.Debug("Debug is true")

	db, err := sqlite.SessionDB(cfg.StoragePath)
	if err != nil {
		log.Error("Failed to initialize database", err)
		os.Exit(1)
	}
	log.Debug("Database initialized successfully")

	bookStorage := &bookQuery.SQLiteBookStorage{DB: db}
	handler := routers.NewHandler(bookStorage)

	router := gin.Default()
	router.POST("/books", handler.CreateBook)
	router.GET("/books", handler.GetAllBooks)
	router.GET("/books/:id", handler.GetBookByID)
	router.DELETE("/books/:id", handler.DeleteBookByID)

	addr := fmt.Sprintf("%s:%s", cfg.HTTPServer.Address, cfg.HTTPServer.Port)
	if err := router.Run(addr); err != nil {
		log.Error("failed to start server: %v", err)
	}
}
