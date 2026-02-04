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

	// Create basic stats with talent information
	physical := Physical{
		strength:   0,
		agility:    0,
		speed:      0,
		resistance: 0,
		isTalented: params.IsPhysicalTalented,
	}
	mental := Mental{
		inteligence:   0,
		wisdom:        0,
		concentration: 0,
		will:          0,
		isTalented:    params.IsMentalTalented,
	}
	coordination := Coordination{
		precision:   0,
		calculation: 0,
		coordRange:  0,
		reflexes:    0,
		isTalented:  params.IsCoordinationTalented,
	}
	basicStats := BasicStats{
		physical:     physical,
		mental:       mental,
		coordination: coordination,
		life:         0,
	}

	// Create special stats with talent information
	physicalSkills := PhysicalSkills{
		empowerment:  0,
		vitalControl: 0,
		isTalented:   params.IsPhysicalSkillsTalented,
	}
	mentalSkills := MentalSkills{
		ilusion:       0,
		mentalControl: 0,
		isTalented:    params.IsMentalSkillsTalented,
	}
	energySkills := EnergySkills{
		objectHandling: 0,
		energyHandling: 0,
		isTalented:     params.IsEnergySkillsTalented,
	}
	specialStats := SpecialStats{
		physical:         physicalSkills,
		mental:           mentalSkills,
		energy:           energySkills,
		energyTank:       0,
		isEnergyTalented: params.IsEnergyTalented,
	}

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
		basicStats:        basicStats,
		specialStats:      specialStats,
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
