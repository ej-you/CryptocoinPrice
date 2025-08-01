// Package validator provides interface to validate struct data by tags.
package validator

import (
	"errors"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	govalidator "github.com/go-playground/validator/v10"
	entranslation "github.com/go-playground/validator/v10/translations/en"
)

var _ Validator = (*valid)(nil)

// Validator provides method to validate any struct.
type Validator interface {
	Validate(s any) error
}

// Validator implementation.
type valid struct {
	validatorInstance *govalidator.Validate
	translator        ut.Translator
}

// New returns new validator.
func New() Validator {
	enTranslator := en.New()
	uni := ut.New(enTranslator, enTranslator)
	trans, _ := uni.GetTranslator("en")

	validate := govalidator.New(govalidator.WithRequiredStructEnabled())
	err := entranslation.RegisterDefaultTranslations(validate, trans)
	if err != nil {
		panic(err)
	}

	return &valid{validate, trans}
}

// Validate validates given struct s (using pointer to this struct).
func (v valid) Validate(s any) error {
	err := v.validatorInstance.Struct(s)
	if err == nil { // NOT err
		return nil
	}

	var validateErrors govalidator.ValidationErrors
	if !errors.As(err, &validateErrors) {
		return err
	}
	// handle error messages
	rawTranstaledMap := validateErrors.Translate(v.translator)

	// sort out errors and concat them into string
	transtaledStringSlice := make([]string, 0, len(rawTranstaledMap))
	for _, v := range rawTranstaledMap {
		transtaledStringSlice = append(transtaledStringSlice, strings.ToLower(v))
	}

	return errors.New(strings.Join(transtaledStringSlice, " && "))
}
