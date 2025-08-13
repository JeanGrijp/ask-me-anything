package validators

import (
	"github.com/go-playground/locales/pt_BR"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	ptBRTranslations "github.com/go-playground/validator/v10/translations/pt_BR"
)

var (
	Validate   *validator.Validate
	Translator ut.Translator
)

func InitValidator() error {
	Validate = validator.New()

	// Configura o tradutor para português brasileiro
	locale := pt_BR.New()
	uni := ut.New(locale, locale)

	trans, found := uni.GetTranslator("pt_BR")
	if !found {
		return nil
	}
	Translator = trans

	// Registra as traduções padrão
	if err := ptBRTranslations.RegisterDefaultTranslations(Validate, Translator); err != nil {
		return err
	}

	return nil
}
