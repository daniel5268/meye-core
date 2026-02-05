package session

import (
	"meye-core/internal/domain/session"
	"time"
)

type XPAmounts struct {
	Basic        uint
	Special      uint
	SuperNatural uint
}

type XPAssignation struct {
	PjID    string
	Amounts XPAmounts
	Reason  string
}

type CreateSessionInput struct {
	CampaignID     string
	Summary        string
	XPAssignations []XPAssignation
}

type SessionOutput struct {
	ID             string
	CampaignID     string
	Summary        string
	XPAssignations []XPAssignation
	CreatedAt      time.Time
}

func MapSessionOutput(s *session.Session) SessionOutput {
	sessionXpAssignations := s.XPAssignations()
	xpAssignations := make([]XPAssignation, 0, len(sessionXpAssignations))

	for _, xpA := range sessionXpAssignations {
		xpAssignations = append(xpAssignations, XPAssignation{
			PjID: xpA.PjID(),
			Amounts: XPAmounts{
				Basic:        xpA.Basic(),
				Special:      xpA.Special(),
				SuperNatural: xpA.SuperNatural(),
			},
			Reason: xpA.Reason(),
		})
	}

	return SessionOutput{
		ID:             s.ID(),
		CampaignID:     s.CampaignID(),
		Summary:        s.Summary(),
		XPAssignations: xpAssignations,
		CreatedAt:      s.CreatedAt(),
	}
}
