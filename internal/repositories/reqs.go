package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/labstack/gommon/log"

	"github.com/polygonid/sh-id-platform/internal/core/domain"
	"github.com/polygonid/sh-id-platform/internal/core/ports"
	"github.com/polygonid/sh-id-platform/internal/db"
)

// ErrAuthRequestDuplication authRequest duplication error
var (
	ErrAuthRequestDuplication = errors.New("authRequest duplication error")
	// ErrAuthRequestDoesNotExist authRequest does not exist
	ErrAuthRequestDoesNotExist = errors.New("authRequest does not exist")
)

type authRequests struct{}

// NewAuthRequests returns a new authRequest repository
func NewAuthRequests() ports.ReqsRepository {
	return &authRequests{}
}

func (c *authRequests) Save(ctx context.Context, conn db.Querier, authRequest *domain.AuthRequest) (uuid.UUID, error) {
	var err error
	id := authRequest.ID
	log.Info(ctx, "Saving the auth request ............. ")
	if err == nil {
		return id, nil
	}
	return uuid.Nil, fmt.Errorf("error saving the authRequest: %w", err)
}
