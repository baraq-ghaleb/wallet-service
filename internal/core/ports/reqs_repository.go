package ports

import (
	"context"

	"github.com/google/uuid"

	"github.com/lastingasset/wallet-service/internal/core/domain"
	"github.com/lastingasset/wallet-service/internal/db"
)

// AuthRequestsRepository is the interface that defines the available methods
type ReqsRepository interface {
	Save(ctx context.Context, conn db.Querier, authRequest *domain.AuthRequest) (uuid.UUID, error)
}
