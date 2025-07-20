package bootstrap

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/exPriceD/simple-marketplace/internal/application/port"
	"github.com/exPriceD/simple-marketplace/internal/application/usecase"
	httpdelivery "github.com/exPriceD/simple-marketplace/internal/delivery/http"
	"github.com/exPriceD/simple-marketplace/internal/delivery/http/handler"
	"github.com/exPriceD/simple-marketplace/internal/domain/shared"
	"github.com/exPriceD/simple-marketplace/internal/infrastructure/database/postgres"
	"github.com/exPriceD/simple-marketplace/internal/infrastructure/repository/listing"
	"github.com/exPriceD/simple-marketplace/internal/infrastructure/repository/user"
	"github.com/exPriceD/simple-marketplace/internal/infrastructure/security/hasher"
	"github.com/exPriceD/simple-marketplace/internal/infrastructure/security/token"
	"github.com/exPriceD/simple-marketplace/internal/infrastructure/system"
	"github.com/exPriceD/simple-marketplace/internal/infrastructure/system/config"
)

type App struct {
	Server *http.Server
	DB     *postgres.DB
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

	// Infra services
	hasherSvc := hasher.NewArgon2idHasher(cfg.ArgonTime, cfg.ArgonMemoryMB*1024, cfg.ArgonThreads, cfg.ArgonKeyLen)
	tokenSvc := token.NewJWTService(cfg.JWTSecret, cfg.JWTTTL)
	clock := system.SystemClock{}
	idGen := shared.RandomIDGenerator{}

	// Repositories
	userRepo := userrepo.NewUserRepositoryPG(db)
	listRepo := listingrepo.NewListingRepositoryPG(db)

	// Use Cases
	regUC := usecase.NewRegisterUser(userRepo, hasherSvc, clock, idGen)
	loginUC := usecase.NewLoginUser(userRepo, hasherSvc, tokenSvc)
	createListingUC := usecase.NewCreateListing(listRepo, clock)
	listListingsUC := usecase.NewListListings(listRepo)

	// Handlers
	authH := handler.NewAuthHandler(regUC, loginUC)
	listH := handler.NewListingHandler(createListingUC, listListingsUC)

	router := httpdelivery.NewRouter(authH, listH, tokenSvc.(tokenParserAdapter))

	srv := &http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: router,
	}

	return &App{
		Server: srv,
		DB:     db,
	}, nil
}

// Адаптер для приведения tokenSvc к нужному интерфейсу router
type tokenParserAdapter interface {
	Parse(token string) (port.Claims, error)
}

func (a *App) Shutdown(ctx context.Context) {
	if err := a.Server.Shutdown(ctx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}
	a.DB.Close()
}
