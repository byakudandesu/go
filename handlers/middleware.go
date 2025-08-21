package handlers

import "github.com/gin-gonic/gin"

func AdminOnly() gin.HandlerFunc {
	return func(request *gin.Context) {
		adminKey := request.GetHeader("admin-key")

		if adminKey != "byakubyaku" {
			request.JSON(403, gin.H{"error": "Admin access required"})
			request.Abort()
			return
		}
		request.Next()
	}
}
