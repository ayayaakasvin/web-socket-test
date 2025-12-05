package handlers

import (
	"net/http"
	"web-socket-test/internal/models/response"
)

const (
	defaultLimit = 10
)

func (h *Handlers) GetChatHistory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		history, err := h.chatRepo.GetRecentMessage(r.Context(), defaultLimit)
		if err != nil {
			response.SendErrorJson(w, http.StatusInternalServerError, "failed to fetch chat history")
			h.logger.WithError(err).Error("Chat Repository failed to fetch chat history")
			return
		}

		data := response.NewData()
		data["chat_history"] = history
		
		response.SendSuccessJson(w, http.StatusOK, data)
	}
}