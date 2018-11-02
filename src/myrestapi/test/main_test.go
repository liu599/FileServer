package test

import (
	"testing"
	"os"
	"net/http"
	"net/http/httptest"
	"github.com/gin-gonic/gin"
	"fmt"
	"myrestapi/controller"
	"bytes"
	"mime/multipart"
	"io"
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
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




