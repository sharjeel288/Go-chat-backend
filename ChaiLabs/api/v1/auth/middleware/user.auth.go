package middleware

import (
	userAuthInterface "ChaiLabs/ChaiLabs/api/v1/auth/interface"
	userTypes "ChaiLabs/ChaiLabs/api/v1/user/types"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	socketio "github.com/googollee/go-socket.io"
)

func AuthMiddleware(userAuthService userAuthInterface.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if the route is public
		if c.GetBool("isPublic") {
			c.Next()
			return
		}

		// Extract token from HTTP header
		rawToken := c.GetHeader("Authorization")
		if rawToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"unable to get token": "Unauthorized"})
			return
		}

		var token userTypes.Token
		err := json.Unmarshal([]byte(rawToken), &token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error unmarshalling user token": err.Error()})
			return
		}

		// Validate user token
		user, err := userAuthService.ValidateUser(&token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Validating user token": err.Error()})
			return
		}

		// Attach the user object to the request
		c.Set("user", user)
		c.Next()
	}
}

func WebSocketAuthMiddleware(userAuthService userAuthInterface.AuthService) func(socketio.Conn, http.Header) error {
	return func(so socketio.Conn, header http.Header) error {
		rawToken := header.Get("Authorization")
		if rawToken == "" {
			return errors.New("WsAuth:token not found")
		}

		var token userTypes.Token
		err := json.Unmarshal([]byte(rawToken), &token)
		if err != nil {
			return err
		}

		user, err := userAuthService.ValidateUser(&token, true)
		if err != nil {
			return err
		}

		// Attach the user object to the socket connection
		so.SetContext(user)

		return nil
	}
}
