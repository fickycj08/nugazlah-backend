package response

import (
	"encoding/json"
	"net/http"

	"github.com/vandenbill/nugazlah-backend/internal/ierr"
)

type (
	ErrorResponse struct {
		Message string `json:"message"`
		Code    string `json:"code"`
	}

	SuccessReponse struct {
		Message string `json:"message"`
		Data    any    `json:"data"`
		Code    string `json:"code"`
	}

	SuccessPageReponse struct {
		Message string `json:"message"`
		Data    any    `json:"data"`
		Meta    Meta   `json:"meta"`
	}

	Meta struct {
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
		Total  int `json:"total"`
	}
)

func ResponseErrWithCode(msg string, code int, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(ErrorResponse{
		Message: msg,
	})
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func ResponseErr(err error, w http.ResponseWriter) {
	statusCode, errCode, msg := ierr.TranslateError(err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err = json.NewEncoder(w).Encode(ErrorResponse{
		Message: msg,
		Code:    errCode,
	})
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func ResponseSuccess(msg string, data any, w http.ResponseWriter) {
	successRes := SuccessReponse{}
	successRes.Message = msg
	successRes.Data = data
	successRes.Code = ierr.Success.Code

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err := json.NewEncoder(w).Encode(successRes)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
