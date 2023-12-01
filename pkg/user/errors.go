package user

import "github.com/dportaluppi/customer-profiles-api/pkg"

var (
	ErrProfileIDMissing            = pkg.NewErrID("missing user id")
	ErrProfileInvalid              = pkg.NewErrInvalid("invalid user data")
	ErrNotFound                    = pkg.NewErrNotFound("user not found")
	ErrProfileConflict             = pkg.NewErrConflict("user conflict occurred")
	ErrInternalError               = pkg.NewErrInternalError("user internal error")
	ErrInvalidPaginationParameters = pkg.NewErrInvalid("invalid user pagination parameters")
)
