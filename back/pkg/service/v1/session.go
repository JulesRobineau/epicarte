package v1

import (
	"gin-template/config"
	"gin-template/pkg/common/session"
	"gin-template/pkg/middleware"
	"gin-template/pkg/model/enum"
	error2 "gin-template/utils/error"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// GetSession returns a session
// @Summary Get a session
// @Description Get a session by ID
// @Tags session
// @Produce json
// @Param session_id path string true "TinySession ID"
// @Security Bearer
// @Success 200 {object} dto.Session
// @Failure 400,404,500 {object} error.MyError
// @Router /sessions/{session_id} [get]
func GetSession(c *gin.Context) {
	var req struct {
		SessionID string `uri:"session_id" binding:"required,uuid"`
	}

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(400, error2.FromBindError(err))
		return
	}

	s, err := session.GetSessionByID(c.MustGet("DB").(*gorm.DB), uuid.Must(uuid.FromString(req.SessionID)))
	if err != nil {
		error2.FromError(err).FillHTTPContextError(c)
		return
	}

	c.JSON(200, s)
}

// DeleteSession deletes a session
// @Summary Delete a session
// @Description Delete a session by ID
// @Tags session
// @Produce json
// @Param session_id path string true "TinySession ID"
// @Security Bearer
// @Success 200 {object} dto.TinySession
// @Failure 400,404,500 {object} error.MyError
// @Router /sessions/{session_id} [delete]
func DeleteSession(c *gin.Context) {
	var req struct {
		SessionID string `uri:"session_id" binding:"required,uuid"`
	}

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(400, error2.FromBindError(err))
		return
	}

	if err := session.DeleteSession(c.MustGet("DB").(*gorm.DB), uuid.Must(uuid.FromString(req.SessionID))); err != nil {
		error2.FromError(err).FillHTTPContextError(c)
		return
	}

	c.JSON(204, nil)
}

// SetSessionRoutes sets the routes for the session service
func SetSessionRoutes(r *gin.RouterGroup, config config.JwtConfig) {
	mdl := middleware.NewJwtMiddleware(config)
	r.Use(mdl.MiddlewareFunc(map[string][]enum.Role{
		"GetSession":    {enum.ADMIN},
		"DeleteSession": {enum.ADMIN},
	}))

	r.GET("/:session_id", GetSession)
	r.DELETE("/:session_id", DeleteSession)
}
