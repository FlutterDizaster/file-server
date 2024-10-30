package handler

import (
	"net/http"
)

func (h Handler) userAuthHandler(w http.ResponseWriter, r *http.Request) {
	h.userHandler(w, r, h.userCtrl.Login)
}
