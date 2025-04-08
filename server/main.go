package main

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"rsc.io/quote"
)

type StudentConsultation struct {
	ID        uint
	Name      string
	StartDate time.Time
	EndDate   time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

type User struct {
	ID        uint
	FirstName string
	LastName  string
	Email     string
	Password  string
	Role      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func main() {

	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	var host, user, password, dbname, port string
	host = os.Getenv("DBHOST")
	user = os.Getenv("DBUSER")
	password = os.Getenv("DBPASSWORD")
	dbname = os.Getenv("DBNAME")
	port = os.Getenv("DBPORT")
	dsn := "host=" + host + " user=" + user + " password=" + password + " dbname=" + dbname + " port=" + port + " sslmode=disable"

	db, err_sql := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err_sql != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&StudentConsultation{}, &User{})
	//db.Create(&StudentConsultation{Name: "test", StartDate: time.Now(), EndDate: time.Now()})
	//db.Create(&User{FirstName: "test", LastName: "test", Email: "test", Password: "test", Role: "test"})
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, "hello world")
	})

	r.GET("/quote", func(c *gin.Context) {
		c.JSON(200, quote.Go())
	})
	r.OPTIONS("/register", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		c.JSON(200, nil)
	})
	r.POST("/register", func(c *gin.Context) {
		var newUser User
		c.Header("Access-Control-Allow-Origin", "http://localhost:5173")
		if err := c.BindJSON(&newUser); err != nil {
			c.JSON(500, "ERROR")
		}
		result := db.Create(&newUser)
		c.JSON(200, result)
	})
	r.OPTIONS("/login", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		c.JSON(200, nil)
	})
	r.POST("/login", func(c *gin.Context) {
		var user User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(500, "ERROR")
		}

		var foundUser User
		result := db.Where("email = ?", user.Email).First(&foundUser)
		if result.Error != nil {
			c.JSON(401, gin.H{"error": "Invalid credentials"})
			return
		}

		if foundUser.Password != user.Password {
			c.JSON(401, gin.H{"error": "Invalid credentials"})
			return
		}
		c.Header("Access-Control-Allow-Origin", "http://localhost:5173")
		c.JSON(200, gin.H{"user": foundUser})

		// usuario por email
		// comparar com o do request
		// se errado faz o L (401)
		// se der certo 200 usuario

	})

	r.Run(":3030")
}
