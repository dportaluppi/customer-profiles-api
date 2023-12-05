package order

import "github.com/dportaluppi/customer-profiles-api/pkg"

var (
	ErrIDMissing                   = pkg.NewErrID("missing profile id")
	ErrInvalid                     = pkg.NewErrInvalid("invalid profile data")
	ErrNotFound                    = pkg.NewErrNotFound("profile not found")
	ErrConflict                    = pkg.NewErrConflict("profile conflict occurred")
	ErrInternalError               = pkg.NewErrInternalError("profile internal error")
	ErrInvalidPaginationParameters = pkg.NewErrInvalid("invalid profile pagination parameters")
)
