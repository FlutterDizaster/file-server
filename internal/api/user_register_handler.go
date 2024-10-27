package api

import (
	"net/http"
)

func (a API) userRegisterHandler(w http.ResponseWriter, r *http.Request) {
	userHandler(w, r, a.userCtrl.Register)
}
