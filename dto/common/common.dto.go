package common_dto

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type StdResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"message"`
}

func BuildErrorResponse(code int, e interface{}) StdResponse {
	xType := fmt.Sprintf("%T", e)
	validationType := "validator.ValidationErrors"
	if validationType == xType {
		msgsOfErrors := []string{}

		if errs, ok := e.(validator.ValidationErrors); ok {
			for _, fieldErr := range []validator.FieldError(errs) {
				msgsOfErrors = append(msgsOfErrors, customValidationError(fieldErr))
			}
		}
		msg := "There is an error with your request"
		if len(msgsOfErrors) > 0 {
			msg = msgsOfErrors[0]
		}
		return StdResponse{
			Code: code,
			Msg:  msg,
		}
	}
	msg := e.(string)
	return StdResponse{
		Code: code,
		Msg:  msg,
	}
}

func customValidationError(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required.", err.Field())
	case "min":
		return fmt.Sprintf("%s must be longer than or equal %s characters.", err.Field(), err.Param())
	case "max":
		return fmt.Sprintf("%s cannot be longer than %s characters.", err.Field(), err.Param())
	default:
		msg := fmt.Sprintf("%v", err)
		return msg
	}
}
