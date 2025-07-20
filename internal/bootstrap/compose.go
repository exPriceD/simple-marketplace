package bootstrap

import (
	"context"
	"net/http"
	"time"

	"github.com/exPriceD/simple-marketplace/internal/application/usecase"
	httpdelivery "github.com/exPriceD/simple-marketplace/internal/delivery/http"
	"github.com/exPriceD/simple-marketplace/internal/delivery/http/handler"
	"github.com/exPriceD/simple-marketplace/internal/domain/shared"
	"github.com/exPriceD/simple-marketplace/internal/infrastructure/database/postgres"
	obslogger "github.com/exPriceD/simple-marketplace/internal/infrastructure/observability/logger"
	obsmetrics "github.com/exPriceD/simple-marketplace/internal/infrastructure/observability/metrics"
	listingrepo "github.com/exPriceD/simple-marketplace/internal/infrastructure/repository/listing"
	userrepo "github.com/exPriceD/simple-marketplace/internal/infrastructure/repository/user"
	"github.com/exPriceD/simple-marketplace/internal/infrastructure/security/hasher"
	"github.com/exPriceD/simple-marketplace/internal/infrastructure/security/token"
	"github.com/exPriceD/simple-marketplace/internal/infrastructure/system"
	"github.com/exPriceD/simple-marketplace/internal/infrastructure/system/config"
	"github.com/exPriceD/simple-marketplace/internal/platform/log"
	"github.com/exPriceD/simple-marketplace/internal/platform/metrics"
)

type App struct {
	Server  *http.Server
	DB      *postgres.DB
	Logger  log.Logger
	Metrics metrics.Metrics
}

func Build(ctx context.Context, cfg config.Config) (*App, error) {
	db, err := postgres.Connect(ctx, postgres.Config{
		DSN:             cfg.DSN,
		MaxConns:        10,
		MinConns:        1,
		MaxConnLifetime: time.Hour,
	})
	if err != nil {
		return nil, err
	}
	if err := postgres.RunMigrations(ctx, db); err != nil {
		db.Close()
		return nil, err
	}

	zl, err := obslogger.New(cfg.LogLevel)
	if err != nil {
		db.Close()
		return nil, err
	}
	mtr := obsmetrics.New()

	hasherSvc := hasher.NewArgon2idHasher(cfg.ArgonTime, cfg.ArgonMemoryMB*1024, cfg.ArgonThreads, cfg.ArgonKeyLen)
	tokenSvc := token.NewJWTService(cfg.JWTSecret, cfg.JWTTTL)
	clock := system.SystemClock{}
	idGen := shared.RandomIDGenerator{}

	userRepo := userrepo.NewUserRepositoryPG(db)
	listRepo := listingrepo.NewListingRepositoryPG(db)

	regUCBase := usecase.NewRegisterUser(userRepo, hasherSvc, clock, idGen)
	regUC := usecase.InstrumentRegisterUser(regUCBase, zl, mtr)

	loginUCBase := usecase.NewLoginUser(userRepo, hasherSvc, tokenSvc)
	loginUC := usecase.InstrumentLoginUser(loginUCBase, zl, mtr)

	createListingUCBase := usecase.NewCreateListing(listRepo, clock)
	createListingUC := usecase.InstrumentCreateListing(createListingUCBase, zl, mtr)

	listListingsUCBase := usecase.NewListListings(listRepo)
	listListingsUC := usecase.InstrumentListListings(listListingsUCBase, zl, mtr)

	authH := handler.NewAuthHandler(regUC, loginUC)
	listH := handler.NewListingHandler(createListingUC, listListingsUC)

	router := httpdelivery.NewRouter(authH, listH, tokenSvc, zl, mtr, mtr.Handler())

	srv := &http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: router,
	}

	return &App{
		Server:  srv,
		DB:      db,
		Logger:  zl,
		Metrics: mtr,
	}, nil
}

func (a *App) Shutdown(ctx context.Context) {
	_ = a.Server.Shutdown(ctx)
	a.DB.Close()
}
