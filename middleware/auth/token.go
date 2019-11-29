package auth

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gin-gonic/gin"
	"nekohandfileserver/middleware/data"
	"nekohandfileserver/middleware/func"
)

func TokenAuthMiddleware() gin.HandlerFunc {
	salt := os.Getenv("NEKO_TOKEN")
	return func(c *gin.Context) {
		requestPath := c.Request.URL.Path
		if c.Request.Method == http.MethodPost {
			if !strings.Contains(requestPath, "auth") {
				c.Next()
			} else {
				requestUser, ok := c.Request.Header["User"]
				if !ok {
					requestUser = []string{"_"}
				}
				tokenRequest, err := request.ParseFromRequest(c.Request, request.AuthorizationHeaderExtractor,
					func(token *jwt.Token) (interface{}, error) {
						return []byte(requestUser[0] + salt), nil
					})
				if err != nil {
					_func.RespondError(c, http.StatusBadRequest,
						data.Error{
							Code: "401",
							Message: fmt.Sprintf("%v", err),
						})
					return
				}
				if tokenRequest.Valid {
					c.Next()
				} else {
					_func.RespondError(c, http.StatusUnauthorized,
						data.Error{
							Message: "此Token已过期或者未授权",
						})
					return
				}
			}
		} else {
			c.Next()
		}
	}
}

func GenerateToken(username string, environmentVariable string) (string, error) {

	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(1)).Unix()
	claims["iat"] = time.Now().Unix()
	token.Claims = claims
	//fmt.Printf("\n%v\n", username + environmentVariable)
	tokenString, err := token.SignedString([]byte(username + environmentVariable))

	//fmt.Printf("\n%v\n", tokenString)

	if err != nil {
		return "Cannot convert secret string", err
	}

	return tokenString, nil
}