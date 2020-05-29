package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"golang.org/x/crypto/bcrypt"
)

var db *gorm.DB

func init() {
	var err error
	db, err = gorm.Open("mysql", "root:karachi123@/hotel-management?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("Unable to connect to database")
	}

	db.AutoMigrate(&hotel{}, &room{}, &user{})
}

func main() {

	router := gin.Default()

	v1 := router.Group("/api")
	{	
		v1.POST("/signup", signUp)
		v1.POST("/signin", signIn)
		v1.POST("/hotel", createHotel)
		v1.POST("/room", createRoom)
		v1.GET("/available-room", availableRoom)
	}
	router.Run()

}

type (
	user struct {
		gorm.Model
		Email string `json:"email"`
		Password   string    `gorm:"size:100;not null;" json:"password"`
		// Password []byte `json:"password"`
	}

	hotel struct {
		gorm.Model
		Name string `json:"name"`
		Area string `json:"area"`
	}

	room struct {
		gorm.Model
		Number    string
		Size      string
		hotel     hotel `gorm:"foreignkey:HotelId"`
		HotelId   int
		StartDate *time.Time
		EndDate   *time.Time
		BookedBy  int
	}
)

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func signUp(c *gin.Context) {
	email := c.PostForm("email")
	hashedPassword, err := Hash(c.PostForm("password"))
	passwordAsString := string(hashedPassword)
	if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Bycrypt error",})
	}	
    var admin []user 
	db.Where("email = ?", email).Find(&admin)
	// fmt.Printf("%+v\n", admin)
	// if len(*admin) {
	// 	c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Email already exist",})
	// }
	user := user{Email: email, Password: passwordAsString}
	db.Save(&user)
	c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "User created successfully", "email": user.Email})
}

func signIn(c *gin.Context) {
	email := c.PostForm("email")
	enteredPassword := c.PostForm("password")
	
    admin :=user{}
	db.Where("email = ?", email).Find(&admin)

	err := VerifyPassword(admin.Password, enteredPassword)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "message": "Wrong credentials"})
		return
	}
	fmt.Printf("%+v\n", admin)
	c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "Logged In"})
	
}



func createHotel(c *gin.Context) {
	hotel := hotel{Name: c.PostForm("name"), Area: c.PostForm("area")}
	db.Save(&hotel)
	c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "Hotel created successfully", "resourceId": hotel.ID})
}

func createRoom(c *gin.Context) {
	hotelId, _ := strconv.Atoi(c.PostForm("hotelId"))
	room := room{Number: c.PostForm("number"), Size: c.PostForm("size"), HotelId: hotelId}
	db.Save(&room)
	c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "Room created", "resourceId": room.ID})
}

func availableRoom(c *gin.Context) {
	var availableRooms []room
	db.Where("booked_by = ?", 0).Find(&availableRooms)

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": availableRooms})
}
