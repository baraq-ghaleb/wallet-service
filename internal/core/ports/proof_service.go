package ports

import (
	"context"
	"fmt"
	"math/big"

	core "github.com/iden3/go-iden3-core"

	"github.com/lastingasset/wallet-service/internal/core/domain"
)

// Query represents structure for query to atomic circuit
type Query struct {
	CircuitID                string
	Challenge                *big.Int
	AllowedIssuers           string                 `json:"allowedIssuers"`
	Req                      map[string]interface{} `json:"req"`
	Context                  string                 `json:"context"`
	Type                     string                 `json:"type"`
	ClaimID                  string                 `json:"claimId"`
	SkipClaimRevocationCheck bool                   `json:"skipClaimRevocationCheck"`
}

// SchemaType returns the schema type
func (q *Query) SchemaType() string {
	return fmt.Sprintf("%s#%s", q.Context, q.Type)
}

// ProofService is the interface implemented by the ProofService service
type ProofService interface {
	PrepareInputs(ctx context.Context, identifier *core.DID, query Query) ([]byte, []*domain.Claim, error)
	GenerateAuthProof(ctx context.Context, identifier *core.DID, challenge *big.Int) (*domain.FullProof, error)
	GenerateAgeProof(ctx context.Context, identifier *core.DID, query Query) (*domain.FullProof, error)
}
