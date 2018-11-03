package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"fmt"
	"os"
	"crypto/md5"
	"io"
	"encoding/hex"
	"gopkg.in/mgo.v2/bson"
	"strings"
	"nekohandfileserver/model"
	"nekohandfileserver/middleware/data"
)

func Pong(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func hashFileMd5(filePath string) (string, error) {
	var returnMD5String string
	file, err := os.Open(filePath)
	if err != nil {
		return returnMD5String, err
	}
	defer file.Close()
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return returnMD5String, err
	}
	hashInBytes := hash.Sum(nil)[:16]
	returnMD5String = hex.EncodeToString(hashInBytes)
	return returnMD5String, nil
}

func genFilePath(filename string, salt string) string {
	rootpath := os.Getenv("SERVER_FILE_PATH")
	return rootpath + strings.Join([]string{salt, filename}, "_")
}

func Upload(c *gin.Context) {
	var md5s []string
	var urls []string
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
		salt := bson.NewObjectId().Hex()
		filePhyUrl := genFilePath(file.Filename, salt)
		if err := c.SaveUploadedFile(file, filePhyUrl); err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}
		if md5st, err := hashFileMd5(filePhyUrl); err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("cannot generate file md5 err: %s", err.Error()))
			return
		} else {
			if err = model.CreateFile(data.NekohandFile{
				FileId:salt,
				HashId:md5st,
				FileName:file.Filename,
			}); err != nil {
				c.String(http.StatusBadRequest, fmt.Sprintf("database error: cannot create file %s", err.Error()))
				return
			}
			md5s = append(md5s, md5st)
			urls = append(urls, filePhyUrl)
		}

	}

	c.JSON(http.StatusOK, gin.H{
		"file_md5_list": md5s,
		"psy_path_list": urls,
		"message": fmt.Sprintf("Uploaded successfully %d files with fields name=%s and email=%s.", len(files), name, email),
	})
}
