package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

func init() {
	var err error
	db, err = gorm.Open("mysql", "<DB-USER>:<DB-PASSWORD>@/<DB-NAME>?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("Unable to connect to database")
	}

	db.AutoMigrate(&hotel{}, &room{})
}

func main() {

	router := gin.Default()

	v1 := router.Group("/api")
	{
		v1.POST("/hotel", createHotel)
		v1.POST("/room", createRoom)
		v1.GET("/available-room", availableRoom)
	}
	router.Run()

}

type (
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
