package handlers

import (
	"net/http"

	"github.com/ayayaakasvin/web-socket-test/internal/models/response"
)

func (h *Handlers) GetClientList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clients := h.clientS.Snapshot()

		data := response.NewData()
		data["clients"] = clients

		response.SendSuccessJson(w, http.StatusOK, data)
	}
}
