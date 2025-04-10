package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"book/api/routers"
	"book/api/services"
	"book/internal/config"
	"book/internal/lib/sl"
	"book/internal/models"
	bookQuery "book/internal/storage/book"
)

func main() {
	cfg := config.MustLoad()
	log := sl.SetupLogger(cfg.Env)

	log.Debug("Debug is true")
	log.Debug("Using storage path: %s", cfg.StoragePath)
	storageDir := filepath.Dir(cfg.StoragePath)
	if err := os.MkdirAll(storageDir, os.ModePerm); err != nil {
		log.Error("Failed to create storage directory", err)
		os.Exit(1)
	}

	db, err := gorm.Open(sqlite.Open(cfg.StoragePath))
	if err != nil {
		log.Error("Failed to initialize database", err)
		os.Exit(1)
	}
	log.Debug("Database initialized successfully")

	err = db.AutoMigrate(&models.Book{})
	if err != nil {
		log.Error("Failed to migrate database", err)
		return
	}

	router := gin.Default()

	bookStorage := &bookQuery.SQLiteBookStorage{DB: db}
	handler := services.NewHandler(bookStorage)

	routers.BookRoutes(router, handler)

	addr := fmt.Sprintf("%s:%s", cfg.HTTPServer.Address, cfg.HTTPServer.Port)
	if err := router.Run(addr); err != nil {
		log.Error("failed to start server: %v", err)
	}
}
