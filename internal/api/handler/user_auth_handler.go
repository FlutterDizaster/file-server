package handler

import (
	"net/http"
)

// FIXME: I'm not sure about this solution, but it seems to be ok.
// Same with userRegisterHandler

func (h Handler) userAuthHandler(w http.ResponseWriter, r *http.Request) {
	h.userHandler(w, r, h.userCtrl.Login)
}
