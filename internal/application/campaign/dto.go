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
	Name                     string
	Weight                   uint
	Height                   uint
	Age                      uint
	Look                     uint
	Charisma                 int
	Villainy                 uint
	Heroism                  uint
	PjType                   campaign.PJType
	IsPhysicalTalented       bool
	IsMentalTalented         bool
	IsCoordinationTalented   bool
	IsPhysicalSkillsTalented bool
	IsMentalSkillsTalented   bool
	IsEnergySkillsTalented   bool
	IsEnergyTalented         bool
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
	IsTalented bool
}

type Coordination struct {
	Precision   uint
	Calculation uint
	Range       uint
	Reflexes    uint
	IsTalented  bool
}

type Mental struct {
	Inteligence   uint
	Wisdom        uint
	Concentration uint
	Will          uint
	IsTalented    bool
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
	IsTalented   bool
}

type MentalSkills struct {
	Ilusion       uint
	MentalControl uint
	IsTalented    bool
}

type EnergySkills struct {
	ObjectHandling uint
	EnergyHandling uint
	IsTalented     bool
}

type SpecialStats struct {
	Physical         PhysicalSkills
	Mental           MentalSkills
	Energy           EnergySkills
	EnergyTank       uint
	IsEnergyTalented bool
}
type Skill struct {
	Transformations []uint
}

type SupernaturalStats struct {
	Skills []Skill
}

type XP struct {
	Basic        uint
	Special      uint
	Supernatural uint
}

type PJOutput struct {
	ID                string
	CampaignID        string
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
	BasicStats        BasicStats
	SpecialStats      SpecialStats
	SupernaturalStats *SupernaturalStats
	XP                XP
}

