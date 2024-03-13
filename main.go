package main

import (
	controller "root/Controller"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.POST("/user/uploadfile", controller.StoreFile)
}
