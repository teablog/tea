package validate

import (
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/pkg/errors"
	"gopkg.in/go-playground/validator.v9"
	zhTranslations "gopkg.in/go-playground/validator.v9/translations/zh"
	"reflect"
	"strings"
)

var (
	uni      *ut.UniversalTranslator
	trans    ut.Translator
	validate = validator.New()
)

func init() {
	// register Zh translations
	zht := zh.New()
	uni = ut.New(zht, zht)
	trans, _ = uni.GetTranslator("zh")
	err := zhTranslations.RegisterDefaultTranslations(validate, trans)
	if err != nil {
		panic(err)
	}
	// 注册tagNameFunc 用于获取validator 的json字段
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

func DoValidate(s interface{}) error {
	err := validate.Struct(s)
	if err != nil {
		if errNew, ok := err.(*validator.InvalidValidationError); ok {
			panic(errNew)
		}
		for _, e := range err.(validator.ValidationErrors) {
			return errors.New(e.Translate(trans))
		}
		return err
	}

	return err
}
