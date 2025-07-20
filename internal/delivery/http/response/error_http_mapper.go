package response

import (
	"net/http"

	apperror "github.com/exPriceD/simple-marketplace/internal/application/error"
)

type errorBody struct {
	Error struct {
		Code      string `json:"code"`
		Message   string `json:"message"`
		RequestID string `json:"request_id"`
	} `json:"error"`
}

func WriteAppError(w http.ResponseWriter, rid string, err error) {
	status, code, msg := mapError(err)
	var body errorBody
	body.Error.Code = code
	body.Error.Message = msg
	body.Error.RequestID = rid
	WriteJSON(w, status, body)
}

func WriteInternalError(w http.ResponseWriter, rid string) {
	var body errorBody
	body.Error.Code = "internal_error"
	body.Error.Message = "internal error"
	body.Error.RequestID = rid
	WriteJSON(w, http.StatusInternalServerError, body)
}

func mapError(err error) (int, string, string) {
	ae, ok := err.(*apperror.AppError)
	if !ok {
		return http.StatusInternalServerError, "internal_error", "internal error"
	}
	switch ae.Kind {
	case apperror.KindValidation:
		return http.StatusUnprocessableEntity, ae.Code, "validation error"
	case apperror.KindConflict:
		return http.StatusConflict, ae.Code, "conflict"
	case apperror.KindAuth:
		return http.StatusUnauthorized, ae.Code, "unauthorized"
	case apperror.KindNotFound:
		return http.StatusNotFound, ae.Code, "not found"
	default:
		return http.StatusInternalServerError, "internal_error", "internal error"
	}
}
