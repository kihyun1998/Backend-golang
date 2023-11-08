package api

import (
	"simplebank/util"

	"github.com/go-playground/validator/v10"
)

var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	// FieldLevel.Field()는 reflect value여서 interface로 가져옴
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		//ok == true면 field가 문자열임
		return util.IsSupportedCurrency(currency)
	} else {
		//ok == false면 field가 문자열이 아님
		return false
	}
}
