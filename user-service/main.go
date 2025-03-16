package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
)

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func main() {
	app := fiber.New(fiber.Config{
		AppName: "user-service",
	})

	// In-memory user store for demo
	users := make(map[string]User)

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "healthy",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// Get all users
	app.Get("/users", func(c *fiber.Ctx) error {
		userList := make([]User, 0, len(users))
		for _, user := range users {
			userList = append(userList, user)
		}
		return c.JSON(userList)
	})

	// Get user by ID
	app.Get("/users/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		if user, exists := users[id]; exists {
			return c.JSON(user)
		}
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	})

	// Create new user
	app.Post("/users", func(c *fiber.Ctx) error {
		var user User
		if err := c.BodyParser(&user); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
		}
		user.ID = time.Now().Format("20060102150405")
		users[user.ID] = user
		return c.Status(201).JSON(user)
	})

	// Setup graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := app.Listen(":3001"); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
