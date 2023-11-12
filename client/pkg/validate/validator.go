package validate

import (
	"log"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var (
	trans ut.Translator
)

func SetUp(validate *validator.Validate) {
	uni := ut.New(en.New())
	trans, _ = uni.GetTranslator("en")
	err := en_translations.RegisterDefaultTranslations(validate, trans)
	if err != nil {
		log.Fatalf("Setup validator failed: %s", err.Error())
	}

	// Rewrite struct field to JSON tag
	validate.RegisterTagNameFunc(func(fl reflect.StructField) string {
		name := strings.SplitN(fl.Tag.Get("json"), ",", 2)
		return name[0]
	})

}

func FormatValidationError(errs validator.ValidationErrors) map[string]string {
	fieldErrors := make(map[string]string)
	for _, e := range errs {
		fieldErrors[e.Field()] = e.Translate(trans)
	}

	return fieldErrors
}
