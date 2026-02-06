package postgres

import (
	"database/sql/driver"
	"encoding/json"
	"meye-core/internal/domain/campaign"
	"time"
)

type SupernaturalStatsJSON struct {
	Skills []SkillJSON `json:"skills"`
}

type SkillJSON struct {
	Transformations []uint `json:"transformations"`
}

func (s SupernaturalStatsJSON) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *SupernaturalStatsJSON) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, s)
}

type PJ struct {
	ID         string `gorm:"primaryKey"`
	CampaignID string
	UserID     string
	Name       string
	Weight     uint
	Height     uint
	Age        uint
	Look       uint
	Charisma   int
	Villainy   uint
	Heroism    uint
	PjType     campaign.PJType `gorm:"column:pj_type"`

	// Basic Stats - Physical
	Strength           uint
	Agility            uint
	Speed              uint
	Resistance         uint
	IsPhysicalTalented bool `gorm:"column:is_physical_talented"`

	// Basic Stats - Mental
	Inteligence      uint
	Wisdom           uint
	Concentration    uint
	Will             uint
	IsMentalTalented bool `gorm:"column:is_mental_talented"`

	// Basic Stats - Coordination
	Precision              uint
	Calculation            uint
	Range                  uint
	Reflexes               uint
	IsCoordinationTalented bool `gorm:"column:is_coordination_talented"`

	// Basic Stats - Life
	Life uint

	// Special Stats - Physical
	Empowerment              uint
	VitalControl             uint `gorm:"column:vital_control"`
	IsPhysicalSkillsTalented bool `gorm:"column:is_physical_skills_talented"`

	// Special Stats - Mental
	Ilusion                uint
	MentalControl          uint `gorm:"column:mental_control"`
	IsMentalSkillsTalented bool `gorm:"column:is_mental_skills_talented"`

	// Special Stats - Energy
	ObjectHandling         uint                   `gorm:"column:object_handling"`
	EnergyHandling         uint                   `gorm:"column:energy_handling"`
	EnergyTank             uint                   `gorm:"column:energy_tank"`
	IsEnergySkillsTalented bool                   `gorm:"column:is_energy_skills_talented"`
	IsEnergyTalented       bool                   `gorm:"column:is_energy_talented"`
	SupernaturalStats      *SupernaturalStatsJSON `gorm:"type:jsonb"`

	XPBasic        uint `gorm:"column:xp_basic"`
	XPSpecial      uint `gorm:"column:xp_special"`
	XPSupernatural uint `gorm:"column:xp_supernatural"`

	CreatedAt time.Time `gorm:"default:current_timestamp"`
	UpdatedAt time.Time `gorm:"default:current_timestamp"`
}

func GetModelFromDomainPJ(pj *campaign.PJ) *PJ {
	model := &PJ{
		ID:         pj.ID(),
		CampaignID: pj.CampaignID(),
		UserID:     pj.UserID(),
		Name:       pj.Name(),
		Weight:     pj.Weight(),
		Height:     pj.Height(),
		Age:        pj.Age(),
		Look:       pj.Look(),
		Charisma:   pj.Charisma(),
		Villainy:   pj.Villainy(),
		Heroism:    pj.Heroism(),
		PjType:     pj.Type(),

		// Basic Stats - Physical
		Strength:           pj.BasicStats().Physical().Strength(),
		Agility:            pj.BasicStats().Physical().Agility(),
		Speed:              pj.BasicStats().Physical().Speed(),
		Resistance:         pj.BasicStats().Physical().Resistance(),
		IsPhysicalTalented: pj.BasicStats().Physical().IsTalented(),

		// Basic Stats - Mental
		Inteligence:      pj.BasicStats().Mental().Inteligence(),
		Wisdom:           pj.BasicStats().Mental().Wisdom(),
		Concentration:    pj.BasicStats().Mental().Concentration(),
		Will:             pj.BasicStats().Mental().Will(),
		IsMentalTalented: pj.BasicStats().Mental().IsTalented(),

		// Basic Stats - Coordination
		Precision:              pj.BasicStats().Coordination().Precision(),
		Calculation:            pj.BasicStats().Coordination().Calculation(),
		Range:                  pj.BasicStats().Coordination().Range(),
		Reflexes:               pj.BasicStats().Coordination().Reflexes(),
		IsCoordinationTalented: pj.BasicStats().Coordination().IsTalented(),

		// Basic Stats - Life
		Life: pj.BasicStats().Life(),

		// Special Stats - Physical
		Empowerment:              pj.SpecialStats().Physical().Empowerment(),
		VitalControl:             pj.SpecialStats().Physical().VitalControl(),
		IsPhysicalSkillsTalented: pj.SpecialStats().Physical().IsTalented(),

		// Special Stats - Mental
		Ilusion:                pj.SpecialStats().Mental().Ilusion(),
		MentalControl:          pj.SpecialStats().Mental().MentalControl(),
		IsMentalSkillsTalented: pj.SpecialStats().Mental().IsTalented(),

		// Special Stats - Energy
		ObjectHandling:         pj.SpecialStats().Energy().ObjectHandling(),
		EnergyHandling:         pj.SpecialStats().Energy().EnergyHandling(),
		EnergyTank:             pj.SpecialStats().EnergyTank(),
		IsEnergySkillsTalented: pj.SpecialStats().Energy().IsTalented(),
		IsEnergyTalented:       pj.SpecialStats().IsEnergyTalented(),

		XPBasic:        pj.XP().Basic(),
		XPSpecial:      pj.XP().Special(),
		XPSupernatural: pj.XP().Supernatural(),
	}

	// Handle supernatural stats if present
	if pj.SupernaturalStats() != nil {
		skills := pj.SupernaturalStats().Skills()
		skillsJSON := make([]SkillJSON, len(skills))
		for i, skill := range skills {
			skillsJSON[i] = SkillJSON{
				Transformations: skill.Transformations(),
			}
		}
		model.SupernaturalStats = &SupernaturalStatsJSON{
			Skills: skillsJSON,
		}
	}

	return model
}

