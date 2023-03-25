package ports

import (
	"context"

	"time"

	core "github.com/iden3/go-iden3-core"
	"github.com/lastingasset/wallet-service/iden3comm/protocol"
)

// CreateAuthRequestRequest struct
type CreateAuthRequestRequest struct {
	DID                   *core.DID
	Schema                string
	CredentialSubject     map[string]any
	Expiration            *time.Time
	Type                  string
	Version               uint32
	SubjectPos            string
	MerklizedRootPosition string
}

// CreateAuthRequestRequest struct
type CreateQueryRequestRequest struct {
	DID                   *core.DID
	Schema                string
	CredentialSubject     map[string]any
	Expiration            *time.Time
	Type                  string
	Version               uint32
	SubjectPos            string
	MerklizedRootPosition string
}

type GenerateProofQuery struct {
	AllowedIssuers        [] string
	Context               string
	CredentialSubject     map[string]any
}

type GenerateProofScope struct {
	ID                    string
	CircuitId             string
	Query                 GenerateProofQuery
}

type GenerateProofBody struct {
	CallbackUrl           string
	Reason                string
	Message               string
	Scope                 [] GenerateProofScope
}

type GenerateProofRequest struct {
	DID                   *core.DID
	ID                    string
	Typ                   string
	Type                  string
	Thid                  string 
	Body                  GenerateProofBody
	From                  string
	To                    string
}

// NewCreateAuthRequestRequest returns a new authRequest object with the given parameters
func NewCreateAuthRequestRequest(did *core.DID, credentialSchema string, credentialSubject map[string]any, expiration *int64, typ string, cVersion *uint32, subjectPos *string, merklizedRootPosition *string) *CreateAuthRequestRequest {
	req := &CreateAuthRequestRequest{
		DID:               did,
		Schema:            credentialSchema,
		CredentialSubject: credentialSubject,
		Type:              typ,
	}
	if expiration != nil {
		t := time.Unix(*expiration, 0)
		req.Expiration = &t
	}
	if cVersion != nil {
		req.Version = *cVersion
	}
	if subjectPos != nil {
		req.SubjectPos = *subjectPos
	}
	if merklizedRootPosition != nil {
		req.MerklizedRootPosition = *merklizedRootPosition
	}
	return req
}

// NewCreateQueryRequestRequest returns a new queryRequest object with the given parameters
func NewCreateQueryRequestRequest(did *core.DID, credentialSchema string, credentialSubject map[string]any, expiration *int64, typ string, cVersion *uint32, subjectPos *string, merklizedRootPosition *string) *CreateQueryRequestRequest {
	req := &CreateQueryRequestRequest{
		DID:               did,
		Schema:            credentialSchema,
		CredentialSubject: credentialSubject,
		Type:              typ,
	}
	if expiration != nil {
		t := time.Unix(*expiration, 0)
		req.Expiration = &t
	}
	if cVersion != nil {
		req.Version = *cVersion
	}
	if subjectPos != nil {
		req.SubjectPos = *subjectPos
	}
	if merklizedRootPosition != nil {
		req.MerklizedRootPosition = *merklizedRootPosition
	}
	return req
}

func NewGenerateProofRequest(did *core.DID, id string, typ string, _type string, thid string) *GenerateProofRequest {
	req := &GenerateProofRequest{
		DID:                   did,
		ID:                    id,
		Typ:                   typ,
		Type:                  _type,
		Thid:                  thid,
		// Body                  GenerateProofBody
		// From                  string
		// To                    string
	}
	// if expiration != nil {
	// 	t := time.Unix(*expiration, 0)
	// 	req.Expiration = &t
	// }
	// if cVersion != nil {
	// 	req.Version = *cVersion
	// }
	// if subjectPos != nil {
	// 	req.SubjectPos = *subjectPos
	// }
	// if merklizedRootPosition != nil {
	// 	req.MerklizedRootPosition = *merklizedRootPosition
	// }
	return req
}

// AuthRequestsService is the interface implemented by the authRequest service
type ReqsService interface {
	CreateAuthRequest(ctx context.Context, authRequestReq *CreateAuthRequestRequest) (protocol.AuthorizationRequestMessage, error)
	VerifyAuthRequestResponse(ctx context.Context, authorizationRequestMessage *protocol.AuthorizationRequestMessage, authorizationResponseMessage *protocol.AuthorizationResponseMessage) bool
	CreateQueryRequest(ctx context.Context, authRequestReq *CreateQueryRequestRequest) (protocol.AuthorizationRequestMessage, error)
	CreateAuthorizationRequestMessage(ctx context.Context, generateProofRequest *CreateQueryRequestRequest) (protocol.AuthorizationRequestMessage, error)
	VerifyQueryRequestResponse(ctx context.Context, authorizationRequestMessage *protocol.AuthorizationRequestMessage, authorizationResponseMessage *protocol.AuthorizationResponseMessage) bool
}
