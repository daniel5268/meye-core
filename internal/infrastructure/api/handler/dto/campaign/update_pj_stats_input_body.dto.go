package campaign

import "meye-core/internal/application/campaign"

type UpdatePJStatsInputBody struct {
	BasicStats        BasicStatsInputBody         `json:"basic_stats" binding:"required"`
	SpecialStats      SpecialStatsInputBody       `json:"special_stats" binding:"required"`
	SupernaturalStats *SupernaturalStatsInputBody `json:"supernatural_stats,omitempty"`
}

type BasicStatsInputBody struct {
	Physical     PhysicalInputBody     `json:"physical" binding:"required"`
	Mental       MentalInputBody       `json:"mental" binding:"required"`
	Coordination CoordinationInputBody `json:"coordination" binding:"required"`
	Life         uint                  `json:"life"`
}

type PhysicalInputBody struct {
	Strength   uint `json:"strength"`
	Agility    uint `json:"agility"`
	Speed      uint `json:"speed"`
	Resistance uint `json:"resistance"`
}

type MentalInputBody struct {
	Intelligence  uint `json:"intelligence"`
	Wisdom        uint `json:"wisdom"`
	Concentration uint `json:"concentration"`
	Will          uint `json:"will"`
}

type CoordinationInputBody struct {
	Precision   uint `json:"precision"`
	Calculation uint `json:"calculation"`
	Range       uint `json:"range"`
	Reflexes    uint `json:"reflexes"`
}

type SpecialStatsInputBody struct {
	Physical   PhysicalSkillsInputBody `json:"physical" binding:"required"`
	Mental     MentalSkillsInputBody   `json:"mental" binding:"required"`
	Energy     EnergySkillsInputBody   `json:"energy" binding:"required"`
	EnergyTank uint                    `json:"energy_tank"`
}

type PhysicalSkillsInputBody struct {
	Empowerment  uint `json:"empowerment"`
	VitalControl uint `json:"vital_control"`
}

type MentalSkillsInputBody struct {
	Illusion      uint `json:"illusion"`
	MentalControl uint `json:"mental_control"`
}

type EnergySkillsInputBody struct {
	ObjectHandling uint `json:"object_handling"`
	EnergyHandling uint `json:"energy_handling"`
}

type SupernaturalStatsInputBody struct {
	Skills []SkillInputBody `json:"skills" binding:"required"`
}

type SkillInputBody struct {
	Transformations []uint `json:"transformations" binding:"required"`
}

func MapUpdatePJStatsInput(pathParams PJPathParams, body UpdatePJStatsInputBody) campaign.UpdatePjStatsInput {
	input := campaign.UpdatePjStatsInput{
		PjID: pathParams.PJID, // Only PJID needed
		Basic: campaign.BasicStats{
			Physical: campaign.Physical{
				Strength:   body.BasicStats.Physical.Strength,
				Agility:    body.BasicStats.Physical.Agility,
				Speed:      body.BasicStats.Physical.Speed,
				Resistance: body.BasicStats.Physical.Resistance,
			},
			Mental: campaign.Mental{
				Inteligence:   body.BasicStats.Mental.Intelligence,
				Wisdom:        body.BasicStats.Mental.Wisdom,
				Concentration: body.BasicStats.Mental.Concentration,
				Will:          body.BasicStats.Mental.Will,
			},
			Coordination: campaign.Coordination{
				Precision:   body.BasicStats.Coordination.Precision,
				Calculation: body.BasicStats.Coordination.Calculation,
				Range:       body.BasicStats.Coordination.Range,
				Reflexes:    body.BasicStats.Coordination.Reflexes,
			},
			Life: body.BasicStats.Life,
		},
		Special: campaign.SpecialStats{
			Physical: campaign.PhysicalSkills{
				Empowerment:  body.SpecialStats.Physical.Empowerment,
				VitalControl: body.SpecialStats.Physical.VitalControl,
			},
			Mental: campaign.MentalSkills{
				Ilusion:       body.SpecialStats.Mental.Illusion,
				MentalControl: body.SpecialStats.Mental.MentalControl,
			},
			Energy: campaign.EnergySkills{
				ObjectHandling: body.SpecialStats.Energy.ObjectHandling,
				EnergyHandling: body.SpecialStats.Energy.EnergyHandling,
			},
			EnergyTank: body.SpecialStats.EnergyTank,
		},
	}

	// Map supernatural stats if provided
	if body.SupernaturalStats != nil {
		skills := make([]campaign.Skill, len(body.SupernaturalStats.Skills))
		for i, skill := range body.SupernaturalStats.Skills {
			skills[i] = campaign.Skill{
				Transformations: skill.Transformations,
			}
		}
		input.Supernatural = &campaign.SupernaturalStats{
			Skills: skills,
		}
	}

	return input
}