func (pj *PJ) ToDomain() *campaign.PJ {
	// Reconstruct Physical
	physical := campaign.CreatePhysicalWithoutValidation(
		pj.Strength,
		pj.Agility,
		pj.Speed,
		pj.Resistance,
		pj.IsPhysicalTalented,
	)

	// Reconstruct Mental
	mental := campaign.CreateMentalWithoutValidation(
		pj.Inteligence,
		pj.Wisdom,
		pj.Concentration,
		pj.Will,
		pj.IsMentalTalented,
	)

	// Reconstruct Coordination
	coordination := campaign.CreateCoordinationWithoutValidation(
		pj.Precision,
		pj.Calculation,
		pj.Range,
		pj.Reflexes,
		pj.IsCoordinationTalented,
	)

	// Reconstruct BasicStats
	basicStats := campaign.CreateBasicStatsWithoutValidation(
		physical,
		mental,
		coordination,
		pj.Life,
	)

	// Reconstruct PhysicalSkills
	physicalSkills := campaign.CreatePhysicalSkillsWithoutValidation(
		pj.Empowerment,
		pj.VitalControl,
		pj.IsPhysicalSkillsTalented,
	)

	// Reconstruct MentalSkills
	mentalSkills := campaign.CreateMentalSkillsWithoutValidation(
		pj.Ilusion,
		pj.MentalControl,
		pj.IsMentalSkillsTalented,
	)

	// Reconstruct EnergySkills
	energySkills := campaign.CreateEnergySkillsWithoutValidation(
		pj.ObjectHandling,
		pj.EnergyHandling,
		pj.IsEnergySkillsTalented,
	)

	// Reconstruct SpecialStats
	specialStats := campaign.CreateSpecialStatsWithoutValidation(
		physicalSkills,
		mentalSkills,
		energySkills,
		pj.EnergyTank,
		pj.IsEnergyTalented,
	)

	// Reconstruct SupernaturalStats if present
	var supernaturalStats *campaign.SupernaturalStats
	if pj.SupernaturalStats != nil {
		skills := make([]campaign.Skill, len(pj.SupernaturalStats.Skills))
		for i, skillJSON := range pj.SupernaturalStats.Skills {
			skills[i] = campaign.CreateSkillWithoutValidation(skillJSON.Transformations)
		}
		supernaturalStats = campaign.CreateSupernaturalStatsWithoutValidation(skills)
	}

	xp := campaign.CreateXPWithoutValidation(
		pj.XPBasic,
		pj.XPSpecial,
		pj.XPSupernatural,
	)

	return campaign.CreatePJWithoutValidation(
		pj.ID,
		pj.CampaignID,
		pj.UserID,
		pj.Name,
		pj.Weight,
		pj.Height,
		pj.Age,
		pj.Look,
		pj.Charisma,
		pj.Villainy,
		pj.Heroism,
		pj.PjType,
		basicStats,
		specialStats,
		supernaturalStats,
		xp,
	)
}
