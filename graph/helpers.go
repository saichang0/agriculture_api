package graph

import (
	"fmt"
	"strconv"

	"agriculture-api/graph/model"

	"go.mongodb.org/mongo-driver/v2/bson"
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
	var unitID *string
	if d.UnitID != nil {
		hex := d.UnitID.Hex()
		unitID = &hex
	}

	factor := d.Factor
	if factor == 0 {
		// Sales created before the Factor field existed default to 1 (base unit).
		factor = 1
	}

	return &model.SaleItem{
		ID:        d.ID.Hex(),
		SaleID:    d.SaleID.Hex(),
		ProductID: d.ProductID.Hex(),
		Quantity:  d.Quantity,
		CostPrice: d.CostPrice,
		UnitPrice: d.UnitPrice,
		PriceType: d.PriceType,
		Subtotal:  d.Subtotal,
		UnitID:    unitID,
		Factor:    factor,
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

func toGraphPackagingUnit(d model.PackagingUnitDoc) *model.PackagingUnit {
	return &model.PackagingUnit{
		UnitID:          d.UnitID.Hex(),
		Factor:          d.Factor,
		CostPrice:       d.CostPrice,
		RetailPrice:     d.RetailPrice,
		WholesalePrice:  d.WholesalePrice,
		WholesaleMinQty: int(d.WholesaleMinQty),
	}
}

func toGraphProduct(d *model.ProductDoc) *model.Product {
	packagingUnits := make([]*model.PackagingUnit, 0, len(d.PackagingUnits))
	for _, pu := range d.PackagingUnits {
		packagingUnits = append(packagingUnits, toGraphPackagingUnit(pu))
	}

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
		PackagingUnits:  packagingUnits,
	}
}

// packagingUnitDocsFromInput converts GraphQL packaging-unit inputs to Mongo docs,
// validating unitId hex and factor > 0.
func packagingUnitDocsFromInput(inputs []*model.PackagingUnitInput) ([]model.PackagingUnitDoc, error) {
	docs := make([]model.PackagingUnitDoc, 0, len(inputs))
	for _, in := range inputs {
		unitID, err := bson.ObjectIDFromHex(in.UnitID)
		if err != nil {
			return nil, fmt.Errorf("invalid packaging unit unitId: %w", err)
		}
		if in.Factor <= 0 {
			return nil, fmt.Errorf("packaging unit factor must be greater than zero")
		}
		docs = append(docs, model.PackagingUnitDoc{
			UnitID:          unitID,
			Factor:          in.Factor,
			CostPrice:       in.CostPrice,
			RetailPrice:     in.RetailPrice,
			WholesalePrice:  in.WholesalePrice,
			WholesaleMinQty: int32(in.WholesaleMinQty),
		})
	}
	return docs, nil
}

// resolveSaleUnit determines the pricing/factor for a sale line, given an optional
// packaging-unit id. unitIDHex == nil means the product's base unit: factor 1, prices
// straight from the product (this keeps single-unit products behaving exactly as
// before this feature existed). Otherwise it looks up the matching packaging unit —
// its own prices/factor are used verbatim, never derived from the base unit's price.
func resolveSaleUnit(product *model.ProductDoc, unitIDHex *string) (factor float64, costPrice, retailPrice, wholesalePrice float64, wholesaleMinQty int32, unitID *bson.ObjectID, err error) {
	if unitIDHex == nil {
		return 1, product.CostPrice, product.RetailPrice, product.WholesalePrice, product.WholesaleMinQty, nil, nil
	}

	oid, err := bson.ObjectIDFromHex(*unitIDHex)
	if err != nil {
		return 0, 0, 0, 0, 0, nil, fmt.Errorf("invalid unitId: %w", err)
	}

	for _, pu := range product.PackagingUnits {
		if pu.UnitID == oid {
			return pu.Factor, pu.CostPrice, pu.RetailPrice, pu.WholesalePrice, pu.WholesaleMinQty, &oid, nil
		}
	}

	return 0, 0, 0, 0, 0, nil, fmt.Errorf("product %s has no packaging unit %s", product.Name, *unitIDHex)
}
