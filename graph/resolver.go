package graph

import (
	"time"

	"agriculture-api/internal/auth"
	"agriculture-api/internal/repository"
	"agriculture-api/internal/scansession"
)

type Resolver struct {
	UserRepo           *repository.UserRepository
	CategoryRepo       *repository.CategoryRepository
	UnitRepo           *repository.UnitRepository
	ProductRepo        *repository.ProductRepository
	CustomerRepo       *repository.CustomerRepository
	ImportRepo         *repository.ImportRepository
	SaleRepo           *repository.SaleRepository
	DebtPaymentRepo    *repository.DebtPaymentRepository
	DamagedProductRepo *repository.DamagedProductRepository
	ExpenseRepo        *repository.ExpenseRepository
	RefreshTokenRepo   *repository.RefreshTokenRepository
	JWT                *auth.JWTManager
	RefreshTokenTTL    time.Duration
	ScanBroker         *scansession.Broker
}
