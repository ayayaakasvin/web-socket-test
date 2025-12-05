package handlers

import (
	"net/http"
	"web-socket-test/internal/libs/validinput"
	"web-socket-test/internal/models"
	"web-socket-test/internal/models/dto"
	"web-socket-test/internal/models/response"
	"web-socket-test/internal/models/token"
)

func (h *Handlers) WS_Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.URL.Query().Get("token")
		if tokenString == "" {
			response.SendErrorJson(w, http.StatusUnauthorized, "missing token")
			h.logger.Warn("WebSocket connection attempt without token")
			return
		}

		cl, err := h.jwtM.Validate(tokenString, &token.AccessTokenClaims{})
		if err != nil {
			response.SendErrorJson(w, http.StatusUnauthorized, "invalid token")
			h.logger.WithError(err).Warn("WebSocket token validation failed")
			return
		}

		claims, ok := cl.(*token.AccessTokenClaims)
		if !ok {
			response.SendErrorJson(w, http.StatusUnauthorized, "invalid claims")
			return
		}

		if claims.UserID == 0 {
			response.SendErrorJson(w, http.StatusUnauthorized, "invalid user_id")
			return
		}

		userInfo, err := h.userRepo.GetPrivateUserInfo(r.Context(), claims.UserID)
		if err != nil {
			response.SendErrorJson(w, http.StatusInternalServerError, "failed to find user or internal server error")
			h.logger.WithError(err).Warn("User not found or Repository is failing")
			return
		}

		conn, err := h.upg.Upgrade(w, r, nil)
		if err != nil {
			response.SendErrorJson(w, http.StatusInternalServerError, "failed to upgrade")
			h.logger.WithError(err).Error("WebSocket upgrade failed")
			return
		}

		nClient := models.NewClient(conn, userInfo)

		h.clientS.Register(nClient)
		h.chatS.PushMessage(dto.SystemMessage(nClient.UserInfo.ID, nClient.UserInfo.Username, dto.ConnectType))

		h.logger.WithField("addr", conn.RemoteAddr().String()).WithField("Connection ID", nClient.ConnectionID).Info("Client registerred")

		go h.readLoop(nClient)
	}
}

func (h *Handlers) readLoop(c *models.Client) {
	for {
		msg := new(dto.WBSMessage)
		if err := c.Conn.ReadJSON(msg); err != nil {
			h.logger.WithError(err).WithField("connection ID", c.ConnectionID).Warn("read failed")
			h.clientS.Unregister(c.UserInfo.ID)
			h.chatS.PushMessage(dto.SystemMessage(c.UserInfo.ID, c.UserInfo.Username, dto.DisconnectType))
			return
		}

		if err := validinput.ValidateWBSMessage(c, msg); err != nil {
			h.logger.WithError(err).WithField("connection ID", c.ConnectionID).Warn("message validating error")
			return
		}

		h.chatS.PushMessage(msg)
	}
}
