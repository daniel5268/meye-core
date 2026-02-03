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

func (c *Campaign) ID() string                { return c.id }
func (c *Campaign) MasterID() string          { return c.masterID }
func (c *Campaign) Name() string              { return c.name }
func (c *Campaign) Invitations() []Invitation { return c.invitations }

func CreateCampaignWithoutValidation(id, masterID, name string, invitations []Invitation) *Campaign {
	return &Campaign{
		id:          id,
		masterID:    masterID,
		name:        name,
		invitations: invitations,
	}
}
