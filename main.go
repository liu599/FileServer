package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/liu599/FileServer/src/controller"
	"github.com/liu599/FileServer/src/data"
	"github.com/liu599/FileServer/src/middleware/func"
	"github.com/liu599/FileServer/src/setting"
	"log"
	"net/http"
)


func main() {

	database := data.Database{
		Driver: "mysql",
		MaxIdle: setting.MaxIdle,
		MaxOpen: setting.MaxOpen,
		Name: "shana",
		Source: setting.Source,
	}

	var Apps = make(map[string]data.Database)

	Apps["nekohand"] = database

	_func.AssignAppDataBaseList(Apps)

	_func.AssignDatabaseFromList([]string{"nekohand"})
	//
	gin.SetMode(setting.RunMode)
	r := gin.Default()
	//
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	sysFilePath := setting.FileRoot
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
	// Router
	r.GET("/ping", controller.Pong)
	r.GET("/filelist", controller.FileList)
	r.POST("/filetype", controller.FileListByType)
	r.GET("/nekofile/:fileid/*size", controller.File)

	r.POST("/upload", controller.Upload)
	r.POST("/fix", controller.Fix)


	er := r.Run(setting.HTTPPort) // 默认为8080端口
	if er != nil {
		log.Fatalf("Server cannot start: ': %v", er)
	}
}