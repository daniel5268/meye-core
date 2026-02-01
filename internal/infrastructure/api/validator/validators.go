package validator

import (
	"meye-core/internal/domain/user"

	"github.com/go-playground/validator/v10"
)

func validateUserRoleEnum(fl validator.FieldLevel) bool {
	if val, ok := fl.Field().Interface().(user.UserRole); ok {
		return val == user.UserRoleAdmin || val == user.UserRoleMaster || val == user.UserRolePlayer
	}
	return false
}

func RegisterCustomValidators(v *validator.Validate) {
	v.RegisterValidation("userrole", validateUserRoleEnum)
}
