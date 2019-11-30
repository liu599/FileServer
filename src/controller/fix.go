package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/liu599/FileServer/src/model"
	"net/http"
)

// 将服务器根目录下某一文件夹写入数据库
func Fix(c *gin.Context) {
	userpath := c.PostForm("userpath")
	filetype := c.PostForm("filetype")
	err := model.FixPath(userpath, filetype)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": fmt.Sprintf("%v", err) ,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"success": "true",
		})
	}
}

