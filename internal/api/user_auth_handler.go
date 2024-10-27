package api

import (
	"net/http"
)

// FIXME: I'm not sure about this solution, but it seems to be ok.
// Same with userRegisterHandler

func (a API) userAuthHandler(w http.ResponseWriter, r *http.Request) {
	userHandler(w, r, a.userCtrl.Login)
}
