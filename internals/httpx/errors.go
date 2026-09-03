package httpx

import (
	"encoding/json"
	"net/http"
)

type Code string

const (
	CodeInvalidId        Code = "invalid_id"
	CodeInternalError    Code = "internal_error"
	CodeNotFound         Code = "not_found"
	CodeMalformedJSON    Code = "malformed_json"
	CodeValidationFailed Code = "validation_failed"
	CodeUnauthenticated  Code = "unauthenticated"
	CodeForbidden        Code = "forbidden"
	CodeConflict         Code = "conflict"
	CodeRateLimited      Code = "rate_limited"
)

type errorEnvelop struct {
	Error errorPayLoad `json:"error"`
}

type errorPayLoad struct {
	Message string `json:"message"`
	Code    Code   `json:"code"`
}

func Error(w http.ResponseWriter, status int, message string, code Code) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)

	err := errorPayLoad{
		Message: message,
		Code:    code,
	}

	errEnv := errorEnvelop{
		Error: err,
	}

	_ = json.NewEncoder(w).Encode(errEnv)
}
