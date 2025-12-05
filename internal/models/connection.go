package models

import (
	"fmt"
	"net/http"
	"strings"
)

type ID string 

func NewID(r *http.Request) ID {
	ip := strings.Split(r.RemoteAddr, ":")[0]
	agent := r.UserAgent()

	return ID(fmt.Sprintf("%s:%s", ip, agent))
}