package test

import (
	"testing"
	"os"
	"net/http"
	"net/http/httptest"
	"github.com/gin-gonic/gin"
	"fmt"
	"nekohandfileserver/controller"
	"bytes"
	"mime/multipart"
	"io"
	"nekohandfileserver/middleware/data"
	"nekohandfileserver/middleware/func"
	"github.com/jmoiron/sqlx"
	"gopkg.in/mgo.v2/bson"
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

var Pass_gen = "000000"

var Neko_token = "4a2c4b"

var NewDate = time.Now().Unix()


func TestMain(m *testing.M) {
	var Apps = make(map[string]data.Database)
	Apps["nekohand"] = GetDataBase()
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
	filename := "D:/Project/MyBlogCMS71/src/myrestapi/test/DqCKyZUVYAEDHxN.jpg"
	/*
		application/x-www-form-urlencoded	在发送前编码所有字符（默认）
		multipart/form-data	 不对字符编码。在使用包含文件上传控件的表单时，必须使用该值。
		text/plain	空格转换为 "+" 加号，但不对特殊字符编码。
	*/
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
    _ = bodyWriter.WriteField("name", "tokei")
    _ = bodyWriter.WriteField("email", "460512944@qq.com")
	//关键的一步操作
	fileWriter, err := bodyWriter.CreateFormFile("files", filename)
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
	bodyWriter.Close()
	req, _ := http.NewRequest("POST", "/upload", bodyBuf)
	//req.Header.Add("Content-Type", "multipart/form-data")
	req.Header.Add("Content-Type", contentType)
	response := executeRequest(req)
	fmt.Println(response.Body)
}

// 检查Response
func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}




