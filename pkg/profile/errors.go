package profile

import "github.com/dportaluppi/customer-profiles-api/pkg"

var (
	ErrIDMissing                   = pkg.NewErrID("missing entity id")
	ErrAccountIDMissing            = pkg.NewErrID("missing account id")
	ErrInvalid                     = pkg.NewErrInvalid("invalid entity data")
	ErrNotFound                    = pkg.NewErrNotFound("entity not found")
	ErrConflict                    = pkg.NewErrConflict("entity conflict occurred")
	ErrInternalError               = pkg.NewErrInternalError("entity internal error")
	ErrInvalidPaginationParameters = pkg.NewErrInvalid("invalid entity pagination parameters")
)
