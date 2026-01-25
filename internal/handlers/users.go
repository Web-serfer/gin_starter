package handlers

import (
	"gin-starter/internal/models"
	"gin-starter/internal/store"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateUserHandler обработчик для создания нового пользователя
func CreateUserHandler(c *gin.Context) {
	// Получаем доступ к базе данных из контекста
	dbStore, exists := c.Get("dbStore")
	if !exists {
		log.Println("Database connection not found in context")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not available"})
		return
	}

	// Приводим к нужному типу
	store := dbStore.(store.Store)

	// Получаем данные из формы
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Создаем пользователя в базе данных
	if err := store.GetUserRepo().Create(&user); err != nil {
		log.Printf("Error creating user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully", "user": user})
}
