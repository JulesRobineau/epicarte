package error

import (
	"errors"
	"fmt"
	"gin-template/logging"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgconn"
	"gorm.io/gorm"
	"strings"
)

type MyError struct {
	Code         int               `json:"-" swaggerignore:"true"`
	ShortMessage string            `json:"short_message"`
	Message      string            `json:"message"`
	Fields       map[string]string `json:"error,omitempty"`
}

// Error returns a string representation of the error
func (e *MyError) Error() string {
	return e.Message
}

// NewError returns a new struct of MyError
func NewError(code int, shortMessage string, message string, fields map[string]string) *MyError {
	return &MyError{
		Code:         code,
		ShortMessage: shortMessage,
		Message:      message,
		Fields:       fields,
	}
}

// UnauthorizedError returns a new struct of MyError with code 401
func UnauthorizedError(message string) *MyError {
	if message == "" {
		message = "You are not authorized to access this resource"
	}
	return NewError(401, "Unauthorized", message, nil)
}

// ForbiddenError returns a new struct of MyError with code 403
func ForbiddenError(message string) *MyError {
	if message == "" {
		message = "You are not allowed to access this resource"
	}
	return NewError(403, "Forbidden", message, nil)
}

// InternalServerError returns a new struct of MyError with code 500
func InternalServerError(message string, err error) *MyError {
	if err != nil {
		logging.Error.Println(err)
	}

	if message == "" {
		message = "Internal Server Error"
	}
	return NewError(500, "internal server error", message, nil)
}

// BadRequestError returns a new struct of MyError with code 400
func BadRequestError(message string, fields map[string]string) *MyError {
	if message == "" {
		message = "Body malformatted"
	}
	return NewError(400, "Bad Request", message, fields)
}

// NotFoundError returns a new struct of MyError with code 404
func NotFoundError(message string) *MyError {
	if message == "" {
		message = "Resource not found"
	}
	return NewError(404, "Not Found", message, nil)
}

// FromBindError returns a new struct of MyError with code 400
func FromBindError(err error) *MyError {
	verr, ok := err.(validator.ValidationErrors)
	if !ok {
		return BadRequestError(err.Error(), nil)
	}

	errs := make(map[string]string)
	for _, f := range verr {
		tag := f.ActualTag()
		switch tag {
		case "required":
			tag = fmt.Sprintf("%s is required", f.Field())
		case "max":
			tag = fmt.Sprintf("%s cannot be longer than %s", f.Field(), f.Param())
		case "min":
			tag = fmt.Sprintf("%s must be longer than %s", f.Field(), f.Param())
		case "email":
			tag = fmt.Sprintf("Invalid email format")
		case "len":
			tag = fmt.Sprintf("%s must be %s characters long", f.Field(), f.Param())
		case "oneof":
			tag = fmt.Sprintf("%s must be one of %s", f.Field(), f.Param())
		case "eqfield":
			tag = fmt.Sprintf("%s must be equal to %s", f.Field(), f.Param())
		case "eqcsfield":
			tag = fmt.Sprintf("%s must be equal to %s", f.Field(), f.Param())
		case "nefield":
			tag = fmt.Sprintf("%s must not be equal to %s", f.Field(), f.Param())
		case "gtfield":
			tag = fmt.Sprintf("%s must be greater than %s", f.Field(), f.Param())
		case "gtefield":
			tag = fmt.Sprintf("%s must be greater than or equal to %s", f.Field(), f.Param())
		case "ltfield":
			tag = fmt.Sprintf("%s must be less than %s", f.Field(), f.Param())
		case "ltefield":
			tag = fmt.Sprintf("%s must be less than or equal to %s", f.Field(), f.Param())
		case "gtcsfield":
			tag = fmt.Sprintf("%s must be greater than %s", f.Field(), f.Param())
		case "gtecsfield":
			tag = fmt.Sprintf("%s must be greater than or equal to %s", f.Field(), f.Param())
		case "ltcsfield":
			tag = fmt.Sprintf("%s must be less than %s", f.Field(), f.Param())
		case "ltecsfield":
			tag = fmt.Sprintf("%s must be less than or equal to %s", f.Field(), f.Param())
		case "eq":
			tag = fmt.Sprintf("%s must be equal to %s", f.Field(), f.Param())
		case "ne":
			tag = fmt.Sprintf("%s must not be equal to %s", f.Field(), f.Param())
		case "gt":
			tag = fmt.Sprintf("%s must be greater than %s", f.Field(), f.Param())
		case "gte":
			tag = fmt.Sprintf("%s must be greater than or equal to %s", f.Field(), f.Param())
		case "lt":
			tag = fmt.Sprintf("%s must be less than %s", f.Field(), f.Param())
		case "lte":
			tag = fmt.Sprintf("%s must be less than or equal to %s", f.Field(), f.Param())
		case "alpha":
			tag = fmt.Sprintf("%s must contain only letters", f.Field())
		case "alphanum":
			tag = fmt.Sprintf("%s must contain only letters and numbers", f.Field())
		case "numeric":
			tag = fmt.Sprintf("%s must contain only numbers", f.Field())
		case "number":
			tag = fmt.Sprintf("%s must contain only numbers", f.Field())
		default:
			tag = fmt.Sprintf("%s is invalid", f.Field())
		}

		errs[f.Field()] = tag
	}

	return &MyError{
		Code:         400,
		ShortMessage: "Bad Request",
		Message:      "Body validation failed",
		Fields:       errs,
	}
}

// FromError returns a new struct of MyError from an error
// If the error is already a MyError, it will be returned as is
// If the error is nil, it will return MyError with code 500
// If the error is not nil and not a MyError, it will return MyError with code 500
func FromError(err error) *MyError {
	if err == nil {
		return InternalServerError("", err)
	}

	if se, ok := err.(*MyError); ok {
		return se
	}

	return InternalServerError("", err)
}

// handlePostgresError returns a new struct of MyError from a GormError
func handlePostgresError(err *pgconn.PgError) *MyError {
	switch err.Code {
	case "23505":
		msg := strings.Replace(err.Detail, "(", "", -1)
		msg = strings.Replace(msg, ")", "", -1)
		return BadRequestError(msg, nil)
	default:
		return InternalServerError("", err)
	}
}

// FromDatabaseError returns a new struct of MyError from a database error
func FromDatabaseError(err error) *MyError {
	if err == nil {
		return InternalServerError("", err)
	}

	pgErr, ok := err.(*pgconn.PgError)
	if ok {
		return handlePostgresError(pgErr)
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return NotFoundError(err.Error())
	}

	return FromError(err)
}

// FillHTTPContextError fills the HTTP context with the error
func (e *MyError) FillHTTPContextError(c *gin.Context) {
	c.AbortWithStatusJSON(e.Code, e)
}
