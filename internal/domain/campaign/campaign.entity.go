package campaign

import "meye-core/internal/domain/shared"

type Campaign struct {
	id       string
	masterID string
	name     string
}

func NewCampaign(masterID, name string, identificationService shared.IdentificationService) *Campaign {
	id := identificationService.GenerateID()

	return &Campaign{
		id:       id,
		masterID: masterID,
		name:     name,
	}
}

func (c *Campaign) ID() string       { return c.id }
func (c *Campaign) MasterID() string { return c.masterID }
func (c *Campaign) Name() string     { return c.name }

func CreateCampaignWithoutValidation(id, masterID, name string) *Campaign {
	return &Campaign{
		id:       id,
		masterID: masterID,
		name:     name,
	}
}
