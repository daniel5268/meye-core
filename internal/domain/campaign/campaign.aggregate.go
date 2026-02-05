package campaign

import (
	"meye-core/internal/domain/shared"
)

type Campaign struct {
	id          string
	masterID    string
	name        string
	invitations []Invitation
	pjs         []PJ
}

func NewCampaign(masterID, name string, identificationService shared.IdentificationService) *Campaign {
	id := identificationService.GenerateID()

	return &Campaign{
		id:       id,
		masterID: masterID,
		name:     name,
	}
}

func (c *Campaign) InviteUser(userID string, identificationService shared.IdentificationService) (*Invitation, error) {
	invitation := NewInvitation(c.id, userID, identificationService)
	c.invitations = append(c.invitations, *invitation)

	return invitation, nil
}

func (c *Campaign) GetPendingUserInvitation(userID string) *Invitation {
	for i := range c.invitations {
		if c.invitations[i].UserID() == userID && c.invitations[i].State() == InvitationStatePending {
			return &c.invitations[i]
		}
	}

	return nil
}

type PJCreateParameters struct {
	Name                     string
	Weight                   uint
	Height                   uint
	Age                      uint
	Look                     uint
	Charisma                 int
	Villainy                 uint
	Heroism                  uint
	PjType                   PJType
	IsPhysicalTalented       bool
	IsMentalTalented         bool
	IsCoordinationTalented   bool
	IsPhysicalSkillsTalented bool
	IsMentalSkillsTalented   bool
	IsEnergySkillsTalented   bool
	IsEnergyTalented         bool
}

func (c *Campaign) AddPJ(userID string, params PJCreateParameters, identificationService shared.IdentificationService) (*PJ, error) {
	inv := c.GetPendingUserInvitation(userID)
	if inv == nil {
		return nil, ErrUserNotInvited
	}

	var supernaturalStats *SupernaturalStats
	if params.PjType == PJTypeSupernatural {
		supernaturalStats = &SupernaturalStats{
			skills: []Skill{
				{
					transformations: []uint{0},
				},
			},
		}
	}

	inv.accept()

	pj := &PJ{
		id:                identificationService.GenerateID(),
		userID:            userID,
		name:              params.Name,
		weight:            params.Weight,
		height:            params.Height,
		age:               params.Age,
		look:              params.Look,
		charisma:          params.Charisma,
		villainy:          params.Villainy,
		heroism:           params.Heroism,
		pjType:            params.PjType,
		supernaturalStats: supernaturalStats,
	}

	pj.basicStats.physical.isTalented = params.IsPhysicalTalented
	pj.basicStats.mental.isTalented = params.IsMentalTalented
	pj.basicStats.coordination.isTalented = params.IsCoordinationTalented

	pj.specialStats.physical.isTalented = params.IsPhysicalSkillsTalented
	pj.specialStats.energy.isTalented = params.IsEnergySkillsTalented
	pj.specialStats.mental.isTalented = params.IsMentalSkillsTalented

	pj.specialStats.isEnergyTalented = params.IsEnergyTalented

	pj.xp = CreateXPWithoutValidation(0, 0, 0)

	c.pjs = append(c.pjs, *pj)

	return pj, nil
}

func (c *Campaign) ID() string                { return c.id }
func (c *Campaign) MasterID() string          { return c.masterID }
func (c *Campaign) Name() string              { return c.name }
func (c *Campaign) Invitations() []Invitation { return c.invitations }
func (c *Campaign) PJs() []PJ                 { return c.pjs }

func CreateCampaignWithoutValidation(id, masterID, name string, invitations []Invitation, pjs []PJ) *Campaign {
	return &Campaign{
		id:          id,
		masterID:    masterID,
		name:        name,
		invitations: invitations,
		pjs:         pjs,
	}
}

func (c *Campaign) FindPjByID(pjID string) *PJ {
	for i := range c.pjs {
		if c.pjs[i].id == pjID {
			return &c.pjs[i]
		}
	}

	return nil
}

func (c *Campaign) MustContainPjs(pjIDs []string) error {
	campaignPJs := make(map[string]struct{}, len(c.pjs))
	for _, pj := range c.pjs {
		campaignPJs[pj.id] = struct{}{}
	}

	for _, pjID := range pjIDs {
		if _, exists := campaignPJs[pjID]; !exists {
			return ErrPJsNotInCampaign
		}
	}

	return nil
}
