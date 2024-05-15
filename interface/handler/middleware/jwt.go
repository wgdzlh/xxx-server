package middleware

import (
	"errors"
	"strings"

	"xxx-server/application/utils"
	"xxx-server/infrastructure/config"
	"xxx-server/interface/resp"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

const (
	USER_ID_KEY = "jti"
	USER_ID_CTX = "user_id"
)

var (
	ErrInvalidUserId = errors.New("invalid user_id in JWT")
)

func AuthJWTChecker() gin.HandlerFunc {
	return func(c *gin.Context) {
		if config.C.Server.DevMode || config.C.Server.JwtSecret == "" || strings.HasSuffix(c.GetHeader("Referer"), "/swagger/index.html") {
			return
		}
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			resp.AbortJSONWithMsg(c, resp.ErrUnauthorized, "Authorization header missing")
			return
		}

		parts := strings.SplitN(authHeader, "Bearer", 2)
		if len(parts) < 2 {
			resp.AbortJSONWithMsg(c, resp.ErrUnauthorized, "wrong token")
			return
		}

		tokenRaw := strings.TrimSpace(parts[1])
		token, err := jwt.Parse(tokenRaw, func(token *jwt.Token) (any, error) {
			return []byte(config.C.Server.JwtSecret), nil
		})
		if err != nil {
			resp.AbortJSONWithMsg(c, resp.ErrUnauthorized, err.Error())
			return
		}

		var userId int
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			switch jwtUserId := claims[USER_ID_KEY].(type) {
			case string:
				userId = utils.StrToInt(jwtUserId)
			case float64:
				userId = int(jwtUserId)
			case int:
				userId = jwtUserId
			default:
				resp.AbortJSONWithMsg(c, resp.ErrUnauthorized, "user_id type unknown in JWT")
				return
			}
		} else {
			resp.AbortJSONWithMsg(c, resp.ErrUnauthorized, "token claims invalid")
			return
		}
		c.Set(USER_ID_CTX, userId)
		c.Next()
	}
}
