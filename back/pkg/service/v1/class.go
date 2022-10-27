package v1

import (
	"gin-template/config"
	"gin-template/pkg/common/class"
	"gin-template/pkg/common/session"
	"gin-template/pkg/dto"
	"gin-template/pkg/middleware"
	"gin-template/pkg/model/enum"
	error2 "gin-template/utils/error"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// GetClass returns a class
// @Summary Get a class
// @Description Get a class by ID
// @Tags class
// @Produce json
// @Param class_id path int true "TinyClass ID"
// @Security Bearer
// @Success 200 {object} dto.Class
// @Failure 400,404,500 {object} error.MyError
// @Router /classes/{class_id} [get]
func GetClass(c *gin.Context) {
	var req struct {
		ClassID uint64 `uri:"class_id" binding:"required"`
	}
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(400, error2.FromBindError(err))
		return
	}

	cl, err := class.GetClassByID(c.MustGet("DB").(*gorm.DB), req.ClassID)
	if err != nil {
		error2.FromError(err).FillHTTPContextError(c)
		return
	}

	c.JSON(200, cl)
}

// ClassList returns a list of classes
// @Summary Get a list of classes
// @Description Get a list of classes
// @Tags class
// @Produce json
// @Security Bearer
// @Success 200 {object} dto.ClassList
// @Failure 400,404,500 {object} error.MyError
// @Router /classes [get]
func ClassList(c *gin.Context) {
	classes, err := class.GetAllClasses(c.MustGet("DB").(*gorm.DB))
	if err != nil {
		error2.FromError(err).FillHTTPContextError(c)
		return
	}

	c.JSON(200, classes)
}

// CreateClass creates a class
// @Summary Create a class
// @Description Create a class
// @Tags class
// @Produce json
// @Param class body dto.CreateClass true "TinyClass"
// @Security Bearer
// @Success 201 {object} dto.TinyClass
// @Failure 400,404,500 {object} error.MyError
// @Router /classes [post]
func CreateClass(c *gin.Context) {
	var req dto.CreateClass
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, error2.FromBindError(err))
		return
	}

	cl, err := class.CreateClass(c.MustGet("DB").(*gorm.DB), req)
	if err != nil {
		error2.FromError(err).FillHTTPContextError(c)
		return
	}

	c.JSON(201, cl)
}

// UpdateClass updates a class
// @Summary Update a class
// @Description Update a class
// @Tags class
// @Produce json
// @Param class_id path int true "TinyClass ID"
// @Param class body dto.UpdateClass true "TinyClass"
// @Security Bearer
// @Success 202
// @Failure 400,404,500 {object} error.MyError
// @Router /classes/{class_id} [put]
func UpdateClass(c *gin.Context) {
	var req dto.UpdateClass
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, error2.FromBindError(err))
		return
	}

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(400, error2.FromBindError(err))
		return
	}

	if err := class.UpdateClass(c.MustGet("DB").(*gorm.DB), req); err != nil {
		error2.FromError(err).FillHTTPContextError(c)
		return
	}

	c.JSON(202, gin.H{"message": "TinyClass updated"})
}

// DeleteClass deletes a class
// @Summary Delete a class
// @Description Delete a class
// @Tags class
// @Produce json
// @Param class_id path int true "TinyClass ID"
// @Security Bearer
// @Success 204
// @Failure 400,404,500 {object} error.MyError
// @Router /classes/{class_id} [delete]
func DeleteClass(c *gin.Context) {
	var req struct {
		ClassID uint64 `uri:"class_id" binding"required"`
	}
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(400, error2.FromBindError(err))
		return
	}

	if err := class.DeleteClass(c.MustGet("DB").(*gorm.DB), req.ClassID); err != nil {
		error2.FromError(err).FillHTTPContextError(c)
		return
	}

	c.JSON(204, nil)
}

// ClassSessionList returns a list of class sessions
// @Summary Get a list of class sessions
// @Description Get a list of class sessions
// @Tags class
// @Produce json
// @Param class_id path int true "Class ID"
// @Security Bearer
// @Success 200 {object} dto.SessionList
// @Failure 400,404,500 {object} error.MyError
// @Router /classes/{class_id}/sessions [get]
func ClassSessionList(c *gin.Context) {
	var req struct {
		ClassID uint64 `uri:"class_id" binding:"required"`
	}
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(400, error2.FromBindError(err))
		return
	}

	s, err := session.GetSessionsByClassID(c.MustGet("DB").(*gorm.DB), req.ClassID)
	if err != nil {
		error2.FromError(err).FillHTTPContextError(c)
		return
	}

	c.JSON(200, s)
}

