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
	User1     uint           `gorm:"not null" json:"user1-id"`
	User2     uint           `gorm:"not null" json:"user2-id"`
}

type Message struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `gorm:"created_at"`
	UpdatedAt time.Time      `gorm:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Text      string         `gorm:"not null" json:"text"`
	Author    uint           `gorm:"not null" json:"author-id"`
	Chat      uint           `gorm:"not null" json:"chat-id"`
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
	r.POST("/create-chat", func(c *gin.Context) {
		createChat(c, db)
	})
	r.POST("/send-message", func(c *gin.Context) {
		sendMessage(c, db)
	})
	r.GET("/chats")
	r.GET("/messages")
	r.Run()
}

func ping(context *gin.Context) {
	context.JSON(http.StatusTeapot, gin.H{"message": "Miumiu"})
}

func register(context *gin.Context, db *gorm.DB) {
	// parse json
	var user struct {
		Username string `json:"username" binding:"required"`
	}
	if err := context.ShouldBindJSON(&user); err != nil {
		log.Println(err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "Unable to take username"})
		return
	}

	// existing username check
	var existingUser User
	if err := db.Where("username = ?", user.Username).First(&existingUser).Error; err == nil {
		context.JSON(http.StatusConflict, gin.H{"message": "Username is already taken"})
		return
	}

	// create account
	newUser := User{
		Username: user.Username,
	}

	if err := db.Create(&newUser).Error; err != nil {
		log.Println(err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to create user"})
		return
	}
	context.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

func createChat(context *gin.Context, db *gorm.DB) {
	// parse json
	var chat struct {
		User1 uint `json:"user1-id" binding:"required"`
		User2 uint `json:"user2-id" binding:"required"`
	}
	if err := context.ShouldBindJSON(&chat); err != nil {
		log.Println(err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "Unable to take chat"})
		return
	}

	// existing user1 check
	var existingUser1 User
	if err := db.Where("id = ?", chat.User1).First(&existingUser1).Error; err != nil {
		log.Println(err)
		context.JSON(http.StatusNotFound, gin.H{"message": "User1 not exist"})
		return
	}

	// existing user2 check
	var existingUser2 User
	if err := db.Where("id = ?", chat.User2).First(&existingUser2).Error; err != nil {
		log.Println(err)
		context.JSON(http.StatusNotFound, gin.H{"message": "User2 not exist"})
		return
	}

	// existing chat check
	var existingChat Chat
	err := db.Where("(user1 = ? AND user2 = ?) OR (user1 = ? AND user2 = ?)",
		chat.User1, chat.User2, chat.User2, chat.User1).First(&existingChat).Error
	if err == nil {
		context.JSON(http.StatusConflict, gin.H{"message": "Chat already exists"})
	}

	// create chat
	var newChat = Chat{
		User1: chat.User1,
		User2: chat.User2,
	}
	if err := db.Create(&newChat).Error; err != nil {
		log.Println(err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to create chat"})
		return
	}
	context.JSON(http.StatusCreated, gin.H{"message": "Chat created successfully"})
}

func sendMessage(context *gin.Context, db *gorm.DB) {
	// parse json
	var message struct {
		Text   string `json:"text" binding:"required"`
		Author uint   `json:"author-id" binding:"required"`
		Chat   uint   `json:"chat-id" binding:"required"`
	}

	if err := context.ShouldBindJSON(&message); err != nil {
		log.Println(err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "Unable to take message"})
		return
	}

	// create message
	newMessage := Message{
		Text:   message.Text,
		Author: message.Author,
		Chat:   message.Chat,
	}

	if err := db.Create(&newMessage).Error; err != nil {
		log.Println(err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to send message"})
		return
	}
	context.JSON(http.StatusCreated, gin.H{"message": "Message sent successfully"})
}

func chats(context *gin.Context, db *gorm.DB) {

}

func messages(context *gin.Context, db *gorm.DB) {

}
