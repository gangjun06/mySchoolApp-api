package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/gangjun06/mySchoolApp-api/middlewares"
	"github.com/gangjun06/mySchoolApp-api/models"
	v1 "github.com/gangjun06/mySchoolApp-api/routes/v1"
	"github.com/gangjun06/mySchoolApp-api/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	etcInit()
	applyConfig()
	initDB()
	startServer()
}

func etcInit() {
	rand.Seed(time.Now().Unix())
}

func applyConfig() {
	rawConfig, err := ioutil.ReadFile("config.toml")
	if err != nil {
		log.Fatalln("Failed to load config.")
	}

	var config models.Config
	if _, err := toml.Decode(string(rawConfig), &config); err != nil {
		log.Fatalln("Failed to parsing config.")
	}
	utils.SetConfig(&config)
}

func startServer() {
	config := utils.GetConfig().Server
	if config.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)

	}
	r := gin.Default()
	r.Use(middlewares.Cors())
	version1 := r.Group("/v1")
	v1.InitRoutes(version1)
	r.Run(":" + strconv.Itoa(config.Port))
}

func initDB() {
	log.Println("Initializing Database...")
	dbConfig := utils.GetConfig().DB

	var db *gorm.DB
	var err error

	if dbConfig.UseSqlite {
		db, err = gorm.Open(sqlite.Open("./data.db"), &gorm.Config{})
	} else {
		c := dbConfig.Server
		connectionInfo := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", c.Username, c.Password, c.Hostname, c.Port, c.DBName)
		db, err = gorm.Open(mysql.Open(connectionInfo), &gorm.Config{})
	}

	if err != nil {
		log.Fatalln("Failed to open database.")
	}
	utils.SetDB(db)
	log.Print("Successfully Connected To Database")

	if err := db.AutoMigrate(); err != nil {
		log.Fatalln("Failed to perform AutoMigrate.")
	}
	log.Print("Successfully performed AutoMigrate")
}
