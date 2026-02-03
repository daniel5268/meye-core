package campaign

import (
	"meye-core/internal/domain/shared"
	"meye-core/internal/domain/user"
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

func (c *Campaign) InviteUser(u *user.User, identificationService shared.IdentificationService) (*Invitation, error) {
	if !u.IsPlayer() {
		return nil, user.ErrUserNotPlayer
	}

	invitation := NewInvitation(c.id, u.ID(), identificationService)
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
		basicTalent:       params.BasicTalent,
		specialTalent:     params.SpecialTalent,
		supernaturalStats: supernaturalStats,
	}

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
