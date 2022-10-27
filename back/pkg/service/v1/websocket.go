package v1

import (
	"context"
	"fmt"
	"gin-template/pkg/common/session"
	"gin-template/pkg/dto"
	"gin-template/pkg/model"
	error2 "gin-template/utils/error"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"net/http"
	"time"
)

const timeout = time.Minute * 20

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebSocketSession struct {
	Session *model.Session
	DB      *gorm.DB
	Cancel  context.CancelFunc
}

func (ws *WebSocketSession) JoinSessionHandler(w http.ResponseWriter, r *http.Request) {
	defer ws.Cancel()
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Failed to set websocket upgrade: %+v\n", err)
		return
	}

	for {
		t, msg, wErr := conn.ReadMessage()
		if wErr != nil {
			break
		}
		fmt.Println(string(msg))
		conn.WriteMessage(t, msg)
	}
}

// JoinSession joins a session
// @Summary Join a session
// @Description Join a session by session ID. This will create a websocket connection.
// @Tags session
// @Produce json
// @Param session_id path string true "Session ID"
// @Param join_session query dto.JoinSession true "Join Session"
// @Success 200
// @Failure 400,404,500 {object} error.MyError
// @Router /ws/sessions/{session_id}/join [get]
func JoinSession(c *gin.Context) {
	var req dto.JoinSession
	if err := c.ShouldBindQuery(&req); err != nil {
		error2.FromBindError(err).FillHTTPContextError(c)
		return
	}

	var path struct {
		ID string `uri:"session_id" binding:"required,uuid"`
	}
	if err := c.ShouldBindUri(&path); err != nil {
		error2.FromBindError(err).FillHTTPContextError(c)
		return
	}

	db := c.MustGet("DB").(*gorm.DB)
	s, err := session.VerifySession(db, uuid.Must(uuid.FromString(path.ID)), req.Password)
	if err != nil {
		error2.FromError(err).FillHTTPContextError(c)
		return
	}

	timeoutContext, cancel := context.WithTimeout(context.Background(), timeout)
	ws := WebSocketSession{
		Session: s,
		DB:      db.WithContext(timeoutContext),
		Cancel:  cancel,
	}

	ws.JoinSessionHandler(c.Writer, c.Request)
}

// SetWebsocketRoutes sets the websocket router
func SetWebsocketRoutes(r *gin.RouterGroup) {
	r.GET("/sessions/:session_id/join", JoinSession)
}
