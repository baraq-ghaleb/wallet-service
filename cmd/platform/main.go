package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ethereum/go-ethereum/common"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	redis2 "github.com/go-redis/redis/v8"

	"github.com/lastingasset/wallet-service/internal/api"
	"github.com/lastingasset/wallet-service/internal/config"
	"github.com/lastingasset/wallet-service/internal/core/services"
	"github.com/lastingasset/wallet-service/internal/db"
	"github.com/lastingasset/wallet-service/internal/errors"
	"github.com/lastingasset/wallet-service/internal/gateways"
	"github.com/lastingasset/wallet-service/internal/health"
	"github.com/lastingasset/wallet-service/internal/kms"
	"github.com/lastingasset/wallet-service/internal/loader"
	"github.com/lastingasset/wallet-service/internal/log"
	"github.com/lastingasset/wallet-service/internal/providers"
	"github.com/lastingasset/wallet-service/internal/providers/blockchain"
	"github.com/lastingasset/wallet-service/internal/redis"
	"github.com/lastingasset/wallet-service/internal/repositories"
	"github.com/lastingasset/wallet-service/pkg/cache"
	"github.com/lastingasset/wallet-service/pkg/loaders"
	"github.com/lastingasset/wallet-service/pkg/protocol"
	"github.com/lastingasset/wallet-service/pkg/reverse_hash"
)

