package campaign

import "meye-core/internal/application/campaign"

type PJOutputBody struct {
	ID                string                 `json:"id"`
	UserID            string                 `json:"user_id"`
	Name              string                 `json:"name"`
	Weight            uint                   `json:"weight"`
	Height            uint                   `json:"height"`
	Age               uint                   `json:"age"`
	Look              uint                   `json:"look"`
	Charisma          int                    `json:"charisma"`
	Villainy          uint                   `json:"villainy"`
	Heroism           uint                   `json:"heroism"`
	PJType            string                 `json:"pj_type"`
	BasicTalent       string                 `json:"basic_talent"`
	SpecialTalent     string                 `json:"special_talent"`
	BasicStats        BasicStatsBody         `json:"basic_stats"`
	SpecialStats      SpecialStatsBody       `json:"special_stats"`
	SupernaturalStats *SupernaturalStatsBody `json:"supernatural_stats,omitempty"`
}

type BasicStatsBody struct {
	Physical     PhysicalBody     `json:"physical"`
	Mental       MentalBody       `json:"mental"`
	Coordination CoordinationBody `json:"coordination"`
	Life         uint             `json:"life"`
}

type PhysicalBody struct {
	Strength   uint `json:"strength"`
	Agility    uint `json:"agility"`
	Speed      uint `json:"speed"`
	Resistance uint `json:"resistance"`
}

type MentalBody struct {
	Intelligence  uint `json:"intelligence"`
	Wisdom        uint `json:"wisdom"`
	Concentration uint `json:"concentration"`
	Will          uint `json:"will"`
}

type CoordinationBody struct {
	Precision   uint `json:"precision"`
	Calculation uint `json:"calculation"`
	Range       uint `json:"range"`
	Reflexes    uint `json:"reflexes"`
}

type SpecialStatsBody struct {
	Physical   PhysicalSkillsBody `json:"physical"`
	Mental     MentalSkillsBody   `json:"mental"`
	Energy     EnergySkillsBody   `json:"energy"`
	EnergyTank uint               `json:"energy_tank"`
}

type PhysicalSkillsBody struct {
	Empowerment  uint `json:"empowerment"`
	VitalControl uint `json:"vital_control"`
}

type MentalSkillsBody struct {
	Illusion      uint `json:"illusion"`
	MentalControl uint `json:"mental_control"`
}

type EnergySkillsBody struct {
	ObjectHandling uint `json:"object_handling"`
	EnergyHandling uint `json:"energy_handling"`
}

type SupernaturalStatsBody struct {
	Skills []SkillBody `json:"skills"`
}

type SkillBody struct {
	Transformations []uint `json:"transformations"`
}

func MapPJOutputBody(output campaign.PJOutput) PJOutputBody {
	body := PJOutputBody{
		ID:            output.ID,
		UserID:        output.UserID,
		Name:          output.Name,
		Weight:        output.Weight,
		Height:        output.Height,
		Age:           output.Age,
		Look:          output.Look,
		Charisma:      output.Charisma,
		Villainy:      output.Villainy,
		Heroism:       output.Heroism,
		PJType:        output.PJType,
		BasicTalent:   output.BasicTalent,
		SpecialTalent: output.SpecialTalent,

		BasicStats: BasicStatsBody{
			Physical: PhysicalBody{
				Strength:   output.BasicStats.Physical.Strength,
				Agility:    output.BasicStats.Physical.Agility,
				Speed:      output.BasicStats.Physical.Speed,
				Resistance: output.BasicStats.Physical.Resistance,
			},
			Mental: MentalBody{
				Intelligence:  output.BasicStats.Mental.Inteligence,
				Wisdom:        output.BasicStats.Mental.Wisdom,
				Concentration: output.BasicStats.Mental.Concentration,
				Will:          output.BasicStats.Mental.Will,
			},
			Coordination: CoordinationBody{
				Precision:   output.BasicStats.Coordination.Precision,
				Calculation: output.BasicStats.Coordination.Calculation,
				Range:       output.BasicStats.Coordination.Range,
				Reflexes:    output.BasicStats.Coordination.Reflexes,
			},
			Life: output.BasicStats.Life,
		},

		SpecialStats: SpecialStatsBody{
			Physical: PhysicalSkillsBody{
				Empowerment:  output.SpecialStats.Physical.Empowerment,
				VitalControl: output.SpecialStats.Physical.VitalControl,
			},
			Mental: MentalSkillsBody{
				Illusion:      output.SpecialStats.Mental.Ilusion,
				MentalControl: output.SpecialStats.Mental.MentalControl,
			},
			Energy: EnergySkillsBody{
				ObjectHandling: output.SpecialStats.Energy.ObjectHandling,
				EnergyHandling: output.SpecialStats.Energy.EnergyHandling,
			},
			EnergyTank: output.SpecialStats.EnergyTank,
		},
	}

	// SupernaturalStats (optional)
	if output.SupernaturalStats != nil {
		supernaturalStats := &SupernaturalStatsBody{
			Skills: make([]SkillBody, len(output.SupernaturalStats.Skills)),
		}
		for i, skill := range output.SupernaturalStats.Skills {
			supernaturalStats.Skills[i] = SkillBody{
				Transformations: skill.Transformations,
			}
		}
		body.SupernaturalStats = supernaturalStats
	}

	return body
}
