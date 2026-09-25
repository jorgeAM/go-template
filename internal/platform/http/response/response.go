package response

import (
	"encoding/json"
	"net/http"
)

type errResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func ErrResponse(w http.ResponseWriter, code string, message string) {
	errResp := &errResponse{
		Code:    code,
		Message: message,
	}

	errResponse, err := json.Marshal(errResp)
	if err != nil {
		w.Write([]byte(""))
		return
	}

	w.Write(errResponse)
}

func Created(w http.ResponseWriter, res any) {
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		InternalServerErr(w, "INTERNAL_ERROR", err.Error())
		return
	}
}

func OK(w http.ResponseWriter, res any) {
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		InternalServerErr(w, "INTERNAL_ERROR", err.Error())
		return
	}
}

func BadRequest(w http.ResponseWriter, code string, message string) {
	w.WriteHeader(http.StatusBadRequest)
	ErrResponse(w, code, message)
}

func InternalServerErr(w http.ResponseWriter, code string, message string) {
	w.WriteHeader(http.StatusInternalServerError)
	ErrResponse(w, code, message)
}

func CustomStatusErrResponse(w http.ResponseWriter, code string, message string, httpStatus int) {
	w.WriteHeader(httpStatus)
	ErrResponse(w, code, message)
}