func main() {
	cfg, err := config.Load("")
	if err != nil {
		log.Error(context.Background(), "cannot load config", err)
		return
	}

	ctx := log.NewContext(context.Background(), cfg.Log.Level, cfg.Log.Mode, os.Stdout)

	if err := cfg.Sanitize(); err != nil {
		log.Error(ctx, "there are errors in the configuration that prevent server to start", err)
		return
	}

	storage, err := db.NewStorage(cfg.Database.URL)
	if err != nil {
		log.Error(ctx, "cannot connect to database", err)
		return
	}

	// Redis cache
	rdb, err := redis.Open(cfg.Cache.RedisUrl)
	if err != nil {
		log.Error(ctx, "cannot connect to redis", err, "host", cfg.Cache.RedisUrl)
		return
	}
	cachex := cache.NewRedisCache(rdb)
	schemaLoader := loader.CachedFactory(loader.HTTPFactory, cachex)

	vaultCli, err := providers.NewVaultClient(cfg.KeyStore.Address, cfg.KeyStore.Token)
	if err != nil {
		log.Error(ctx, "cannot init vault client: ", err)
		return
	}

	keyStore, err := kms.Open(cfg.KeyStore.PluginIden3MountPath, vaultCli)
	if err != nil {
		log.Error(ctx, "cannot initialize kms", err)
		return
	}

	ethereumClient, err := blockchain.Open(cfg)
	if err != nil {
		log.Error(ctx, "error dialing with ethereum client", err)
		return
	}

	stateContract, err := blockchain.InitEthClient(cfg.Ethereum.URL, cfg.Ethereum.ContractAddress)
	if err != nil {
		log.Error(ctx, "failed init ethereum client", err)
		return
	}

	ethConn, err := blockchain.InitEthConnect(cfg.Ethereum)
	if err != nil {
		log.Error(ctx, "failed init ethereum connect", err)
		return
	}

	circuitsLoaderService := loaders.NewCircuits(cfg.Circuit.Path)

	rhsp := reverse_hash.NewRhsPublisher(nil, false)

	// repositories initialization
	identityRepository := repositories.NewIdentity()
	claimsRepository := repositories.NewClaims()
	reqsRepository := repositories.NewAuthRequests()
	mtRepository := repositories.NewIdentityMerkleTreeRepository()
	identityStateRepository := repositories.NewIdentityState()
	revocationRepository := repositories.NewRevocation()

	// services initialization
	mtService := services.NewIdentityMerkleTrees(mtRepository)
	identityService := services.NewIdentity(keyStore, identityRepository, mtRepository, identityStateRepository, mtService, claimsRepository, revocationRepository, storage, rhsp)
	schemaService := services.NewSchema(schemaLoader)
	claimsService := services.NewClaim(
		claimsRepository,
		schemaService,
		identityService,
		mtService,
		identityStateRepository,
		storage,
		services.ClaimCfg{
			RHSEnabled: cfg.ReverseHashService.Enabled,
			RHSUrl:     cfg.ReverseHashService.URL,
			Host:       cfg.ServerUrl,
		},
	)
	reqsService := services.NewAuthRequest(
		reqsRepository,
		schemaService,
		identityService,
		mtService,
		identityStateRepository,
		storage,
		services.AuthRequestCfg{
			RHSEnabled: cfg.ReverseHashService.Enabled,
			RHSUrl:     cfg.ReverseHashService.URL,
			Host:       cfg.ServerUrl,
		},
	)
	proofService := gateways.NewProver(ctx, cfg, circuitsLoaderService)
	revocationService := services.NewRevocationService(ethConn, common.HexToAddress(cfg.Ethereum.ContractAddress))
	zkProofService := services.NewProofService(claimsService, revocationService, identityService, mtService, claimsRepository, proofService, keyStore, storage, stateContract, schemaLoader)
	transactionService, err := gateways.NewTransaction(ethereumClient, cfg.Ethereum.ConfirmationBlockCount)
	if err != nil {
		log.Error(ctx, "error creating transaction service", err)
		return
	}

	publisherGateway, err := gateways.NewPublisherEthGateway(ethereumClient, common.HexToAddress(cfg.Ethereum.ContractAddress), keyStore, cfg.PublishingKeyPath)
	if err != nil {
		log.Error(ctx, "error creating publish gateway", err)
		return
	}

	publisher := gateways.NewPublisher(storage, identityService, claimsService, reqsService, mtService, keyStore, transactionService, proofService, publisherGateway, cfg.Ethereum.ConfirmationTimeout)

	packageManager, err := protocol.InitPackageManager(ctx, stateContract, zkProofService, cfg.Circuit.Path)
	if err != nil {
		log.Error(ctx, "failed init package protocol", err)
		return
	}

	serverHealth := health.New(health.Monitors{
		"postgres": storage.Ping,
		"redis": func(rdb *redis2.Client) health.Pinger {
			return func(ctx context.Context) error { return rdb.Ping(ctx).Err() }
		}(rdb),
	})
	serverHealth.Run(ctx, health.DefaultPingPeriod)

	mux := chi.NewRouter()
	mux.Use(
		chiMiddleware.RequestID,
		log.ChiMiddleware(ctx),
		chiMiddleware.Recoverer,
		cors.Handler(cors.Options{AllowedOrigins: []string{"*"}}),
		chiMiddleware.NoCache,
	)
	api.HandlerFromMux(
		api.NewStrictHandlerWithOptions(
			api.NewServer(cfg, identityService, zkProofService, claimsService, reqsService, schemaService, publisher, packageManager, serverHealth),
			middlewares(ctx, cfg.HTTPBasicAuth),
			api.StrictHTTPServerOptions{
				RequestErrorHandlerFunc:  errors.RequestErrorHandlerFunc,
				ResponseErrorHandlerFunc: errors.ResponseErrorHandlerFunc,
			}),
		mux)
	api.RegisterStatic(mux)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.ServerPort),
		Handler: mux,
	}
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Info(ctx, "server started", "port", cfg.ServerPort)
		if err := server.ListenAndServe(); err != nil {
			log.Error(ctx, "Starting http server", err)
		}
	}()

	<-quit
	log.Info(ctx, "Shutting down")
}

func middlewares(ctx context.Context, auth config.HTTPBasicAuth) []api.StrictMiddlewareFunc {
	return []api.StrictMiddlewareFunc{
		api.LogMiddleware(ctx),
		api.BasicAuthMiddleware(ctx, auth.User, auth.Password),
	}
}
