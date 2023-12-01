package profile

import "github.com/dportaluppi/customer-profiles-api/pkg"

var (
	ErrProfileIDMissing            = pkg.NewErrID("missing profile id")
	ErrProfileInvalid              = pkg.NewErrInvalid("invalid profile data")
	ErrNotFound                    = pkg.NewErrNotFound("profile not found")
	ErrProfileConflict             = pkg.NewErrConflict("profile conflict occurred")
	ErrInternalError               = pkg.NewErrInternalError("profile internal error")
	ErrInvalidPaginationParameters = pkg.NewErrInvalid("invalid profile pagination parameters")
)
