package campaign

import "meye-core/internal/domain/campaign"

type CreateCampaignInput struct {
	Name     string
	MasterID string
}

func MapCampaignOutput(c *campaign.Campaign) CampaignOutput {
	return CampaignOutput{
		ID:       c.ID(),
		Name:     c.Name(),
		MasterID: c.MasterID(),
	}
}

type CampaignOutput struct {
	ID       string
	Name     string
	MasterID string
}

type InviteUserInput struct {
	CampaignID string
	UserID     string
}

func MapInvitationOutput(i *campaign.Invitation) InvitationOutput {
	return InvitationOutput{
		ID:         i.ID(),
		CampaignID: i.CampaignID(),
		UserID:     i.UserID(),
		State:      i.State(),
	}
}

type InvitationOutput struct {
	ID         string
	CampaignID string
	UserID     string
	State      campaign.InvitationState
}

type UserCampaignIDs struct {
	UserID     string
	CampaignID string
}

type CreatePJInfo struct {
	Name          string
	Weight        uint
	Height        uint
	Age           uint
	Look          uint
	Charisma      int
	Villainy      uint
	Heroism       uint
	PjType        campaign.PJType
	BasicTalent   campaign.BasicTalentType
	SpecialTalent campaign.SpecialTalentType
}

type CreatePJInput struct {
	IDs    UserCampaignIDs
	PJInfo CreatePJInfo
}

type Physical struct {
	Strength   uint
	Agility    uint
	Speed      uint
	Resistance uint
}

type Coordination struct {
	Precision   uint
	Calculation uint
	Range       uint
	Reflexes    uint
}

type Mental struct {
	Inteligence   uint
	Wisdom        uint
	Concentration uint
	Will          uint
}

type BasicStats struct {
	Physical     Physical
	Mental       Mental
	Coordination Coordination
	Life         uint
}

type PhysicalSkills struct {
	Empowerment  uint
	VitalControl uint
}

type MentalSkills struct {
	Ilusion       uint
	MentalControl uint
}

type EnergySkills struct {
	ObjectHandling uint
	EnergyHandling uint
}

type SpecialStats struct {
	Physical   PhysicalSkills
	Mental     MentalSkills
	Energy     EnergySkills
	EnergyTank uint
}
type Skill struct {
	Transformations []uint
}

type SupernaturalStats struct {
	Skills []Skill
}

type PJOutput struct {
	ID                string
	UserID            string
	Name              string
	Weight            uint
	Height            uint
	Age               uint
	Look              uint
	Charisma          int
	Villainy          uint
	Heroism           uint
	PJType            string
	BasicTalent       string
	SpecialTalent     string
	BasicStats        BasicStats
	SpecialStats      SpecialStats
	SupernaturalStats *SupernaturalStats
}

func MapPJOutput(pj *campaign.PJ) PJOutput {
	output := PJOutput{
		ID:            pj.ID(),
		UserID:        pj.UserID(),
		Name:          pj.Name(),
		Weight:        pj.Weight(),
		Height:        pj.Height(),
		Age:           pj.Age(),
		Look:          pj.Look(),
		Charisma:      pj.Charisma(),
		Villainy:      pj.Villainy(),
		Heroism:       pj.Heroism(),
		PJType:        string(pj.Type()),
		BasicTalent:   string(pj.BasicTalent()),
		SpecialTalent: string(pj.SpecialTalent()),
		BasicStats: BasicStats{
			Physical: Physical{
				Strength:   pj.BasicStats().Physical().Strength(),
				Agility:    pj.BasicStats().Physical().Agility(),
				Speed:      pj.BasicStats().Physical().Speed(),
				Resistance: pj.BasicStats().Physical().Resistance(),
			},
			Mental: Mental{
				Inteligence:   pj.BasicStats().Mental().Inteligence(),
				Wisdom:        pj.BasicStats().Mental().Wisdom(),
				Concentration: pj.BasicStats().Mental().Concentration(),
				Will:          pj.BasicStats().Mental().Will(),
			},
			Coordination: Coordination{
				Precision:   pj.BasicStats().Coordination().Precision(),
				Calculation: pj.BasicStats().Coordination().Calculation(),
				Range:       pj.BasicStats().Coordination().Range(),
				Reflexes:    pj.BasicStats().Coordination().Reflexes(),
			},
			Life: pj.BasicStats().Life(),
		},
		SpecialStats: SpecialStats{
			Physical: PhysicalSkills{
				Empowerment:  pj.SpecialStats().Physical().Empowerment(),
				VitalControl: pj.SpecialStats().Physical().VitalControl(),
			},
			Mental: MentalSkills{
				Ilusion:       pj.SpecialStats().Mental().Ilusion(),
				MentalControl: pj.SpecialStats().Mental().MentalControl(),
			},
			Energy: EnergySkills{
				ObjectHandling: pj.SpecialStats().Energy().ObjectHandling(),
				EnergyHandling: pj.SpecialStats().Energy().EnergyHandling(),
			},
			EnergyTank: pj.SpecialStats().EnergyTank(),
		},
	}

	if pj.SupernaturalStats() != nil {
		skills := make([]Skill, len(pj.SupernaturalStats().Skills()))
		for i, skill := range pj.SupernaturalStats().Skills() {
			skills[i] = Skill{
				Transformations: skill.Transformations(),
			}
		}
		output.SupernaturalStats = &SupernaturalStats{
			Skills: skills,
		}
	}

	return output
}
