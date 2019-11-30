package controller

import (
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/liu599/FileServer/src/data"
	"github.com/liu599/FileServer/src/model"
	"github.com/liu599/FileServer/src/setting"
	"github.com/liu599/FileServer/src/utils"
	"gopkg.in/mgo.v2/bson"
	"mime"
	"net/http"
	"os"
	"path"
	"strings"
)

func Pong(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func GenFilePath(filename string, salt string) string {
	return strings.Join([]string{salt, filename}, "_")
}

func deleteFile(fileUrl string) error {
	return os.Remove(fileUrl)
}

func Upload(c *gin.Context) {
	var md5s []string
	var urls []string
	var relativePath = ""
	name := c.PostForm("name")
	email := c.PostForm("email")
	relativePath = c.PostForm("relativePath")
	rootPath := setting.FileRoot
    fileRootPath := rootPath+relativePath
	//fmt.Println(fileRootPath) // /a/b/c/d 或 \a\b\c\d
	//fmt.Println(fileRootPath[1:]) // /a/b/c/d 或 \a\b\c\d
	// 创建目录试试
	if err := os.MkdirAll(fileRootPath, 0777); err != nil {
		fmt.Println(err)
	}
	// Multipart form
	form, err := c.MultipartForm()
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}
	files := form.File["files"]

	for _, file := range files {
		salt := bson.NewObjectId().Hex()
		filePhyUrl := rootPath + relativePath + GenFilePath(file.Filename, salt)
		fmt.Println(filePhyUrl)
		if err := c.SaveUploadedFile(file, filePhyUrl); err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}
		if md5st, err := utils.HashFileMd5(filePhyUrl); err != nil {
			_ = deleteFile(filePhyUrl)
			c.String(http.StatusBadRequest, fmt.Sprintf("cannot generate file md5 err: %s", err.Error()))
			return
		} else {
			if err = model.CreateFile(data.NekohandFile{
				FileId:salt,
				HashId:md5st,
				FileName:file.Filename,
			}); err != nil {
				_ = deleteFile(filePhyUrl)
				c.String(http.StatusBadRequest, fmt.Sprintf("database error: cannot create file %s", err.Error()))
				return
			}
			md5s = append(md5s, md5st)
			urls = append(urls, salt + "_" + file.Filename)
		}

	}

	c.JSON(http.StatusOK, gin.H{
		"router": "/nekofile/",
		"file_md5_list": md5s,
		"psy_path_list": urls,
		"message": fmt.Sprintf("Uploaded successfully %d files with fields name=%s and email=%s.", len(files), name, email),
	})
}

const base64Table = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"

var coder = base64.NewEncoding(base64Table)

func Base64Encode(encodeByte []byte) []byte {
	return []byte(coder.EncodeToString(encodeByte))
}

func Base64Decode(decodeByte []byte) ([]byte, error) {
	return coder.DecodeString(string(decodeByte))

}

func File(c *gin.Context) {
	fileid := c.Param("fileid")
	fmt.Println(fileid)
	rootpath := setting.FileRoot
	relativePath, err := model.FetchFile(fileid)
	if err != nil {
		c.String(http.StatusNotFound, fmt.Sprintf("Cannot find file, wrong id. %s", err.Error()))
	}
	physicalPath := rootpath + relativePath
	//fmt.Println(physicalPath)
	fileInfo, err2 := os.Stat(physicalPath)
	//fmt.Println(fileInfo)
	if err2 != nil {
		if os.IsNotExist(err2) {
			c.String(http.StatusBadRequest, fmt.Sprintf("Cannot find file, probably physical deleted. %s", err2.Error()))
		}
	}
	//file, err := os.OpenFile(physicalPath, os.O_RDONLY, 0666)
	//if err != nil {
	//	c.String(http.StatusBadRequest, fmt.Sprintf("Cannot open file, this file is locked. %s", err.Error()))
	//	return
	//}
	mimetype := mime.TypeByExtension(path.Ext(physicalPath))

	//fmt.Println(mimetype)
	//data, err := ioutil.ReadAll(file)
	//fmt.Println(data)
	// https://stackoverflow.com/questions/31638447/how-to-server-a-file-from-a-handler-in-golang
	c.Header("content-length", string(fileInfo.Size()))
	c.Header("date", fileInfo.ModTime().String())
	c.Header("accept-ranges", "bytes")
	c.Header("content-type", mimetype)
	c.File(physicalPath)
}