func MapPJOutput(pj *campaign.PJ) PJOutput {
	output := PJOutput{
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
		PJType:     string(pj.Type()),
		BasicStats: BasicStats{
			Physical: Physical{
				Strength:   pj.BasicStats().Physical().Strength(),
				Agility:    pj.BasicStats().Physical().Agility(),
				Speed:      pj.BasicStats().Physical().Speed(),
				Resistance: pj.BasicStats().Physical().Resistance(),
				IsTalented: pj.BasicStats().Physical().IsTalented(),
			},
			Mental: Mental{
				Inteligence:   pj.BasicStats().Mental().Inteligence(),
				Wisdom:        pj.BasicStats().Mental().Wisdom(),
				Concentration: pj.BasicStats().Mental().Concentration(),
				Will:          pj.BasicStats().Mental().Will(),
				IsTalented:    pj.BasicStats().Mental().IsTalented(),
			},
			Coordination: Coordination{
				Precision:   pj.BasicStats().Coordination().Precision(),
				Calculation: pj.BasicStats().Coordination().Calculation(),
				Range:       pj.BasicStats().Coordination().Range(),
				Reflexes:    pj.BasicStats().Coordination().Reflexes(),
				IsTalented:  pj.BasicStats().Coordination().IsTalented(),
			},
			Life: pj.BasicStats().Life(),
		},
		SpecialStats: SpecialStats{
			Physical: PhysicalSkills{
				Empowerment:  pj.SpecialStats().Physical().Empowerment(),
				VitalControl: pj.SpecialStats().Physical().VitalControl(),
				IsTalented:   pj.SpecialStats().Physical().IsTalented(),
			},
			Mental: MentalSkills{
				Ilusion:       pj.SpecialStats().Mental().Ilusion(),
				MentalControl: pj.SpecialStats().Mental().MentalControl(),
				IsTalented:    pj.SpecialStats().Mental().IsTalented(),
			},
			Energy: EnergySkills{
				ObjectHandling: pj.SpecialStats().Energy().ObjectHandling(),
				EnergyHandling: pj.SpecialStats().Energy().EnergyHandling(),
				IsTalented:     pj.SpecialStats().Energy().IsTalented(),
			},
			EnergyTank:       pj.SpecialStats().EnergyTank(),
			IsEnergyTalented: pj.SpecialStats().IsEnergyTalented(),
		},
		XP: XP{
			Basic:        pj.XP().Basic(),
			Special:      pj.XP().Special(),
			Supernatural: pj.XP().Supernatural(),
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

type XpAmounts struct {
	Basic        uint
	Special      uint
	Supernatural uint
}

type ConsumeXpInput struct {
	PjID string
	Xp   XpAmounts
}

type UpdatePjStatsInput struct {
	PjID         string
	Basic        BasicStats
	Special      SpecialStats
	Supernatural *SupernaturalStats
}

func MapToPhysicalParameters(p Physical) campaign.PhysicalParameters {
	return campaign.PhysicalParameters{
		Strength:   p.Strength,
		Agility:    p.Agility,
		Speed:      p.Speed,
		Resistance: p.Resistance,
	}
}

func MapToMentalParameters(m Mental) campaign.MentalParameters {
	return campaign.MentalParameters{
		Intelligence:  m.Inteligence,
		Wisdom:        m.Wisdom,
		Concentration: m.Concentration,
		Will:          m.Will,
	}
}

func MapToCoordinationParameters(c Coordination) campaign.CoordinationParameters {
	return campaign.CoordinationParameters{
		Precision:   c.Precision,
		Calculation: c.Calculation,
		Range:       c.Range,
		Reflexes:    c.Reflexes,
	}
}

func MapToBasicStatsParameters(bs BasicStats) campaign.BasicStatsParameters {
	return campaign.BasicStatsParameters{
		Physical:     MapToPhysicalParameters(bs.Physical),
		Mental:       MapToMentalParameters(bs.Mental),
		Coordination: MapToCoordinationParameters(bs.Coordination),
		Life:         bs.Life,
	}
}

func MapToPhysicalSkillsParameters(ps PhysicalSkills) campaign.PhysicalSkillsParameters {
	return campaign.PhysicalSkillsParameters{
		Empowerment:  ps.Empowerment,
		VitalControl: ps.VitalControl,
	}
}

func MapToMentalSkillsParameters(ms MentalSkills) campaign.MentalSkillsParameters {
	return campaign.MentalSkillsParameters{
		Illusion:      ms.Ilusion,
		MentalControl: ms.MentalControl,
	}
}

func MapToEnergySkillsParameters(es EnergySkills) campaign.EnergySkillsParameters {
	return campaign.EnergySkillsParameters{
		ObjectHandling: es.ObjectHandling,
		EnergyHandling: es.EnergyHandling,
	}
}

func MapToSpecialStatsParameters(ss SpecialStats) campaign.SpecialStatsParameters {
	return campaign.SpecialStatsParameters{
		Physical:   MapToPhysicalSkillsParameters(ss.Physical),
		Mental:     MapToMentalSkillsParameters(ss.Mental),
		Energy:     MapToEnergySkillsParameters(ss.Energy),
		EnergyTank: ss.EnergyTank,
	}
}

func MapToSkillParameters(s Skill) campaign.SkillParameters {
	transformations := make([]uint, len(s.Transformations))
	copy(transformations, s.Transformations)
	return campaign.SkillParameters{
		Transformations: transformations,
	}
}

func MapToSupernaturalStatsParameters(ss SupernaturalStats) campaign.SupernaturalStatsParameters {
	skills := make([]campaign.SkillParameters, len(ss.Skills))
	for i, skill := range ss.Skills {
		skills[i] = MapToSkillParameters(skill)
	}
	return campaign.SupernaturalStatsParameters{
		Skills: skills,
	}
}

func MapToUpdatePjStatsParameters(input UpdatePjStatsInput) campaign.PjUpdateParameters {
	params := campaign.PjUpdateParameters{
		BasicStats:   MapToBasicStatsParameters(input.Basic),
		SpecialStats: MapToSpecialStatsParameters(input.Special),
	}

	// Only map supernatural stats if provided
	if input.Supernatural != nil {
		supernaturalParams := MapToSupernaturalStatsParameters(*input.Supernatural)
		params.SupernaturalStats = &supernaturalParams
	}

	return params
}
