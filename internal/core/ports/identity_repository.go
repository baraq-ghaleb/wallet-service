package ports

import (
	"context"

	core "github.com/iden3/go-iden3-core"

	"github.com/lastingasset/wallet-service/internal/core/domain"
	"github.com/lastingasset/wallet-service/internal/db"
)

// IndentityRepository is the interface implemented by the identity service
type IndentityRepository interface {
	Save(ctx context.Context, conn db.Querier, identity *domain.Identity) error
	GetByID(ctx context.Context, conn db.Querier, identifier *core.DID) (*domain.Identity, error)
	Get(ctx context.Context, conn db.Querier) (identities []string, err error)
	GetUnprocessedIssuersIDs(ctx context.Context, conn db.Querier) (issuersIDs []*core.DID, err error)
	HasUnprocessedStatesByID(ctx context.Context, conn db.Querier, identifier *core.DID) (bool, error)
}
