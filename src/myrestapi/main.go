package main

import (
	"myrestapi/controller"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
)


func main() {
	r := gin.Default()
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
	r.StaticFS("/files", http.Dir("D:/Project/MyBlogCMS71/src/myrestapi/files"))
	r.GET("/ping", controller.Pong)
	r.POST("/upload", controller.Upload)
	r.Run(":4545") // 默认为8080端口
}