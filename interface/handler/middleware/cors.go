package middleware

import (
	"net/http"

	"xxx-server/infrastructure/config"

	"github.com/gin-gonic/gin"
)

// cors middleware
func CorsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// origin := c.GetHeader("Origin")
		// if origin == "" || !strings.Contains(origin, "localhost") || !strings.Contains(origin, "127.0.0.1") {
		// 	origin = config.C.Server.DefaultDomain
		// }
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(204)
		}

		if !config.C.Server.DevMode {
			return
		}
		c.Header("Access-Control-Allow-Origin", "*")
		// c.Header("Access-Control-Allow-Credentials", "true")
		// c.Header("Access-Control-Allow-Headers",
		// 	"Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		// c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		c.Header("Vary", "Origin,Access-Control-Request-Method,Access-Control-Request-Headers")

	}
}
