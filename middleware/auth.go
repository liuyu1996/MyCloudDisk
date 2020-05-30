package middleware

import (
	"MyCloudDisk/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func CheckAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
			username := c.Request.FormValue("username")
			token := c.Request.FormValue("token")
			claims, err := utils.ParseToken(token)
			if len(username) < 3 || err != nil || claims == nil {
				c.Abort()
				resp := utils.NewRespMsg(-1, "token无效", nil)
				c.JSON(http.StatusOK, resp)
				return
			}else if time.Now().Unix() > claims.ExpiresAt {
				c.Abort()
				resp := utils.NewRespMsg(-2, "token已过期", nil)
				c.JSON(http.StatusOK, resp)
				return
			}
			c.Next()
		}
}
