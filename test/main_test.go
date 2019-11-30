package test

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/liu599/FileServer/src/controller"
	"github.com/liu599/FileServer/src/data"
	"github.com/liu599/FileServer/src/middleware/func"
	"github.com/liu599/FileServer/src/setting"
	"gopkg.in/mgo.v2/bson"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

const filesTableCreationQuery = `
CREATE TABLE IF NOT EXISTS files
(
    fid      INT(11) NOT NULL AUTO_INCREMENT PRIMARY KEY,
    fileId   VARCHAR(50) UNIQUE NOT NULL,
    hashId   VARCHAR(50) UNIQUE NOT NULL,
	fileName   VARCHAR(100) UNIQUE NOT NULL,
	createdAt  INT(64)  NOT NULL,
	modifiedAt INT(64) NOT NULL
) character set = utf8`

var db *sqlx.DB

var NewDate = time.Now().Unix()


func TestMain(m *testing.M) {
	var Apps = make(map[string]data.Database)
	Apps["nekohand"] = data.Database{
		Driver: "mysql",
		MaxIdle: 2,
		MaxOpen: 15,
		Name: "nekohand",
		Source: setting.Source,
	}
	_func.AssignAppDataBaseList(Apps)
	_func.AssignDatabaseFromList([]string{"nekohand"})
	db, err := _func.MySqlGetDB("nekohand")
	if err != nil {
		fmt.Println("Error Database Connection")
		return
	}
	insertOneData(db)
	code := m.Run()
	os.Exit(code)
}

func insertOneData(db *sqlx.DB) {
	//fmt.Println(id, "Inserted")
	statement := fmt.Sprintf("INSERT INTO files (fileId, hashId, fileName, createAt, modifiedAt) VALUES('%s', '%s', '%s', '%d', '%d')", bson.NewObjectId().Hex(), "123ajkdf3afdsaf", "afd.jpg", NewDate, NewDate)
	_, err := db.Exec(statement)

	if err != nil {
		fmt.Println("Database error")
	}
}


// 发送请求
func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.GET("/filelist", controller.FileList)
	engine.GET("/ping", controller.Pong)
	engine.POST("/upload", controller.Upload)
	engine.ServeHTTP(rr, req)

	return rr
}

func TestBasicServer(t *testing.T) {
	req, _ := http.NewRequest("GET", "/ping", nil)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	response := executeRequest(req)
	fmt.Println(response.Body)
}


func TestUploadFile(t *testing.T) {
	filename := "D:/Project/PictureServer/test/QQ图片20180812010011.jpg"
	filenamex := "QQ图片20180812010011.jpg"
	/*
		application/x-www-form-urlencoded	在发送前编码所有字符（默认）
		multipart/form-data	 不对字符编码。在使用包含文件上传控件的表单时，必须使用该值。
		text/plain	空格转换为 "+" 加号，但不对特殊字符编码。
	*/
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
    _ = bodyWriter.WriteField("name", "tokei")
    _ = bodyWriter.WriteField("email", "460512944@qq.com")
    _ = bodyWriter.WriteField("relativePath", "blue/red/")
	//关键的一步操作
	fileWriter, err := bodyWriter.CreateFormFile("files", filenamex)
	if err != nil {
		fmt.Println("error writing to buffer")
	}

	//打开文件句柄操作
	fh, err := os.Open(filename)
	if err != nil {
		fmt.Println("error opening file")
	}
	defer fh.Close()

	//iocopy
	_, err = io.Copy(fileWriter, fh)

	contentType := bodyWriter.FormDataContentType()
	_ = bodyWriter.Close()
	req, _ := http.NewRequest("POST", "/upload", bodyBuf)
	//req.Header.Add("Content-Type", "multipart/form-data")
	req.Header.Add("Content-Type", contentType)
	response := executeRequest(req)
	fmt.Println(response.Body)
}

func TestFileList(t *testing.T) {
	//req, _ := http.NewRequest("GET", "/filelist", nil)
	//req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	//response := executeRequest(req)
	//fmt.Println(response.Body)
	fmt.Println(1)
}

func TestFileCatch(t *testing.T) {
	req, _ := http.NewRequest("GET", "/nekofile/5db884d458adfe413b4be2a2/", nil)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	response := executeRequest(req)
	fmt.Println(response.Body)
}

// 检查Response
func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}




