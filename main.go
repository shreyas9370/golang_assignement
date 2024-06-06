package main

import (
	"context"
	"database/sql"
	"golang-assignment/config"
	"golang-assignment/handlers"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	var err error
	config.DB, err = sql.Open("mysql", "root:Shreyas9370@@tcp(localhost:3306)/golang_assignement")
	if err != nil {
		log.Fatal("Failed to connect to MySQL:", err)
	}
	if err = config.DB.Ping(); err != nil {
		log.Fatal("Failed to ping MySQL:", err)
	}
	defer config.DB.Close()
	createTableQuery := `
        CREATE TABLE IF NOT EXISTS employees (
            id INT AUTO_INCREMENT PRIMARY KEY,
            first_name VARCHAR(255) NOT NULL,
            last_name VARCHAR(255) NOT NULL,
            company_name VARCHAR(255),
            address VARCHAR(255),
            city VARCHAR(255),
            county VARCHAR(255),
            postal VARCHAR(255),
            phone VARCHAR(255),
            email VARCHAR(255),
            web VARCHAR(255)
        )
    `
	_, err = config.DB.Exec(createTableQuery)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}

	config.RedisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	config.Ctx = context.Background()
	if err := config.RedisClient.Ping(config.Ctx).Err(); err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	defer config.RedisClient.Close()

	router := gin.Default()

	router.POST("/upload", handlers.UploadFile)
	router.GET("/employees", handlers.GetEmployees)
	router.PUT("/employee/:id", handlers.UpdateEmployee)
	router.DELETE("/employee/:id", handlers.DeleteEmployee)

	router.Run(":8080")
}
