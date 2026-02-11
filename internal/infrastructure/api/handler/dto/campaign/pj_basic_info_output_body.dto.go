package campaign

import "meye-core/internal/application/campaign"

type PjBasicInfoOutputBody struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func MapPjBasicInfoOutputBody(c campaign.PjBasicInfoOutput) PjBasicInfoOutputBody {
	return PjBasicInfoOutputBody{
		ID:   c.ID,
		Name: c.Name,
	}
}
