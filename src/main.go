package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

type Weather struct {
	ID          int       `gorm:"column:id;primary_key" json:"id"`
	RecordedAt  time.Time `gorm:"column:recorded_at" json:"datetime"`
	Temperature float32   `gorm:"column:temperature" json:"temperature"`
	Humidity    float32   `gorm:"column:humidity"  json:"humidity"`
}

type SensorUploadData struct {
	Weathers []Weather `json:"data"`
}

func InitDb() *gorm.DB {
	// Openning file
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "data.db"
	}
	db, err := gorm.Open("sqlite3", dbPath)
	// Display SQL queries
	// db.LogMode(true)

	// Error
	if err != nil {
		panic(err)
	}

	// Migrate the schema
	db.AutoMigrate(&Weather{})

	return db
}

func main() {
	db := InitDb()
	defer db.Close()

	r := gin.Default()

	v1 := r.Group("api/v1")
	{
		v1.POST("/echo", Echo)
		v1.POST("/sensor/upload", UploadSensorData)
		v1.GET("/weather/query", QueryWeather)
		v1.GET("/weather/latest", LatestWeather)
	}

	r.Run(":8080")
}

func Echo(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	c.Writer.Write(body)
}

func UploadSensorData(c *gin.Context) {
	db := InitDb()
	defer db.Close()

	var uploadData SensorUploadData
	c.Bind(&uploadData)

	for _, weather := range uploadData.Weathers {
		db.Create(&weather)
	}
	c.JSON(201, gin.H{"status": "ok"})
}

func QueryWeather(c *gin.Context) {
	db := InitDb()
	defer db.Close()

	param := c.Request.URL.Query()
	fmt.Println(param)
	startDatetime := param.Get("start_datetime")
	endDatetime := param.Get("end_datetime")

	if startDatetime == "" || endDatetime == "" {
		c.JSON(403, gin.H{"status": "error", "message": "must provide start_datetime and end_datetime"})
	}

	weathers := []Weather{}
	db.Find(&weathers, "recorded_at >= ? AND recorded_at <= ?", startDatetime, endDatetime).Order("recorded_at asc")
	c.JSON(200, gin.H{"status": "ok", "data": weathers})
}

func LatestWeather(c *gin.Context) {
	db := InitDb()
	defer db.Close()

	weather := Weather{}
	db.Last(&weather).Order("recorded_at asc")

	if weather.ID == 0 {
		c.JSON(200, gin.H{"status": "no_data"})
	} else {
		c.JSON(200, gin.H{"status": "ok", "data": weather})
	}
}
