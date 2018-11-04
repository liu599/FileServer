package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"nekohandfileserver/controller"
	"os"
	"strconv"
	"nekohandfileserver/middleware/data"
	"nekohandfileserver/middleware/func"
	"net/http"
)


func main() {

	//Configure()
	//os.Setenv("SERVER_FILE_PATH",  "D:/Pictures/")
	maxIdle, _ := strconv.Atoi(os.Getenv("SERVER_DB_MAX_IDLE"))
	maxOpen, _ := strconv.Atoi(os.Getenv("SERVER_DB_MAX_OPEN"))
	source := os.Getenv("SERVER_DB_URL")

	database := data.Database{
		Driver: "mysql",
		MaxIdle: maxIdle,
		MaxOpen: maxOpen,
		Name: "nekohand",
		Source: source,
	}

	var Apps = make(map[string]data.Database)

	Apps["nekohand"] = database

	_func.AssignAppDataBaseList(Apps)

	_func.AssignDatabaseFromList([]string{"nekohand"})
	//
	r := gin.New()
	//
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	sysFilePath := os.Getenv("SERVER_FILE_PATH")
	r.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "OPTIONS", "PUT"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "User-Agent", "Referrer", "Host", "Token", "User", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "X-Real-Ip"},
		AllowCredentials: true,
		AllowAllOrigins:  false,
		AllowOriginFunc:  func(origin string) bool { return true },
		MaxAge:           86400,
	}))
	r.MaxMultipartMemory = 8 << 20 // 8 MiB
	r.Static("/files", sysFilePath)
	r.StaticFS("/nhfiles", http.Dir(sysFilePath))
	r.GET("/ping", controller.Pong)
	r.GET("/filelist", controller.FileList)
	r.GET("/nekofile/:fileid/*size", controller.File)
	r.POST("/upload", controller.Upload)
	r.Run(":17699") // 默认为8080端口
}