package postgres

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"meye-core/internal/domain/session"
	"time"
)

type Session struct {
	ID             string
	CampaignID     string
	Summary        string
	XPAssignations XPAssignations
	CreatedAt      time.Time
}

// Custom type para JSONB
type XPAssignations []XPAssignation

type XPAmounts struct {
	Basic        uint `json:"basic"`
	Special      uint `json:"special"`
	SuperNatural uint `json:"supernatural"`
}

type XPAssignation struct {
	PjID    string    `json:"pj_id"`
	Amounts XPAmounts `json:"amounts"`
	Reason  string    `json:"reason"`
}

// Implementar driver.Valuer para escribir a JSONB
func (x XPAssignations) Value() (driver.Value, error) {
	return json.Marshal(x)
}

// Implementar sql.Scanner para leer desde JSONB
func (x *XPAssignations) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan XPAssignations")
	}
	return json.Unmarshal(bytes, x)
}

// Conversión a dominio
func (s *Session) ToDomain() *session.Session {
	assignations := make([]session.XPAssignation, len(s.XPAssignations))
	for i, a := range s.XPAssignations {
		assignations[i] = session.NewXPAssignation(a.PjID, a.Amounts.Basic, a.Amounts.Special, a.Amounts.SuperNatural, a.Reason)
	}

	return session.CreateSessionWithoutValidation(
		s.ID,
		s.CampaignID,
		s.Summary,
		s.CreatedAt,
		assignations,
	)
}

// Conversión desde dominio
func GetModelFromDomainSession(s *session.Session) *Session {
	assignations := make(XPAssignations, len(s.XPAssignations()))
	for i, a := range s.XPAssignations() {
		assignations[i] = XPAssignation{
			PjID: a.PjID(),
			Amounts: XPAmounts{
				Basic:        a.Basic(),
				Special:      a.Special(),
				SuperNatural: a.SuperNatural(),
			},
			Reason: a.Reason(),
		}
	}

	return &Session{
		ID:             s.ID(),
		CampaignID:     s.CampaignID(),
		Summary:        s.Summary(),
		XPAssignations: assignations,
		CreatedAt:      s.CreatedAt(),
	}
}
