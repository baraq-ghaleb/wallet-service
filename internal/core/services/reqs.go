package services

import (
	"context"
	"errors"
	"net/url"

	"github.com/ethereum/go-ethereum/common"
	auth "github.com/iden3/go-iden3-auth"
	"github.com/iden3/go-iden3-auth/loaders"
	"github.com/iden3/go-iden3-auth/pubsignals"
	"github.com/iden3/go-iden3-auth/state"
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
	// r := &protocol.AuthorizationRequestMessage {}
	
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

// TODO: remove or update CreateAuthRequestRequest
func (a *authRequest) CreateAuthRequest(ctx context.Context, req *ports.CreateAuthRequestRequest) (protocol.AuthorizationRequestMessage, error) {
	err := a.guardCreateAuthRequestRequest(req);
	if err != nil {
		log.Warn(ctx, "validating create authRequest request", "req", req)
	}

	const CallBackUrl      = "http:localhost:8001/call-back"
	const VerifierIdentity = "1125GJqgw6YEsKFwj63GY87MMxPL9kwDKxPUiwMLNZ"

    request := auth.CreateAuthorizationRequest("10", VerifierIdentity, "")
	request.ID = "7f38a193-0918-4a48-9fac-36adfdb8b542"
	request.ThreadID = "7f38a193-0918-4a48-9fac-36adfdb8b542"

	return request, err
	// if claimID NOT exist in request select all claims and filter it.
	// claims, err = auth.findClaimForQuery(ctx, identifier, query)
	// if err != nil {
	// 	return claim, nil, err
	// }

	// var mtpProofRequest protocol.ZeroKnowledgeProofRequest
	// mtpProofRequest.ID = 1
	// mtpProofRequest.CircuitID = string(circuits.AtomicQuerySigV2CircuitID)
	// mtpProofRequest.Query = map[string]interface{}{
	// 	"allowedIssuers": []string{"*"},
	// 	"credentialSubject": map[string]interface{}{
	// 		"birthday": map[string]interface{}{
	// 			"$lt": 20000101,
	// 		},
	// 	},
	// 	"context": "https://raw.githubusercontent.com/iden3/claim-schema-vocab/main/schemas/json-ld/kyc-v3.json-ld",
	// 	"type":    "KYCAgeCredential",
	// }



	// authRequest, err := domain.FromAuthRequester()
	// if err != nil {
	// 	log.Error(ctx, "Can not obtain the claim from claimer", err)
	// 	return nil, err
	// }

	// claimResp, err := a.save(ctx, authRequest)
	// if err != nil {
	// 	log.Error(ctx, "Can not save the claim", err)
	// 	return nil, err
	// }
	// return claimResp, err
}

func (a *authRequest) VerifyAuthRequestResponse(ctx context.Context, authorizationRequestMessage *protocol.AuthorizationRequestMessage, authorizationResponseMessage *protocol.AuthorizationResponseMessage) (bool) {
	keyDIR := "/home/zakwan/wallet-service/pkg/credentials/circuits/authV2"
	// circuitsLoaderService := pkgloader.NewCircuits("/home/zakwan/wallet-service/pkg/credentials/circuits")

	// authV2Set, err := circuitsLoaderService.Load(circuits.AuthV2CircuitID)

	var verificationKeyloader = &loaders.FSKeyLoader{Dir: keyDIR}

	URL := "https://polygon-mumbai.g.alchemy.com/v2/jNN6BxHCdHmxeTFcHHz-6DG7VTqX1tPY"
	contractAddress := "0x134B1BE34911E39A8397ec6289782989729807a4"

	resolver := state.ETHResolver{
		RPCUrl:          URL,
		ContractAddress:  common.HexToAddress(contractAddress),
	}

	resolvers := map[string]pubsignals.StateResolver{
		"polygon:mumbai": resolver,
	}
	verifier := auth.NewVerifier(verificationKeyloader, loaders.DefaultSchemaLoader{IpfsURL: "ipfs.io"}, resolvers)		

	err := verifier.VerifyAuthResponse(ctx, *authorizationResponseMessage, *authorizationRequestMessage)
	if err != nil {
		return false
	}
	return true
}

func (c *authRequest) guardCreateAuthRequestRequest(req *ports.CreateAuthRequestRequest) error {
	if _, err := url.ParseRequestURI(req.Schema); err != nil {
		return ErrMalformedURL
	}
	return nil
}

func (c *authRequest) guardAuthorizationRequestMessage(authorizationRequestMessage *protocol.AuthorizationRequestMessage) error {
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
