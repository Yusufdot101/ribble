package api

import (
	"strconv"
	"strings"

	"github.com/Yusufdot101/ripple/shared/middleware"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var NewMessageRequest struct {
	ChatID  uint   `json:"chatID"`
	Content string `json:"string"`
}

func (h *handler) newMessage(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	// auth first message
	var authMsg struct {
		Token string `json:"token"`
	}

	if err := conn.ReadJSON(&authMsg); err != nil {
		return
	}

	// validate the token
	token, err := middleware.ValidateJWT(authMsg.Token)
	if err != nil {
		wsError(conn, middleware.ErrInvalidJWT.Error())
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		wsError(conn, middleware.ErrInvalidJWT.Error())
		return
	}

	// extract the user id from it
	userID, ok := claims["sub"].(string)
	if !ok || userID == "" {
		wsError(conn, middleware.ErrInvalidJWT.Error())
		return
	}

	// register connection
	userIDint, err := strconv.Atoi(userID)
	if err != nil {
		conn.WriteJSON(map[string]string{
			"type":    "error",
			"message": middleware.ErrInvalidJWT.Error(),
		})
		return
	}

	h.hub.addClient(uint(userIDint), conn)
	defer h.hub.removeClient(uint(userIDint), conn)

	for {
		var msg struct {
			Type    string `json:"type"`
			ChatID  uint   `json:"chatID"`
			Content string `json:"content"`
		}

		if err := conn.ReadJSON(&msg); err != nil {
			break
		}

		if msg.ChatID == 0 || strings.TrimSpace(msg.Content) == "" {
			wsError(conn, "invalid message")
			continue
		} else {
			participants, err := h.csvc.GetChatParticipants(msg.ChatID)
			if err != nil {
				conn.WriteJSON(map[string]string{
					"type":    "error",
					"message": "chat not found",
				})
				continue
			}

			// make sure the user is in the chat
			userInChat := false
			for _, p := range participants {
				if p.ID == uint(userIDint) {
					userInChat = true
					break
				}
			}
			if !userInChat {
				break
			}

			err = h.csvc.NewMessage(uint(userIDint), msg.ChatID, msg.Content)
			if err != nil {
				conn.WriteJSON(map[string]string{
					"type":    "error",
					"message": "chat not found",
				})
				continue
			}

			for _, p := range participants {
				h.hub.SendToUser(p.UserID, msg)
			}
		}
	}
}
