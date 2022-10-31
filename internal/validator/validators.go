package validator

import "net/mail"

type Validator struct {
	Errors map[string]string
}

func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

func (v *Validator) IsValid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) AddErrorToMap(key, message string) {
	if _, errs := v.Errors[key]; !errs {
		v.Errors[key] = message
	}
}

func (v *Validator) Check(isvalid bool, key, message string) {
	if !isvalid {
		v.AddErrorToMap(key, message)
	}
}

func (v *Validator) IsMailValid(value, key, message string) {
	_, err := mail.ParseAddress(value)
	if err != nil {
		v.AddErrorToMap(key, message)
	}
}
