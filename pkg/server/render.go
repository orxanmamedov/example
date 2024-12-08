package server

import (
	"net/http"

	"github.com/go-chi/render"
)

type httpError struct {
	Meta meta `json:"meta"`
}

type meta struct {
	Message string        `json:"message"`
	DebugID string        `json:"debug_id,omitempty"`
	Errors  []ErrorDetail `json:"errors,omitempty"`
	Code    int           `json:"code"`
}

type ErrorDetail struct {
	Code    int    `json:"code,omitempty"`
	Field   string `json:"field,omitempty"`
	Message string `json:"message,omitempty"`
}

func ResponseJSON(w http.ResponseWriter, r *http.Request, obj interface{}) {
	if obj == nil {
		obj = struct {
		}{}
	}
	render.JSON(w, r, obj)
}

func ResponseJSONWithCode(w http.ResponseWriter, r *http.Request, code int, obj interface{}) {
	if obj == nil {
		obj = struct {
		}{}
	}
	render.Status(r, code)
	render.JSON(w, r, obj)
}

func ErrorJSON(w http.ResponseWriter, r *http.Request, code int, err error, errs ...ErrorDetail) {
	var resp httpError

	resp.Meta.Code = code
	resp.Meta.Message = err.Error()
	resp.Meta.Errors = errs

	render.Status(r, code)
	render.JSON(w, r, &resp)
}

func ErrorCustomJSON(w http.ResponseWriter, r *http.Request, code int, obj interface{}) {
	render.Status(r, code)
	render.JSON(w, r, &obj)
}
