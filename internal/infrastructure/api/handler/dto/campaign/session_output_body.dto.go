package campaign

import (
	"meye-core/internal/application/session"
	"time"
)

type SessionOutput struct {
	ID             string          `json:"id"`
	CampaignID     string          `json:"campaign_id"`
	Summary        string          `json:"summary"`
	XPAssignations []XPAssignation `json:"xp_assignations"`
	CreatedAt      time.Time       `json:"created_at"`
}

func MapSessionOutput(output session.SessionOutput) SessionOutput {
	xpAss := make([]XPAssignation, 0, len(output.XPAssignations))
	for _, xpA := range output.XPAssignations {
		xpAss = append(xpAss, XPAssignation{
			PjID: xpA.PjID,
			Amounts: XPAmounts{
				Basic:        xpA.Amounts.Basic,
				Special:      xpA.Amounts.Special,
				SuperNatural: xpA.Amounts.SuperNatural,
			},
			Reason: xpA.Reason,
		})
	}

	return SessionOutput{
		ID:             output.ID,
		CampaignID:     output.CampaignID,
		Summary:        output.Summary,
		XPAssignations: xpAss,
		CreatedAt:      output.CreatedAt,
	}
}