// CreateClassSession creates a class session
// @Summary Create a class session
// @Description Create a class session
// @Tags class
// @Produce json
// @Param class_id path int true "Class ID"
// @Param class body dto.CreateSession true "Class Session"
// @Security Bearer
// @Success 201 {object} dto.TinySession
// @Failure 400,404,500 {object} error.MyError
// @Router /classes/{class_id}/sessions [post]
func CreateClassSession(c *gin.Context) {
	var req dto.CreateSession
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, error2.FromBindError(err))
		return
	}

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(400, error2.FromBindError(err))
		return
	}

	s, err := session.CreateSession(c.MustGet("DB").(*gorm.DB), req)
	if err != nil {
		error2.FromError(err).FillHTTPContextError(c)
		return
	}

	c.JSON(201, s)
}

// CloseClassSession closes a class session
// @Summary Close a class session
// @Description Close a class session
// @Tags class
// @Produce json
// @Param class_id path int true "Class ID"
// @Param session_id path string true "Session ID"
// @Security Bearer
// @Success 202
// @Failure 400,404,500 {object} error.MyError
// @Router /classes/{class_id}/sessions/{session_id} [put]
func CloseClassSession(c *gin.Context) {
	var req struct {
		ClassID   uint64 `uri:"class_id" binding:"required"`
		SessionID string `uri:"session_id" binding:"required,uuid"`
	}

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(400, error2.FromBindError(err))
		return
	}

	if err := session.CloseSession(
		c.MustGet("DB").(*gorm.DB),
		req.ClassID,
		uuid.Must(uuid.FromString(req.SessionID)),
	); err != nil {
		error2.FromError(err).FillHTTPContextError(c)
		return
	}

	c.JSON(202, gin.H{"message": "Session closed"})
}

// DeleteClassSession deletes a class session
// @Summary Delete a class session
// @Description Delete a class session
// @Tags class
// @Produce json
// @Param class_id path int true "Class ID"
// @Security Bearer
// @Success 204
// @Failure 400,404,500 {object} error.MyError
// @Router /classes/{class_id}/sessions [delete]
func DeleteClassSession(c *gin.Context) {
	var req struct {
		ClassID uint64 `uri:"class_id" binding:"required"`
	}
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(400, error2.FromBindError(err))
		return
	}

	if err := session.DeleteSessionsByClassID(c.MustGet("DB").(*gorm.DB), req.ClassID); err != nil {
		error2.FromError(err).FillHTTPContextError(c)
		return
	}

	c.JSON(204, nil)
}

// AddStudentToClass adds a student to a class
// @Summary Add a student to a class
// @Description Add a student to a class
// @Tags class
// @Produce json
// @Param class_id path int true "Class ID"
// @Param student body dto.AddStudentToClass true "Student"
// @Security Bearer
// @Success 201
// @Failure 400,404,500 {object} error.MyError
// @Router /classes/{class_id}/students [post]
func AddStudentToClass(c *gin.Context) {
	var req dto.AddStudentToClass
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, error2.FromBindError(err))
		return
	}

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(400, error2.FromBindError(err))
		return
	}

	st, err := class.AddStudentToClass(c.MustGet("DB").(*gorm.DB), req)
	if err != nil {
		error2.FromError(err).FillHTTPContextError(c)
		return
	}

	c.JSON(201, st)
}

// SetClassRoutes sets up the class routes
func SetClassRoutes(r *gin.RouterGroup, jwtConfig config.JwtConfig) {
	mdl := middleware.NewJwtMiddleware(jwtConfig)
	r.Use(mdl.MiddlewareFunc(map[string][]enum.Role{
		"GetClass":           {enum.ADMIN},
		"ClassList":          {enum.ADMIN},
		"CreateClass":        {enum.ADMIN},
		"UpdateClass":        {enum.ADMIN},
		"DeleteClass":        {enum.ADMIN},
		"ClassSessionList":   {enum.ADMIN},
		"CreateClassSession": {enum.ADMIN},
		"CloseClassSession":  {enum.ADMIN},
		"DeleteClassSession": {enum.ADMIN},
		"AddStudentToClass":  {enum.ADMIN},
	}))
	r.GET("", ClassList)
	r.GET("/:class_id", GetClass)
	r.POST("", CreateClass)
	r.PUT("/:class_id", UpdateClass)
	r.DELETE("/:class_id", DeleteClass)
	r.GET("/:class_id/sessions", ClassSessionList)
	r.POST("/:class_id/sessions", CreateClassSession)
	r.PUT("/:class_id/sessions/:session_id", CloseClassSession)
	r.DELETE("/:class_id/sessions", DeleteClassSession)
	r.POST("/:class_id/students", AddStudentToClass)
}
