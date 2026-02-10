package postgres

import (
	"database/sql/driver"
	"encoding/json"
	"meye-core/internal/domain/session"
	"time"
)

type XPAssignationJSON struct {
	PjID    string        `json:"pj_id"`
	Amounts XPAmountsJSON `json:"amounts"`
	Reason  string        `json:"reason"`
}

type XPAmountsJSON struct {
	Basic        uint `json:"basic"`
	Special      uint `json:"special"`
	Supernatural uint `json:"supernatural"`
}

type XPAssignationsJSON []XPAssignationJSON

func (x XPAssignationsJSON) Value() (driver.Value, error) {
	return json.Marshal(x)
}

func (x *XPAssignationsJSON) Scan(value interface{}) error {
	if value == nil {
		*x = XPAssignationsJSON{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, x)
}

type Session struct {
	ID             string             `gorm:"primaryKey"`
	CampaignID     string             `gorm:"column:campaign_id"`
	Summary        string             `gorm:"column:summary"`
	XPAssignations XPAssignationsJSON `gorm:"type:jsonb;column:xp_assignations"`
	CreatedAt      time.Time          `gorm:"column:created_at"`
}

func (s *Session) TableName() string {
	return "sessions"
}

func GetModelFromDomainSession(s *session.Session) *Session {
	xpAssignations := s.XPAssignations()
	assignationsJSON := make(XPAssignationsJSON, len(xpAssignations))

	for i, assignation := range xpAssignations {
		assignationsJSON[i] = XPAssignationJSON{
			PjID: assignation.PjID(),
			Amounts: XPAmountsJSON{
				Basic:        assignation.Amounts().Basic(),
				Special:      assignation.Amounts().Special(),
				Supernatural: assignation.Amounts().SuperNatural(),
			},
			Reason: assignation.Reason(),
		}
	}

	return &Session{
		ID:             s.ID(),
		CampaignID:     s.CampaignID(),
		Summary:        s.Summary(),
		XPAssignations: assignationsJSON,
		CreatedAt:      s.CreatedAt(),
	}
}

func (s *Session) ToDomain() *session.Session {
	xpAssignations := make([]session.XPAssignation, len(s.XPAssignations))

	for i, assignationJSON := range s.XPAssignations {
		xpAssignations[i] = session.NewXPAssignation(
			assignationJSON.PjID,
			assignationJSON.Amounts.Basic,
			assignationJSON.Amounts.Special,
			assignationJSON.Amounts.Supernatural,
			assignationJSON.Reason,
		)
	}

	return session.CreateSessionWithoutValidation(
		s.ID,
		s.CampaignID,
		s.Summary,
		s.CreatedAt,
		xpAssignations,
	)
}
