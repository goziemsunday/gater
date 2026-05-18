package validator

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	v *validator.Validate
}

func New() *Validator {
	return &Validator{
		v: validator.New(validator.WithRequiredStructEnabled()),
	}
}

func (v *Validator) ValidateStruct(s any) ([]string, bool) {
	err := v.v.Struct(s)
	if err == nil {
		return nil, true
	}

	errs := v.v.Struct(s).(validator.ValidationErrors)
	var errMsgs []string
	for _, e := range errs {
		field := strings.ToLower(e.Field())
		switch e.Tag() {
		case "required":
			errMsgs = append(errMsgs, field+" is required")
		case "email":
			errMsgs = append(errMsgs, "invalid email format")
		case "min":
			errMsgs = append(errMsgs, field+" is too short")
		case "max":
			errMsgs = append(errMsgs, field+" is too long")
		default:
			errMsgs = append(errMsgs, field+" is invalid")
		}
	}

	return errMsgs, false
}

func (v *Validator) ValidateVar(variable any, tag string) error {
	return v.v.Var(variable, tag)
}
