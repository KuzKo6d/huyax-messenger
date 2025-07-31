package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `gorm:"created_at"`
	UpdatedAt time.Time      `gorm:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Username  string         `gorm:"unique;not null" json:"username"`
}

type Chat struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `gorm:"created_at"`
	UpdatedAt time.Time      `gorm:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	User1     uint           `gorm:"not null" json:"user_id1"`
	User2     uint           `gorm:"not null" json:"user_id2"`
}

type Message struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `gorm:"created_at"`
	UpdatedAt time.Time      `gorm:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Text      string         `gorm:"not null" json:"text"`
	Author    uint           `gorm:"not null" json:"author"`
	Chat      uint           `gorm:"not null" json:"chat"`
}

func main() {
	dsn := "host=localhost user=kuzko password=penis dbname=messenger port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&User{}, &Chat{}, &Message{})

	r := gin.Default()
	r.GET("/ping", ping)
	r.POST("/register", func(c *gin.Context) {
		register(c, db)
	})
	r.Run()
}

func ping(context *gin.Context) {
	context.JSON(http.StatusTeapot, gin.H{"message": "Miumiu"})
}

func register(context *gin.Context, db *gorm.DB) {
	// parse json
	// 1. data transfer object
	var user struct {
		Username string `json:"username" binding:"required"`
	}
	if err := context.ShouldBindJSON(&user); err != nil {
		log.Println(err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "Unable to take username"})
		return
	}

	// 2. Existing username check
	var existingUser User
	if err := db.Where("username = ?", user.Username).First(&existingUser).Error; err == nil {
		context.JSON(http.StatusConflict, gin.H{"message": "Username is already taken"})
		return
	}

	// 3. fill user info
	newUser := User{
		Username: user.Username,
	}

	// 4. create account
	if err := db.Create(&newUser).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to create user"})
		return
	}
	context.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

func sendMessage(context *gin.Context, db *gorm.DB, message Message) {

}

func chats(context *gin.Context, db *gorm.DB, user User) {

}

func messages(context *gin.Context, db *gorm.DB, chat Chat) {

}
