package server

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

func NewHandler(cfg config.Config) (http.Handler, error) {
	client, err := db.Connect(cfg.MongoURI)
	if err != nil {
		return nil, err
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
	srv.SetErrorPresenter(func(ctx context.Context, err error) *gqlerror.Error {
		gqlErr := graphql.DefaultErrorPresenter(ctx, err)
		if errors.Is(err, auth.ErrUnauthenticated) {
			errcode.Set(gqlErr, "UNAUTHENTICATED")
		}
		return gqlErr
	})

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	mux.Handle("/", playground.Handler("GraphQL playground", "/query"))
	mux.Handle("/query", auth.Middleware(jwtManager)(srv))

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   cfg.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	return corsHandler.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("incoming method=%s path=%s origin=%q", r.Method, r.URL.Path, r.Header.Get("Origin"))
		mux.ServeHTTP(w, r)
	})), nil
}
