package validators

import (
	"errors"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var (
	Validator  *validator.Validate
	Translator ut.Translator

	// ErrTranslatorNotFound is returned when the translator is not found
	ErrTranslatorNotFound = errors.New("translator not found")
)

// InitValidator initializes the validator and translator
func InitValidator() error {
	// Create a new validator instance
	Validator = validator.New()

	// Create a new translator
	en := en.New()
	uni := ut.New(en, en)

	// Get the English translator
	var found bool
	Translator, found = uni.GetTranslator("en")
	if !found {
		return ErrTranslatorNotFound
	}

	// Register the English translations
	if err := en_translations.RegisterDefaultTranslations(Validator, Translator); err != nil {
		return err
	}

	return nil
}
