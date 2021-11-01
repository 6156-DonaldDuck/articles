package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	message string
}

// New security middleware
func Security() gin.HandlerFunc {
	return func (c *gin.Context) {
		whiteListMap := map[string] []string {
			"/swagger/*any": {"GET"},
			"/api/v1/articles": {"GET", "POST"},
			"/api/v1/articles/:articleId": {"GET","PUT"},
		}

		// Figure out whether current request is in white list.
		allowed := false
		allowedMethods, ok := whiteListMap[c.FullPath()]
		if ok {
			for _, method := range allowedMethods {
				if method == c.Request.Method {
					allowed = true
					break
				}
			}
		}

		if allowed {
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": false, "message":"Unauthorized"})
		}
	}
}