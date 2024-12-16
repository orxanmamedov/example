package auth

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"
)

type httpError struct {
	Meta metaError `json:"meta"`
}

type metaError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func errorJSON(w http.ResponseWriter, r *http.Request, code int, err error) {
	var resp httpError

	resp.Meta.Code = code
	resp.Meta.Message = fmt.Sprintf("error: %s", err.Error())
	render.Status(r, code)
	render.JSON(w, r, resp)
}
