package campaign

type XPAmounts struct {
	Basic        uint `json:"basic" binding:"min=0"`
	Special      uint `json:"special" binding:"min=0"`
	SuperNatural uint `json:"supernatural" binding:"min=0"`
}

type XPAssignation struct {
	PjID    string    `json:"pj_id" binding:"required,uuid"`
	Amounts XPAmounts `json:"amounts" binding:"required"`
	Reason  string    `json:"reason"`
}

type CreateSessionInputBody struct {
	Summary        string          `json:"summary"`
	XPAssignations []XPAssignation `json:"xp_assignations" binding:"omitempty,dive"`
}
