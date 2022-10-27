package v1

import (
	"gin-template/config"
	"gin-template/pkg/common/auth"
	"gin-template/pkg/dto"
	"gin-template/pkg/middleware"
	"gin-template/pkg/model"
	"gin-template/pkg/model/enum"
	error2 "gin-template/utils/error"
	jwt2 "gin-template/utils/jwt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthService struct {
	jwt config.JwtConfig
}

// Login login a user
// @Summary Login
// @Description Logs in a user
// @Tags auth
// @Accept json
// @Produce json
// @Param login body dto.Login true "Login request"
// @Success 202 {object} dto.AuthResponse
// @Failure 400,404,500 {object} error.MyError
// @Router /auth/login [post]
func (a *AuthService) Login(c *gin.Context) {
	var req dto.Login
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, error2.FromBindError(err))
		return
	}

	db := c.MustGet("DB").(*gorm.DB)
	token, err := auth.Login(db, req, a.jwt)
	if err != nil {
		error2.FromError(err).FillHTTPContextError(c)
		return
	}

	c.JSON(202, token)
}

// Logout logout a user
// @Summary Logout
// @Description Logs out a user
// @Tags auth
// @Success 204
// @Failure 400,404,500 {object} error.MyError
// @Security Bearer
// @Router /auth/logout [delete]
func (*AuthService) Logout(c *gin.Context) {
	claims := c.MustGet("claims").(*jwt2.Claims)

	db := c.MustGet("DB").(*gorm.DB)
	err := auth.Logout(db, claims.ID)
	if err != nil {
		error2.FromError(err).FillHTTPContextError(c)
		return
	}
	c.JSON(204, gin.H{})
}

// Register register a user
// @Summary Register
// @Description Registers a user
// @Tags auth
// @Accept json
// @Produce json
// @Param register body dto.Register true "Register request"
// @Success 201 {object} dto.AuthResponse
// @Failure 400,404,500 {object} error.MyError
// @Router /auth/register [post]
func (a *AuthService) Register(c *gin.Context) {
	var req dto.Register
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, error2.FromBindError(err))
		return
	}

	db := c.MustGet("DB").(*gorm.DB)
	token, err := auth.Register(db, req, a.jwt)
	if err != nil {
		error2.FromError(err).FillHTTPContextError(c)
		return
	}

	c.JSON(201, token)
}

// ChangePassword change a user's password
// @Summary Change password
// @Description Changes a user's password
// @Tags auth
// @Accept json
// @Produce json
// @Param changePassword body dto.ChangePassword true "Change password request"
// @Success 202 {object} error.MyError
// @Failure 400,404,500 {object} error.MyError
// @Security Bearer
// @Router /auth/change-password [put]
func (*AuthService) ChangePassword(c *gin.Context) {
	claims := c.MustGet("claims").(*jwt2.Claims)
	var req dto.ChangePassword
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, error2.FromBindError(err))
		return
	}

	req.UserId = claims.UserId
	db := c.MustGet("DB").(*gorm.DB)
	err := auth.ChangePassword(db, req)
	if err != nil {
		error2.FromError(err).FillHTTPContextError(c)
		return
	}

	c.JSON(202, gin.H{
		"message":       "password has been changed",
		"short_message": "password changed",
	})
}

// RefreshToken refresh a user's token
// @Summary Refresh token
// @Description Refreshes a user's token
// @Tags auth
// @Accept json
// @Produce json
// @Success 202 {object} dto.AuthResponse
// @Failure 400,404,500 {object} error.MyError
// @Security Bearer
// @Router /auth/refresh-token [post]
func (a *AuthService) RefreshToken(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	refreshToken := c.MustGet("token").(string)

	token, err := auth.RefreshToken(db, a.jwt, refreshToken)
	if err != nil {
		error2.FromError(err).FillHTTPContextError(c)
		return
	}

	c.JSON(202, token)
}

// GetAccount get a user's account
// @Summary Get account
// @Description Gets a user's account information from the token
// @Tags auth
// @Produce json
// @Success 200 {object} dto.User
// @Failure 400,404,500 {object} error.MyError
// @Security Bearer
// @Router /auth/account [get]
func (*AuthService) GetAccount(c *gin.Context) {
	claims := c.MustGet("claims").(*jwt2.Claims)

	am := model.AccountModel{Tx: c.MustGet("DB").(*gorm.DB)}
	a, err := auth.GetAccount(am.Tx, claims.UserId)
	if err != nil {
		error2.FromError(err).FillHTTPContextError(c)
		return
	}

	c.JSON(200, a)
}

// SetAuthService Add auth service to gin engine
func SetAuthService(r *gin.RouterGroup, jwt config.JwtConfig) {
	as := AuthService{jwt: jwt}
	// Set middleware
	mddl := middleware.NewJwtMiddleware(jwt)

	// Set routes
	// Routes without privileges required
	r.POST("/login", as.Login)
	r.POST("/register", as.Register)
	// Add privileges on routes
	r.Use(mddl.MiddlewareFunc(map[string][]enum.Role{
		"Logout":         {enum.STUDENT},
		"ChangePassword": {enum.STUDENT},
		"RefreshToken":   {enum.STUDENT},
		"GetAccount":     {enum.STUDENT},
	}))
	// Routes with privileges required
	r.DELETE("/logout", as.Logout)
	r.PUT("/change-password", as.ChangePassword)
	r.POST("/refresh-token", as.RefreshToken)
	r.GET("/account", as.GetAccount)
}
