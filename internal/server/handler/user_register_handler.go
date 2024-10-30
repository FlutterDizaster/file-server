package handler

import (
	"net/http"
)

func (h Handler) userRegisterHandler(w http.ResponseWriter, r *http.Request) {
	h.userHandler(w, r, h.userCtrl.Register)
}
