package main

import (
	"database/sql"
	"fmt"
	controller "root/Controller"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
)

var db *sql.DB

func main() {
	router := gin.Default()
	gin.SetMode(gin.DebugMode)
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("error loading .env file")
	}

	router.POST("/user/create", controller.StoreData)

	if err := router.Run("0.0.0.0:8000"); err != nil {
		panic(fmt.Sprintf("failed to start the server: %v", err))
	}
}
