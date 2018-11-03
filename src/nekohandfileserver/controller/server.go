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
	"encoding/base64"
	"mime"
	"path"
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

func GenFilePath(filename string, salt string) string {
	rootpath := os.Getenv("SERVER_FILE_PATH")
	return rootpath + strings.Join([]string{salt, filename}, "_")
}

func deleteFile(fileUrl string) error {
	return os.Remove(fileUrl)
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
		filePhyUrl := GenFilePath(file.Filename, salt)
		if err := c.SaveUploadedFile(file, filePhyUrl); err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}
		if md5st, err := hashFileMd5(filePhyUrl); err != nil {
			deleteFile(filePhyUrl)
			c.String(http.StatusBadRequest, fmt.Sprintf("cannot generate file md5 err: %s", err.Error()))
			return
		} else {
			if err = model.CreateFile(data.NekohandFile{
				FileId:salt,
				HashId:md5st,
				FileName:file.Filename,
			}); err != nil {
				deleteFile(filePhyUrl)
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

const base64Table = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"

var coder = base64.NewEncoding(base64Table)

func Base64Encode(encode_byte []byte) []byte {
	return []byte(coder.EncodeToString(encode_byte))
}

func Base64Decode(decode_byte []byte) ([]byte, error) {
	return coder.DecodeString(string(decode_byte))

}

func File(c *gin.Context) {
	fileid := c.Param("fileid")
	rootpath := os.Getenv("SERVER_FILE_PATH")
	relativePath := model.FetchFile(fileid)
	physicalPath := rootpath + relativePath
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