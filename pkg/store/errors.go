package store

import "github.com/dportaluppi/customer-profiles-api/pkg"

var (
	ErrIDMissing                   = pkg.NewErrID("missing store id")
	ErrInvalid                     = pkg.NewErrInvalid("invalid store data")
	ErrNotFound                    = pkg.NewErrNotFound("store not found")
	ErrConflict                    = pkg.NewErrConflict("store conflict occurred")
	ErrInternalError               = pkg.NewErrInternalError("store internal error")
	ErrInvalidPaginationParameters = pkg.NewErrInvalid("invalid store pagination parameters")
)
