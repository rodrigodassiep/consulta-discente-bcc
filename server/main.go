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

	r.Run(":3030")
}
