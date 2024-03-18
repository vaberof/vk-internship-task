package http

import (
	"bytes"
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

type validationErrors []error

func (ve validationErrors) Error() string {
	buff := bytes.NewBufferString("")

	for i := 0; i < len(ve); i++ {

		buff.WriteString(ve[i].Error())
		buff.WriteString("\n")
	}

	return strings.TrimSpace(buff.String())
}

func validateRequestBody(err error) validationErrors {
	var verr validationErrors

	for _, err := range err.(validator.ValidationErrors) {
		var e error
		switch err.Tag() {
		case "required":
			e = fmt.Errorf("Field '%s' must be nonempty", err.Field())
		case "min":
			e = fmt.Errorf("Field '%s' must be equal or greater than '%v' characters long", err.Field(), err.Param())
		case "max":
			e = fmt.Errorf("Field '%s' must be lower or equal than '%v' characters long", err.Field(), err.Param())
		case "gt":
			e = fmt.Errorf("Field '%s' must be greater than '%v'", err.Field(), err.Param())
		case "oneof":
			e = fmt.Errorf("Field '%s' must have one of acceptable values: '%v'", err.Field(), err.Param())
		case "numeric":
			e = fmt.Errorf("Field '%s' must contain numeric values", err.Field())
		default:
			e = fmt.Errorf("Field '%s': '%v' must satisfy '%s' '%v' criteria", err.Field(), err.Value(), err.Tag(), err.Param())
		}
		verr = append(verr, e)
	}

	return verr
}
