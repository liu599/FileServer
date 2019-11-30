package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_func "github.com/liu599/FileServer/src/middleware/func"
	"github.com/liu599/FileServer/src/model"
	"net/http"
)

func FileList(c *gin.Context) {
	err, filelist := model.FetchFileList("ALL")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("Error filelist %s", err.Error()))
		return
	}
	list := make(map[string]interface{})
	list["data"] = filelist
	_func.Respond(c, http.StatusOK, list)
}

func FileListByType(c *gin.Context) {
	filetype := c.PostForm("filetype")
	err, filelist := model.FetchFileList(filetype)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("Error filelist %s", err.Error()))
		return
	}
	list := make(map[string]interface{})
	list["data"] = filelist
	_func.Respond(c, http.StatusOK, list)
}

