package validation

import "github.com/go-playground/validator/v10"

var validate = validator.New()

type ErrorReponse struct {
  FailedField string
  Tag         string
  Value       string
}

func ValidateStruct(data interface{}) []*ErrorReponse {
  var errors []*ErrorReponse
  err := validate.Struct(data)
  if err != nil {
    for _, err := range err.(validator.ValidationErrors) {
      var element ErrorReponse
      element.FailedField = err.StructNamespace()
      element.Tag = err.Tag()
      element.Value = err.Param()
      errors = append(errors, &element)
    }
  }

  return errors
}
