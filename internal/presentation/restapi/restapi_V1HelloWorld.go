package restapi

import (
	"net/http"
)

func (h *restApi) HelloWorld(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("hello world toml"))
}
