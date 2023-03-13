package api_admin

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/lastingasset/wallet-service/internal/core/services"
	"github.com/lastingasset/wallet-service/internal/health"
	"github.com/lastingasset/wallet-service/internal/loader"
	"github.com/lastingasset/wallet-service/internal/repositories"
	"github.com/lastingasset/wallet-service/pkg/reverse_hash"
)

func TestServer_CheckStatus(t *testing.T) {
	identityRepo := repositories.NewIdentity()
	claimsRepo := repositories.NewClaims()
	reqsRepo := repositories.NewAuthRequests()
	identityStateRepo := repositories.NewIdentityState()
	mtRepo := repositories.NewIdentityMerkleTreeRepository()
	mtService := services.NewIdentityMerkleTrees(mtRepo)
	revocationRepository := repositories.NewRevocation()
	rhsp := reverse_hash.NewRhsPublisher(nil, false)
	identityService := services.NewIdentity(&KMSMock{}, identityRepo, mtRepo, identityStateRepo, mtService, claimsRepo, revocationRepository, storage, rhsp)
	schemaService := services.NewSchema(loader.CachedFactory(loader.HTTPFactory, cachex))

	claimsConf := services.ClaimCfg{
		RHSEnabled: false,
		Host:       "host",
	}
	authRequestsConf := services.AuthRequestCfg{
		RHSEnabled: false,
		Host:       "host",
	}
	claimsService := services.NewClaim(claimsRepo, schemaService, identityService, mtService, identityStateRepo, storage, claimsConf)
	reqsService := services.NewAuthRequest(reqsRepo, schemaService, identityService, mtService, identityStateRepo, storage, authRequestsConf)

	server := NewServer(&cfg, identityService, claimsService, reqsService, schemaService, NewPublisherMock(), NewPackageManagerMock(), &health.Status{})
	handler := getHandler(context.Background(), server)

	t.Run("should return 200", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/status", nil)
		require.NoError(t, err)

		handler.ServeHTTP(rr, req)
		require.Equal(t, http.StatusOK, rr.Code)
		var response Health200JSONResponse
		require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &response))
	})
}
