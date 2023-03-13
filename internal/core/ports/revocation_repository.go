package ports

import (
	"context"

	core "github.com/iden3/go-iden3-core"

	"github.com/lastingasset/wallet-service/internal/core/domain"
	"github.com/lastingasset/wallet-service/internal/db"
)

// RevocationRepository interface that defines the available methods
type RevocationRepository interface {
	UpdateStatus(ctx context.Context, conn db.Querier, did *core.DID) ([]*domain.Revocation, error)
}
