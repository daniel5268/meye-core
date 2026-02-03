package validator

import (
	"meye-core/internal/domain/campaign"
	"meye-core/internal/domain/user"

	"github.com/go-playground/validator/v10"
)

func validateUserRoleEnum(fl validator.FieldLevel) bool {
	if val, ok := fl.Field().Interface().(user.UserRole); ok {
		return val == user.UserRoleAdmin || val == user.UserRoleMaster || val == user.UserRolePlayer
	}
	return false
}

func validatePJTypeEnum(fl validator.FieldLevel) bool {
	if val, ok := fl.Field().Interface().(campaign.PJType); ok {
		return val == campaign.PJTypeHuman || val == campaign.PJTypeSupernatural
	}
	return false
}

func validateBasicTalentEnum(fl validator.FieldLevel) bool {
	if val, ok := fl.Field().Interface().(campaign.BasicTalentType); ok {
		return val == campaign.BasicTalentCoordination ||
			val == campaign.BasicTalentEnergy ||
			val == campaign.BasicTalentMental ||
			val == campaign.BasicTalentPhysical
	}
	return false
}

func validateSpecialTalent(fl validator.FieldLevel) bool {
	if val, ok := fl.Field().Interface().(campaign.SpecialTalentType); ok {
		return val == campaign.SpecialTalentEnergy ||
			val == campaign.SpecialTalentMental ||
			val == campaign.SpecialTalentPhysical
	}
	return false
}

func RegisterCustomValidators(v *validator.Validate) {
	v.RegisterValidation("userrole", validateUserRoleEnum)
	v.RegisterValidation("pjtype", validatePJTypeEnum)
	v.RegisterValidation("basictalent", validateBasicTalentEnum)
	v.RegisterValidation("specialtalent", validateSpecialTalent)
}
