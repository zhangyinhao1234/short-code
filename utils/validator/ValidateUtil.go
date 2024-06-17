package validator

import (
	"errors"
	"fmt"
	"github.com/go-playground/locales/zh_Hans_CN"
	unTrans "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTrans "github.com/go-playground/validator/v10/translations/zh"
	"reflect"
	"strings"
)

var validate_ *validator.Validate
var trans_ unTrans.Translator

func init() {
	validate_ = validator.New()
	uni := unTrans.New(zh_Hans_CN.New())
	trans_, _ = uni.GetTranslator("zh_Hans_CN")
	err := zhTrans.RegisterDefaultTranslations(validate_, trans_)
	if err != nil {
		fmt.Println("err:", err)
	}
	//将验证法字段名 映射为中文名
	validate_.RegisterTagNameFunc(func(field reflect.StructField) string {
		label := field.Tag.Get("label")
		return label
	})
}

func Validate(data interface{}) error {
	err := validate_.Struct(data)
	var sbuild strings.Builder
	if err != nil {
		for _, v := range err.(validator.ValidationErrors) {
			sbuild.WriteString(v.Translate(trans_))
			sbuild.WriteString(";")
		}
		return errors.New(sbuild.String())
	}
	return nil
}
