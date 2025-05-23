package handlers

import (
	"errors"
	"io"
	"net/http"
	"sentimenta/internal/utils"
	"sentimenta/internal/ws"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type WSHandler struct {
	logger  *zap.SugaredLogger
	connMgr *ws.ConnectionManager
	// config *config.Config
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *WSHandler) HandleWS(c echo.Context) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		h.logger.Error("failed to upgrade to websocket: ", err)
		return err
	}
	defer func() {
		if err := conn.Close(); err != nil {
			h.logger.Errorf("Failed to close ws connection: %v\n", err)
		}
	}()
	h.connMgr.Add(userID, conn)
	defer h.connMgr.Remove(userID)

	h.logger.Infof("User %s connected via WebSocket", userID)
	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			if closeErr, ok := err.(*websocket.CloseError); ok {
				switch closeErr.Code {
				case websocket.CloseNormalClosure, websocket.CloseGoingAway:
					h.logger.Infof("WS: normal closure by user %s: %v", userID, closeErr)
				default:
					h.logger.Warnf("WS: abnormal closure by user %s: %v", userID, closeErr)
				}
			} else if errors.Is(err, io.EOF) {
				h.logger.Warnf("WS: EOF from user %s: %v", userID, err)
			} else {
				h.logger.Errorf("WS: read error from user %s: %v", userID, err)
			}
			break
		}
		reply := "Echo from server: " + string(msg)
		if err := conn.WriteMessage(msgType, []byte(reply)); err != nil {
			h.logger.Error("write error:", err)
			break
		}
	}
	return nil
}

func NewWSHandler(logger *zap.SugaredLogger, connMgr *ws.ConnectionManager) *WSHandler {
	return &WSHandler{logger: logger, connMgr: connMgr}
}
