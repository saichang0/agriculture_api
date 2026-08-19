package main

import (
	"context"
	"errors"
	"log"
	"net/http"

	"agriculture-api/graph"
	"agriculture-api/internal/auth"
	"agriculture-api/internal/config"
	"agriculture-api/internal/db"
	"agriculture-api/internal/repository"
	"agriculture-api/internal/scansession"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/errcode"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/rs/cors"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func main() {
	cfg := config.Load()

	client, err := db.Connect(cfg.MongoURI)
	if err != nil {
		log.Fatalf("failed to connect to mongo: %v", err)
	}
	database := client.Database(cfg.MongoDBName)

	jwtManager := auth.NewJWTManager(cfg.JWTSecret, cfg.JWTTTL)

	resolver := &graph.Resolver{
		UserRepo:           repository.NewUserRepository(database),
		CategoryRepo:       repository.NewCategoryRepository(database),
		UnitRepo:           repository.NewUnitRepository(database),
		ProductRepo:        repository.NewProductRepository(database),
		CustomerRepo:       repository.NewCustomerRepository(database),
		ImportRepo:         repository.NewImportRepository(database),
		SaleRepo:           repository.NewSaleRepository(database),
		DebtPaymentRepo:    repository.NewDebtPaymentRepository(database),
		DamagedProductRepo: repository.NewDamagedProductRepository(database),
		ExpenseRepo:        repository.NewExpenseRepository(database),
		RefreshTokenRepo:   repository.NewRefreshTokenRepository(database),
		JWT:                jwtManager,
		RefreshTokenTTL:    cfg.RefreshTokenTTL,
		ScanBroker:         scansession.NewBroker(),
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	// Marks auth.ErrUnauthenticated with extensions.code = "UNAUTHENTICATED" so the
	// frontend can distinguish "access token expired, please refresh" from any
	// other resolver error without relying on a fragile message-string match.
	srv.SetErrorPresenter(func(ctx context.Context, err error) *gqlerror.Error {
		gqlErr := graphql.DefaultErrorPresenter(ctx, err)
		if errors.Is(err, auth.ErrUnauthenticated) {
			errcode.Set(gqlErr, "UNAUTHENTICATED")
		}
		return gqlErr
	})

	mux := http.NewServeMux()
	mux.Handle("/", playground.Handler("GraphQL playground", "/query"))
	mux.Handle("/query", srv)

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   cfg.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, corsHandler.Handler(auth.Middleware(jwtManager)(mux))))
}
