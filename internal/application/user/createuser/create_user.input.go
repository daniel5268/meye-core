package createuser

import domainuser "meye-core/internal/domain/user"

type Input struct {
	Username string
	Password string
	Role     domainuser.UserRole
}
