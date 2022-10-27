package v1

import (
	"gin-template/config"
	user2 "gin-template/pkg/common/user"
	"gin-template/pkg/dto"
	"gin-template/pkg/middleware"
	"gin-template/pkg/model/enum"
	error2 "gin-template/utils/error"
	jwt2 "gin-template/utils/jwt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetUser returns a user
// @Summary Get a user
// @Description Get a user
// @Tags user
// @Produce json
// @Param user_id path int true "User ID"
// Security Bearer
// @Success 200 {object} dto.User
// @Failure 400,404,500 {object} error.MyError
// @Security Bearer
// @Router /users/{user_id} [get]
func GetUser(c *gin.Context) {
	// Get the user if from the path
	var req struct {
		UserId uint64 `uri:"user_id" binding:"required"`
	}
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(400, error2.FromBindError(err))
		return
	}

	claims := c.MustGet("claims").(*jwt2.Claims)
	if req.UserId != claims.UserId && claims.Role != enum.SUPERADMIN {
		error2.ForbiddenError("You can only access your own account").FillHTTPContextError(c)
		return
	}

	db := c.MustGet("DB").(*gorm.DB)
	// Get the user from common
	user, err := user2.GetUser(db, req.UserId)
	if err != nil {
		error2.FromError(err).FillHTTPContextError(c)
		return
	}

	// Return the user
	c.JSON(200, user)
}

// UserList returns a list of users
// @Summary Get a list of users
// @Description Get a list of users with pagination
// @Param params query dto.UserQueryParams false "..."
// @Tags user
// @Produce json
// @Security Bearer
// @Success 200 {object} []dto.User
// @Failure 400,404,500 {object} error.MyError
// @Router /users [get]
func UserList(c *gin.Context) {
	var req dto.UserQueryParams
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(400, error2.FromBindError(err))
		return
	}

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(400, error2.FromBindError(err))
		return
	}

	db := c.MustGet("DB").(*gorm.DB)
	users, err := user2.GetUserList(db, req)
	if err != nil {
		error2.FromError(err).FillHTTPContextError(c)
		return
	}

	c.JSON(200, gin.H{"users": users})
}

// UpdateUser updates a user
// @Summary Update a user
// @Description Update a user by id
// @Tags user
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Param user body dto.UpdateUser true "User"
// Security Bearer
// @Success 200 {object} dto.User
// @Failure 400,404,500 {object} error.MyError
// @Security Bearer
// @Router /users/{user_id} [put]
func UpdateUser(c *gin.Context) {
	var req dto.UpdateUser
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, error2.FromBindError(err))
		return
	}

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(400, error2.FromBindError(err))
		return
	}

	claims := c.MustGet("claims").(*jwt2.Claims)
	if req.Id != claims.UserId && claims.Role != enum.SUPERADMIN {
		error2.ForbiddenError("You can only update your account").FillHTTPContextError(c)
		return
	}

	db := c.MustGet("DB").(*gorm.DB)
	var err error
	switch claims.Role {
	case enum.SUPERADMIN:
		// Superadmin can update any user
		err = user2.SuperAdminUpdateUser(db, req)
	default:
		// User can only update himself
		err = user2.UpdateUser(db, req)
	}
	if err != nil {
		error2.FromError(err).FillHTTPContextError(c)
		return
	}

	c.JSON(202, gin.H{"message": "User updated"})
}

func SetUserRoutes(r *gin.RouterGroup, config config.JwtConfig) {
	// Setup middleware
	mdl := middleware.NewJwtMiddleware(config)

	// Add privileges on routes
	r.Use(mdl.MiddlewareFunc(map[string][]enum.Role{
		"GetUser":    {enum.SUPERADMIN, enum.STUDENT},
		"UserList":   {enum.SUPERADMIN},
		"UpdateUser": {enum.SUPERADMIN, enum.STUDENT},
	}))

	// Add routes to the router
	r.GET("/:user_id", GetUser)
	r.GET("", UserList)
	r.PUT("/:user_id", UpdateUser)
}
