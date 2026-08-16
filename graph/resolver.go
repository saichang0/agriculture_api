package graph

import (
	"agriculture-api/internal/auth"
	"agriculture-api/internal/repository"
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
	JWT                *auth.JWTManager
}
