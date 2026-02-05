package campaign

import "errors"

var (
	ErrUserNotInvited   = errors.New("ERR_USER_NOT_INVITED")
	ErrPjNotFound       = errors.New("ERR_PJ_NOT_FOUND")
	ErrPJsNotInCampaign = errors.New("ERR_PJS_NOT_IN_CAMPAIGN")
)
