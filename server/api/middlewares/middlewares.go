package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/victorsteven/server/api/authen"
	
)


func TokenAuthMiddleware() gin.HandlerFunc {
	// Tạo một map gắp lỗi xuất hiên
	errList := make(map[string]string)
	// Trả về một bao đóng
	return func(c *gin.Context) {
		err := authen.TokenValid(c.Request)
		if err != nil {
			errList["unauthorized"] = "Unauthorized"
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": http.StatusUnauthorized,
				"error":  errList,
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// Thằng này giúp tương tác với react Frontend
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
