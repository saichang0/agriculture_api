package graph

import (
	"strconv"

	"agriculture-api/graph/model"
)

func toGraphUser(d *model.UserDoc) *model.User {
	return &model.User{
		ID:        d.ID.Hex(),
		Username:  d.Username,
		Role:      d.Role,
		FirstName: d.FirstName,
		LastName:  d.LastName,
		Phone:     d.Phone,
		Status:    d.Status,
	}
}

func toGraphCategory(d *model.CategoryDoc) *model.Category {
	return &model.Category{
		ID:   d.ID.Hex(),
		Name: d.Name,
	}
}

func toGraphUnit(d *model.UnitDoc) *model.Unit {
	return &model.Unit{
		ID:   d.ID.Hex(),
		Name: d.Name,
	}
}

func toGraphCustomer(d *model.CustomerDoc) *model.Customer {
	return &model.Customer{
		ID:      d.ID.Hex(),
		Name:    d.Name,
		Phone:   d.Phone,
		Address: d.Address,
		Debt:    d.Debt,
		Status:  d.Status,
	}
}

func toGraphImport(d *model.ImportDoc) *model.Import {
	return &model.Import{
		ID:        d.ID.Hex(),
		ProductID: d.ProductID.Hex(),
		Quantity:  d.Quantity,
		CostPrice: d.CostPrice,
		UserID:    d.UserID.Hex(),
		Note:      d.Note,
		Date:      strconv.FormatInt(d.Date, 10),
	}
}

func toGraphSaleItem(d *model.SaleItemDoc) *model.SaleItem {
	return &model.SaleItem{
		ID:        d.ID.Hex(),
		SaleID:    d.SaleID.Hex(),
		ProductID: d.ProductID.Hex(),
		Quantity:  d.Quantity,
		CostPrice: d.CostPrice,
		UnitPrice: d.UnitPrice,
		PriceType: d.PriceType,
		Subtotal:  d.Subtotal,
	}
}

func toGraphSale(d *model.SaleDoc, items []*model.SaleItemDoc) *model.Sale {
	var customerID *string
	if d.CustomerID != nil {
		hex := d.CustomerID.Hex()
		customerID = &hex
	}

	graphItems := make([]*model.SaleItem, 0, len(items))
	for _, item := range items {
		graphItems = append(graphItems, toGraphSaleItem(item))
	}

	return &model.Sale{
		ID:            d.ID.Hex(),
		Code:          d.Code,
		CustomerID:    customerID,
		UserID:        d.UserID.Hex(),
		SaleDate:      strconv.FormatInt(d.SaleDate, 10),
		Total:         d.Total,
		Paid:          d.Paid,
		Debt:          d.Debt,
		PaymentStatus: d.PaymentStatus,
		DueDate:       d.DueDate,
		PaymentMethod: d.PaymentMethod,
		Items:         graphItems,
	}
}

func toGraphDebtPayment(d *model.DebtPaymentDoc) *model.DebtPayment {
	return &model.DebtPayment{
		ID:          d.ID.Hex(),
		SaleID:      d.SaleID.Hex(),
		CustomerID:  d.CustomerID.Hex(),
		UserID:      d.UserID.Hex(),
		AmountPaid:  d.AmountPaid,
		PaymentDate: strconv.FormatInt(d.PaymentDate, 10),
		Note:        d.Note,
	}
}

func toGraphDamagedProduct(d *model.DamagedProductDoc) *model.DamagedProduct {
	return &model.DamagedProduct{
		ID:        d.ID.Hex(),
		ProductID: d.ProductID.Hex(),
		UserID:    d.UserID.Hex(),
		Quantity:  d.Quantity,
		CostPrice: d.CostPrice,
		Reason:    d.Reason,
		Note:      d.Note,
		Date:      strconv.FormatInt(d.Date, 10),
	}
}

func toGraphExpense(d *model.ExpenseDoc) *model.Expense {
	return &model.Expense{
		ID:     d.ID.Hex(),
		UserID: d.UserID.Hex(),
		Title:  d.Title,
		Type:   d.Type,
		Amount: d.Amount,
		Date:   strconv.FormatInt(d.Date, 10),
	}
}

func toGraphProduct(d *model.ProductDoc) *model.Product {
	return &model.Product{
		ID:              d.ID.Hex(),
		Barcode:         d.Barcode,
		Name:            d.Name,
		ImageURL:        d.ImageURL,
		CategoryID:      d.CategoryID.Hex(),
		UnitID:          d.UnitID.Hex(),
		CostPrice:       d.CostPrice,
		RetailPrice:     d.RetailPrice,
		WholesalePrice:  d.WholesalePrice,
		WholesaleMinQty: int(d.WholesaleMinQty),
		StockQty:        d.StockQty,
		MinStockAlert:   d.MinStockAlert,
		Status:          d.Status,
	}
}
