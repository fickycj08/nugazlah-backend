package ierr

import (
	"fmt"
	"log"
	"net/http"

	"github.com/pkg/errors"
)

type customError struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

func (e customError) Error() string {
	return e.Message
}

func ExtendErr(err customError, msg string) error {
	err.Message = fmt.Sprintf("%s, err : %s", err.Message, msg)
	return err
}

var (
	Success             = customError{Code: "C-0000", Message: "Success"}
	ErrInternal         = customError{Code: "C-0002", Message: "Sorry, an internal server error occurred. Please try again later."}
	ErrDuplicate        = customError{Code: "C-0003", Message: "The data you provided conflicts with existing data. Please review the information you entered"}
	ErrNotFound         = customError{Code: "C-0004", Message: "Sorry, the resource you requested could not be found."}
	ErrBadRequest       = customError{Code: "C-0005", Message: "Sorry, the request is invalid. Please check your input and try again."}
	ErrForbidden        = customError{Code: "C-0006", Message: "You do not have permission to access or edit this resource."}
	ErrAlreadyJoinClass = customError{Code: "C-0007", Message: "User already join the class."}
)

func TranslateError(err error) (statusCode int, errCode string, msg string) {
	log.Println(err)

	switch errors.Cause(err) {
	case ErrDuplicate:
		return http.StatusBadRequest, ErrDuplicate.Code, err.Error()
	case ErrNotFound:
		return http.StatusBadRequest, ErrNotFound.Code, err.Error()
	case ErrForbidden:
		return http.StatusBadRequest, ErrForbidden.Code, err.Error()
	case ErrBadRequest:
		return http.StatusBadRequest, ErrBadRequest.Code, err.Error()
	case ErrAlreadyJoinClass:
		return http.StatusBadRequest, ErrAlreadyJoinClass.Code, err.Error()
	}

	return http.StatusInternalServerError, "C-0001", ErrInternal.Message
}
