package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
)

type DatabaseMiddleware struct {
	DB *gorm.DB
}

func (dm *DatabaseMiddleware) SetDBMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// If the request method is OPTIONS, we don't need to set the database
		if c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		var db *gorm.DB
		timeoutContext, cancel := context.WithTimeout(context.Background(), time.Hour)
		defer cancel()
		// We use the timeout context to avoid the database connection to be blocked
		db = dm.DB.WithContext(timeoutContext).Begin()
		if db.Error != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": "database error"})
			return
		}

		c.Set("DB", db)
		c.Next()

		// If the request method is GET, we don't need to commit the transaction
		// If the request method is POST, PUT, DELETE, we need to commit the transaction
		// If the request method is POST, PUT, DELETE and there is an error, we need to rollback the transaction
		if c.Request.Method != "GET" {
			if c.IsAborted() {
				db.Rollback()
				if db.Error != nil {
					c.AbortWithStatusJSON(500, gin.H{"error": "database error"})
				}
				return
			}
			db.Commit()
			if db.Error != nil {
				c.AbortWithStatusJSON(500, gin.H{"error": "database error"})
				return
			}
		}
	}
}
