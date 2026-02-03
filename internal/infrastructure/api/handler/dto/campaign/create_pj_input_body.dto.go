package campaign

import domaincampaign "meye-core/internal/domain/campaign"

type CreatePJInputBody struct {
	Name          string                           `json:"name" binding:"required"`
	Weight        uint                             `json:"weight" binding:"gt=0"`
	Height        uint                             `json:"height" binding:"gt=0"`
	Age           uint                             `json:"age" binding:"gte=0"`
	Look          uint                             `json:"look" binding:"required,gt=0,lte=20"`
	Charisma      int                              `json:"charisma" binding:"required,gte=-10,lte=10"`
	Villainy      uint                             `json:"villainy" binding:"required,lte=10"`
	Heroism       uint                             `json:"heroism" binding:"required,lte=10"`
	PjType        domaincampaign.PJType            `json:"type" binding:"required,pjtype"`
	BasicTalent   domaincampaign.BasicTalentType   `json:"basic_talent" binding:"required,basictalent"`
	SpecialTalent domaincampaign.SpecialTalentType `json:"special_talent" binding:"required,specialtalent"`
}
