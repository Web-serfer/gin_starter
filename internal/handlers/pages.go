package handlers

import (
	"gin-starter/templates"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HomeHandler(c *gin.Context) {
	canonicalURL := c.Request.URL.Scheme + "://" + c.Request.Host + c.Request.URL.Path
	menuItems := templates.GetDefaultMenuItems()
	c.Status(http.StatusOK)
	if err := templates.IndexPage(canonicalURL, menuItems).Render(c.Request.Context(), c.Writer); err != nil {
		log.Printf("Template render error: %v", err)
		c.String(http.StatusInternalServerError, "Internal Server Error")
	}
}

func AboutHandler(c *gin.Context) {
	canonicalURL := c.Request.URL.Scheme + "://" + c.Request.Host + c.Request.URL.Path
	menuItems := templates.GetDefaultMenuItems()
	c.Status(http.StatusOK)
	if err := templates.AboutPage(canonicalURL, menuItems).Render(c.Request.Context(), c.Writer); err != nil {
		log.Printf("Template render error: %v", err)
		c.String(http.StatusInternalServerError, "Internal Server Error")
	}
}

func ContactHandler(c *gin.Context) {
	canonicalURL := c.Request.URL.Scheme + "://" + c.Request.Host + c.Request.URL.Path
	menuItems := templates.GetDefaultMenuItems()
	c.Status(http.StatusOK)
	if err := templates.ContactPage(canonicalURL, menuItems).Render(c.Request.Context(), c.Writer); err != nil {
		log.Printf("Template render error: %v", err)
		c.String(http.StatusInternalServerError, "Internal Server Error")
	}
}

// NotFoundHandler обработчик для страницы 404
func NotFoundHandler(c *gin.Context) {
	// Получаем меню
	menuItems := templates.GetDefaultMenuItems()

	// Формируем canonical URL
	canonicalURL := c.Request.URL.Scheme + "://" + c.Request.Host + c.Request.URL.Path

	// Устанавливаем статус 404 Not Found
	c.Status(http.StatusNotFound)

	// Рендерим шаблон
	if err := templates.NotFoundPage(canonicalURL, menuItems).Render(c.Request.Context(), c.Writer); err != nil {
		log.Printf("Template render error: %v", err)
		c.String(http.StatusInternalServerError, "Internal Server Error")
	}
}
