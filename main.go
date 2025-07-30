package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()
	r.GET("/ping", Ping)
	r.Run()
}

func Ping(context *gin.Context) {
	context.JSON(http.StatusTeapot, gin.H{"message": "Penis"})

}
  db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})