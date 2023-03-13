package ports

import (
	"context"

	"github.com/google/uuid"

	"github.com/polygonid/sh-id-platform/internal/core/domain"
	"github.com/polygonid/sh-id-platform/internal/db"
)

// AuthRequestsRepository is the interface that defines the available methods
type ReqsRepository interface {
	Save(ctx context.Context, conn db.Querier, authRequest *domain.AuthRequest) (uuid.UUID, error)
}
