package http_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	apperror "github.com/exPriceD/simple-marketplace/internal/application/error"
	"github.com/exPriceD/simple-marketplace/internal/delivery/http/response"
)

func TestErrorMapper_Kinds(t *testing.T) {
	tests := []struct {
		kind     apperror.Kind
		code     string
		expected int
	}{
		{apperror.KindValidation, "invalid_title", http.StatusUnprocessableEntity},
		{apperror.KindConflict, "login_taken", http.StatusConflict},
		{apperror.KindAuth, "unauthorized", http.StatusUnauthorized},
		{apperror.KindNotFound, "not_found", http.StatusNotFound},
		{apperror.KindInternal, "internal_error", http.StatusInternalServerError},
	}

	for _, tc := range tests {
		ae := apperror.New("op", tc.code, tc.kind, nil)
		rec := httptest.NewRecorder()
		response.WriteAppError(rec, "rid123", ae)
		if rec.Code != tc.expected {
			t.Fatalf("expected %d got %d for kind %s", tc.expected, rec.Code, tc.kind)
		}
		var body map[string]any
		if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
			t.Fatalf("json err: %v", err)
		}
		if body["error"].(map[string]any)["code"] != tc.code {
			t.Fatalf("wrong code field")
		}
	}
}
