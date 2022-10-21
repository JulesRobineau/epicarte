package middleware

import (
	"errors"
	"gin-template/config"
	token2 "gin-template/pkg/model"
	"gin-template/pkg/model/enum"
	error2 "gin-template/utils/error"
	jwt2 "gin-template/utils/jwt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

type Claims struct {
	jwt.RegisteredClaims
	UserId int64     `json:"user_id"`
	Role   enum.Role `json:"role"`
}

type JwtMiddleware struct {
	Conf config.JwtConfig
	Jwt  jwt2.JwtManager
}

func NewJwtMiddleware(conf config.JwtConfig) JwtMiddleware {
	return JwtMiddleware{
		Conf: conf,
		Jwt:  jwt2.NewJwtManager(conf.Secret, conf.Expiration),
	}
}

func (j *JwtMiddleware) MiddlewareFunc(accessRoles map[string][]enum.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		db := c.MustGet("DB").(*gorm.DB)
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, error2.UnauthorizedError("token is required"))
			return
		}
		if !strings.HasPrefix(token, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, error2.UnauthorizedError("token must be prefixed with Bearer"))
			return
		}

		token = strings.TrimPrefix(token, "Bearer ")
		c.Set("token", token)

		rClaims, err := jwt2.ParseToken(token, j.Conf.Secret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, error2.UnauthorizedError(err.Error()))
			return
		}
		c.Set("claims", rClaims)

		tokenModel := token2.TokenModel{Tx: db}
		err = tokenModel.FindToken(&token2.Token{TokenID: rClaims.ID}).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, error2.UnauthorizedError("token is expired"))
			return
		}

		// Get last segment of handler path
		// e.g. gin-template/pkg/service/v1.Login-fm => Login
		// e.g. gin-template/pkg/common/AuthInterface.Login => Login
		handler := c.HandlerName()
		if strings.HasSuffix(handler, "-fm") {
			handler = handler[:len(handler)-3]
		}
		handler = handler[strings.LastIndex(handler, ".")+1:]

		// check if user has access to this route
		access, ok := accessRoles[handler]
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, error2.ForbiddenError("access denied"))
			return
		}

		// check if role is allowed to access this route
		for _, role := range access {
			if rClaims.Role.HasPermission(role) {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, error2.ForbiddenError("access denied"))
	}
}
