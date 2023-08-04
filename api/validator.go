package api

import (
	"github.com/Gokul-B12/money-txn/util"
	"github.com/go-playground/validator/v10"
)

var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		return util.IsValidCurrency(currency)
	}
	return false
}
