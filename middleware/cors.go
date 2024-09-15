package middleware

import "github.com/gin-gonic/gin"

// CORSMiddleware is a function that sets up Cross-Origin Resource Sharing (CORS) headers for a gin context.
//
// It sets the following headers:
// - Access-Control-Allow-Origin: "*"
// - Access-Control-Allow-Credentials: "true"
// - Access-Control-Allow-Headers: "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With"
// - Access-Control-Allow-Methods: "POST, OPTIONS, GET, PUT"
//
// It checks if the request method is "OPTIONS" and if so, aborts with a status code of 204.
//
// The function takes in a gin context (ctx *gin.Context) and does not return anything.
func CORSMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Set CORS
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(204)
			return
		}
		ctx.Next()
	}
}
