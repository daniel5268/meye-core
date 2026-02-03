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
	ID            string `gorm:"primaryKey"`
	CampaignID    string
	UserID        string
	Name          string
	Weight        uint
	Height        uint
	Age           uint
	Look          uint
	Charisma      int
	Villainy      uint
	Heroism       uint
	PjType        campaign.PJType            `gorm:"column:pj_type"`
	BasicTalent   campaign.BasicTalentType   `gorm:"column:basic_talent"`
	SpecialTalent campaign.SpecialTalentType `gorm:"column:special_talent"`

	// Basic Stats - Physical
	Strength   uint
	Agility    uint
	Speed      uint
	Resistance uint

	// Basic Stats - Mental
	Inteligence   uint
	Wisdom        uint
	Concentration uint
	Will          uint

	// Basic Stats - Coordination
	Precision   uint
	Calculation uint
	Range       uint
	Reflexes    uint

	// Basic Stats - Life
	Life uint

	// Special Stats - Physical
	Empowerment  uint
	VitalControl uint `gorm:"column:vital_control"`

	// Special Stats - Mental
	Ilusion       uint
	MentalControl uint `gorm:"column:mental_control"`

	// Special Stats - Energy
	ObjectHandling uint `gorm:"column:object_handling"`
	EnergyHandling uint `gorm:"column:energy_handling"`
	EnergyTank     uint `gorm:"column:energy_tank"`

	// Supernatural Stats (nullable)
	SupernaturalStats *SupernaturalStatsJSON `gorm:"type:jsonb"`

	CreatedAt time.Time `gorm:"default:current_timestamp"`
	UpdatedAt time.Time `gorm:"default:current_timestamp"`
}

func GetModelFromDomainPJ(pj *campaign.PJ, campaignID string) *PJ {
	model := &PJ{
		ID:            pj.ID(),
		CampaignID:    campaignID,
		UserID:        pj.UserID(),
		Name:          pj.Name(),
		Weight:        pj.Weight(),
		Height:        pj.Height(),
		Age:           pj.Age(),
		Look:          pj.Look(),
		Charisma:      pj.Charisma(),
		Villainy:      pj.Villainy(),
		Heroism:       pj.Heroism(),
		PjType:        pj.Type(),
		BasicTalent:   pj.BasicTalent(),
		SpecialTalent: pj.SpecialTalent(),

		// Basic Stats - Physical
		Strength:   pj.BasicStats().Physical().Strength(),
		Agility:    pj.BasicStats().Physical().Agility(),
		Speed:      pj.BasicStats().Physical().Speed(),
		Resistance: pj.BasicStats().Physical().Resistance(),

		// Basic Stats - Mental
		Inteligence:   pj.BasicStats().Mental().Inteligence(),
		Wisdom:        pj.BasicStats().Mental().Wisdom(),
		Concentration: pj.BasicStats().Mental().Concentration(),
		Will:          pj.BasicStats().Mental().Will(),

		// Basic Stats - Coordination
		Precision:   pj.BasicStats().Coordination().Precision(),
		Calculation: pj.BasicStats().Coordination().Calculation(),
		Range:       pj.BasicStats().Coordination().Range(),
		Reflexes:    pj.BasicStats().Coordination().Reflexes(),

		// Basic Stats - Life
		Life: pj.BasicStats().Life(),

		// Special Stats - Physical
		Empowerment:  pj.SpecialStats().Physical().Empowerment(),
		VitalControl: pj.SpecialStats().Physical().VitalControl(),

		// Special Stats - Mental
		Ilusion:       pj.SpecialStats().Mental().Ilusion(),
		MentalControl: pj.SpecialStats().Mental().MentalControl(),

		// Special Stats - Energy
		ObjectHandling: pj.SpecialStats().Energy().ObjectHandling(),
		EnergyHandling: pj.SpecialStats().Energy().EnergyHandling(),
		EnergyTank:     pj.SpecialStats().EnergyTank(),
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
	)

	// Reconstruct Mental
	mental := campaign.CreateMentalWithoutValidation(
		pj.Inteligence,
		pj.Wisdom,
		pj.Concentration,
		pj.Will,
	)

	// Reconstruct Coordination
	coordination := campaign.CreateCoordinationWithoutValidation(
		pj.Precision,
		pj.Calculation,
		pj.Range,
		pj.Reflexes,
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
	)

	// Reconstruct MentalSkills
	mentalSkills := campaign.CreateMentalSkillsWithoutValidation(
		pj.Ilusion,
		pj.MentalControl,
	)

	// Reconstruct EnergySkills
	energySkills := campaign.CreateEnergySkillsWithoutValidation(
		pj.ObjectHandling,
		pj.EnergyHandling,
	)

	// Reconstruct SpecialStats
	specialStats := campaign.CreateSpecialStatsWithoutValidation(
		physicalSkills,
		mentalSkills,
		energySkills,
		pj.EnergyTank,
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

	return campaign.CreatePJWithoutValidation(
		pj.ID,
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
		pj.BasicTalent,
		pj.SpecialTalent,
		basicStats,
		specialStats,
		supernaturalStats,
	)
}
