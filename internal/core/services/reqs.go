package services

import (
	"context"
	"errors"
	"net/url"

	"github.com/iden3/go-circuits"
	"github.com/iden3/iden3comm/protocol"
	"github.com/lastingasset/wallet-service/internal/core/domain"
	"github.com/lastingasset/wallet-service/internal/core/ports"
	"github.com/lastingasset/wallet-service/internal/db"
	"github.com/lastingasset/wallet-service/internal/log"
)

var (
	ErrAuthRequestNotFound = errors.New("authRequest not found") // ErrAuthRequestNotFound Cannot retrieve the given authRequest 	// ErrProcessSchema Cannot process schema
)

// AuthRequestCfg authRequest service configuration
type AuthRequestCfg struct {
	RHSEnabled bool // ReverseHash Enabled
	RHSUrl     string
	Host       string
}

type authRequest struct {
	cfg                     AuthRequestCfg
	icRepo                  ports.ReqsRepository
	schemaSrv               ports.SchemaService
	identitySrv             ports.IdentityService
	mtService               ports.MtService
	identityStateRepository ports.IdentityStateRepository
	storage                 *db.Storage
}

// NewAuthRequest creates a new authRequest service
func NewAuthRequest(repo ports.ReqsRepository, schemaSrv ports.SchemaService, idenSrv ports.IdentityService, mtService ports.MtService, identityStateRepository ports.IdentityStateRepository, storage *db.Storage, cfg AuthRequestCfg) ports.ReqsService {
	s := &authRequest{
		cfg: AuthRequestCfg{
			RHSEnabled: cfg.RHSEnabled,
			RHSUrl:     cfg.RHSUrl,
			Host:       cfg.Host,
		},
		icRepo:                  repo,
		schemaSrv:               schemaSrv,
		identitySrv:             idenSrv,
		mtService:               mtService,
		identityStateRepository: identityStateRepository,
		storage:                 storage,
	}
	return s
}

// CreateAuthRequest creates a new authRequest
// 1.- Creates document
// 2.- Signature proof
// 3.- MerkelTree proof
func (a *authRequest) CreateAuthRequest(ctx context.Context, req *ports.CreateAuthRequestRequest) (*domain.AuthRequest, error) {
	if err := a.guardCreateAuthRequestRequest(req); err != nil {
		log.Warn(ctx, "validating create authRequest request", "req", req)
		return nil, err
	}

	var mtpProofRequest protocol.ZeroKnowledgeProofRequest
	mtpProofRequest.ID = 1
	mtpProofRequest.CircuitID = string(circuits.AtomicQuerySigV2CircuitID)
	mtpProofRequest.Query = map[string]interface{}{
		"allowedIssuers": []string{"*"},
		"credentialSubject": map[string]interface{}{
			"birthday": map[string]interface{}{
				"$lt": 20000101,
			},
		},
		"context": "https://raw.githubusercontent.com/iden3/claim-schema-vocab/main/schemas/json-ld/kyc-v3.json-ld",
		"type":    "KYCAgeCredential",
	}

	authRequest, err := domain.FromAuthRequester()
	if err != nil {
		log.Error(ctx, "Can not obtain the claim from claimer", err)
		return nil, err
	}

	claimResp, err := a.save(ctx, authRequest)
	if err != nil {
		log.Error(ctx, "Can not save the claim", err)
		return nil, err
	}
	return claimResp, err
}

func (c *authRequest) guardCreateAuthRequestRequest(req *ports.CreateAuthRequestRequest) error {
	if _, err := url.ParseRequestURI(req.Schema); err != nil {
		return ErrMalformedURL
	}
	return nil
}

func (a *authRequest) save(ctx context.Context, authRequest *domain.AuthRequest) (*domain.AuthRequest, error) {
	id, err := a.icRepo.Save(ctx, a.storage.Pgx, authRequest)
	if err != nil {
		return nil, err
	}

	authRequest.ID = id

	return authRequest, nil
}
