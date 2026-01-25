package handlers

import (
	"fmt"
	"gin-starter/internal/models"
	"gin-starter/internal/store"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserHandler структура для обработчиков API
type UserHandler struct{}

// NewUserHandler создает новый экземпляр UserHandler
func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

// GetUsers обработчик для получения списка пользователей
func (h *UserHandler) GetUsers(c *gin.Context) {
	// Получаем доступ к базе данных из контекста
	dbStore, exists := c.Get("dbStore")
	if !exists {
		log.Println("Database connection not found in context")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not available"})
		return
	}

	// Приводим к нужному типу
	store := dbStore.(store.Store)

	// Получаем всех пользователей из базы данных
	users, err := store.GetUserRepo().GetAll()
	if err != nil {
		log.Printf("Error getting users: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get users"})
		return
	}

	c.JSON(http.StatusOK, users)
}

// CreateTestUsers обработчик для создания тестовых пользователей
func (h *UserHandler) CreateTestUsers(c *gin.Context) {
	// Получаем доступ к базе данных из контекста
	dbStore, exists := c.Get("dbStore")
	if !exists {
		log.Println("Database connection not found in context")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not available"})
		return
	}

	// Приводим к нужному типу
	store := dbStore.(store.Store)

	// Создаем тестовые данные
	testUsers := []models.User{
		{Name: "Иван Иванов", Email: "ivan@example.com"},
		{Name: "Мария Смирнова", Email: "maria@example.com"},
		{Name: "Алексей Попов", Email: "alexey@example.com"},
		{Name: "Елена Кузнецова", Email: "elena@example.com"},
		{Name: "Дмитрий Волков", Email: "dmitry@example.com"},
	}

	createdCount := 0
	for _, user := range testUsers {
		// Проверяем, существует ли уже пользователь с таким email
		existingUser, err := store.GetUserRepo().GetByEmail(user.Email)
		if err != nil || existingUser == nil {
			// Создаем нового пользователя
			if err := store.GetUserRepo().Create(&user); err != nil {
				log.Printf("Error creating test user: %v", err)
				continue
			}
			createdCount++
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Test users created successfully", "count": createdCount})
}

// CreateUser обработчик для создания нового пользователя
func (h *UserHandler) CreateUser(c *gin.Context) {
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

	c.JSON(http.StatusOK, user)
}

// DeleteUser обработчик для удаления пользователя
func (h *UserHandler) DeleteUser(c *gin.Context) {
	// Получаем ID пользователя из параметров
	userID := c.Param("id")

	// Получаем доступ к базе данных из контекста
	dbStore, exists := c.Get("dbStore")
	if !exists {
		log.Println("Database connection not found in context")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not available"})
		return
	}

	// Приводим к нужному типу
	store := dbStore.(store.Store)

	// Преобразуем ID в uint
	var id uint
	_, err := fmt.Sscanf(userID, "%d", &id)
	if err != nil {
		log.Printf("Error parsing user ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Удаляем пользователя из базы данных
	if err := store.GetUserRepo().Delete(id); err != nil {
		log.Printf("Error deleting user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}