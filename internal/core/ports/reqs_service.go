package ports

import (
	"context"

	"time"

	core "github.com/iden3/go-iden3-core"
	"github.com/iden3/iden3comm/protocol"

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

// AuthRequestsService is the interface implemented by the authRequest service
type ReqsService interface {
	CreateAuthRequest(ctx context.Context, authRequestReq *CreateAuthRequestRequest) (protocol.AuthorizationRequestMessage, error)
	VerifyAuthRequestResponse(ctx context.Context, authorizationRequestMessage *protocol.AuthorizationRequestMessage, authorizationResponseMessage *protocol.AuthorizationResponseMessage) (bool)
}