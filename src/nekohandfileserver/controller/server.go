package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"fmt"
)

func Pong(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func Upload(c *gin.Context) {
	name := c.PostForm("name")
	email := c.PostForm("email")

	// Multipart form
	form, err := c.MultipartForm()
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}
	files := form.File["file"]

	for _, file := range files {
		if err := c.SaveUploadedFile(file, "D:/Project/MyBlogCMS71/src/nekohandfileserver/files/" + file.Filename); err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Uploaded successfully %d files with fields name=%s and email=%s.", len(files), name, email),
	})
}