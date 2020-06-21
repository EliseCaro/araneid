package service

import (
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator"
	translations "github.com/go-playground/validator/translations/zh"
	"reflect"
)

type DefaultBaseVerify struct {
	validate    *validator.Validate
	translation ut.Translator
}

/** 初始化验证器 **/
func (v *DefaultBaseVerify) Begin() *validator.Validate {
	v.translation, _ = ut.New(zh.New(), zh.New()).GetTranslator("zh")
	v.validate = validator.New()
	v.validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		label := field.Tag.Get("label")
		if label == "" {
			return field.Name
		}
		return label
	})
	_ = translations.RegisterDefaultTranslations(v.validate, v.translation)
	return v.validate
}

/** 翻译错误为中文 **/
func (v *DefaultBaseVerify) Translate(errs validator.ValidationErrors) string {
	var errorMessage string
	for _, e := range errs {
		errorMessage = e.Translate(v.translation)
		break
	}
	return errorMessage
}
