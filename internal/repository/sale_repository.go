package repository

import (
	"context"
	"fmt"
	"time"

	"agriculture-api/graph/model"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type SaleRepository struct {
	salesCol     *mongo.Collection
	saleItemsCol *mongo.Collection
	countersCol  *mongo.Collection
}

func NewSaleRepository(db *mongo.Database) *SaleRepository {
	return &SaleRepository{
		salesCol:     db.Collection("sales"),
		saleItemsCol: db.Collection("saleItems"),
		countersCol:  db.Collection("counters"),
	}
}

// nextSaleCode atomically reserves the next sequence number for today (Vientiane time)
// and returns a code like "260813-0001". The counter resets automatically each day
// because it is keyed by date string, not shared across days.
func (r *SaleRepository) nextSaleCode(ctx context.Context) (string, error) {
	loc, err := time.LoadLocation("Asia/Vientiane")
	if err != nil {
		loc = time.UTC
	}
	datePart := time.Now().In(loc).Format("060102")
	counterID := "sale-" + datePart

	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
	result := r.countersCol.FindOneAndUpdate(
		ctx,
		bson.M{"_id": counterID},
		bson.M{"$inc": bson.M{"seq": 1}},
		opts,
	)

	var counter struct {
		Seq int `bson:"seq"`
	}
	if err := result.Decode(&counter); err != nil {
		return "", err
	}

	return fmt.Sprintf("%s-%04d", datePart, counter.Seq), nil
}

func (r *SaleRepository) FindAll(ctx context.Context) ([]*model.SaleDoc, error) {
	cur, err := r.salesCol.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var sales []*model.SaleDoc
	if err := cur.All(ctx, &sales); err != nil {
		return nil, err
	}
	return sales, nil
}

func (r *SaleRepository) FindByID(ctx context.Context, id bson.ObjectID) (*model.SaleDoc, error) {
	var sale model.SaleDoc
	err := r.salesCol.FindOne(ctx, bson.M{"_id": id}).Decode(&sale)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &sale, nil
}

func (r *SaleRepository) FindByCustomer(ctx context.Context, customerID bson.ObjectID) ([]*model.SaleDoc, error) {
	cur, err := r.salesCol.Find(ctx, bson.M{"customerId": customerID})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var sales []*model.SaleDoc
	if err := cur.All(ctx, &sales); err != nil {
		return nil, err
	}
	return sales, nil
}

func (r *SaleRepository) ItemsBySale(ctx context.Context, saleID bson.ObjectID) ([]*model.SaleItemDoc, error) {
	cur, err := r.saleItemsCol.Find(ctx, bson.M{"saleId": saleID})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var items []*model.SaleItemDoc
	if err := cur.All(ctx, &items); err != nil {
		return nil, err
	}
	return items, nil
}

// Create inserts the sale and its items. Not run in a Mongo transaction (requires a
// replica set, which a single local `docker compose` instance doesn't provide by
// default) — stock/debt adjustments happen in the resolver alongside this call.
func (r *SaleRepository) Create(ctx context.Context, sale *model.SaleDoc, items []*model.SaleItemDoc) error {
	code, err := r.nextSaleCode(ctx)
	if err != nil {
		return fmt.Errorf("generate sale code: %w", err)
	}

	sale.ID = bson.NewObjectID()
	sale.Code = code

	if _, err := r.salesCol.InsertOne(ctx, sale); err != nil {
		return err
	}

	docs := make([]interface{}, 0, len(items))
	for _, item := range items {
		item.ID = bson.NewObjectID()
		item.SaleID = sale.ID
		docs = append(docs, item)
	}
	if len(docs) > 0 {
		if _, err := r.saleItemsCol.InsertMany(ctx, docs); err != nil {
			return err
		}
	}

	return nil
}

// ApplyPayment adds amount to the sale's paid total, reduces debt accordingly (never
// below zero), and recalculates paymentStatus. Returns the updated sale, and the
// amount actually applied (capped at the remaining debt, in case of overpayment).
func (r *SaleRepository) ApplyPayment(ctx context.Context, id bson.ObjectID, amount float64) (*model.SaleDoc, float64, error) {
	sale, err := r.FindByID(ctx, id)
	if err != nil {
		return nil, 0, err
	}
	if sale == nil {
		return nil, 0, nil
	}

	applied := amount
	if applied > sale.Debt {
		applied = sale.Debt
	}

	newPaid := sale.Paid + applied
	newDebt := sale.Debt - applied

	status := "PARTIAL"
	if newDebt <= 0 {
		status = "PAID"
	} else if newPaid <= 0 {
		status = "UNPAID"
	}

	_, err = r.salesCol.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{
		"paid":          newPaid,
		"debt":          newDebt,
		"paymentStatus": status,
	}})
	if err != nil {
		return nil, 0, err
	}

	sale.Paid = newPaid
	sale.Debt = newDebt
	sale.PaymentStatus = status
	return sale, applied, nil
}

// ReplaceItems swaps out a sale's line items and updates its totals/payment fields in
// place — the sale keeps its original _id, code, and saleDate. Stock and customer debt
// deltas are the caller's responsibility (see UpdateSale resolver): this only writes
// the sale + saleItems documents.
func (r *SaleRepository) ReplaceItems(ctx context.Context, id bson.ObjectID, items []*model.SaleItemDoc, subtotal, discount, total, paid, debt float64, paymentStatus string) error {
	if _, err := r.saleItemsCol.DeleteMany(ctx, bson.M{"saleId": id}); err != nil {
		return err
	}

	docs := make([]interface{}, 0, len(items))
	for _, item := range items {
		item.ID = bson.NewObjectID()
		item.SaleID = id
		docs = append(docs, item)
	}
	if len(docs) > 0 {
		if _, err := r.saleItemsCol.InsertMany(ctx, docs); err != nil {
			return err
		}
	}

	_, err := r.salesCol.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{
		"subtotal":      subtotal,
		"discount":      discount,
		"total":         total,
		"paid":          paid,
		"debt":          debt,
		"paymentStatus": paymentStatus,
	}})
	return err
}

func (r *SaleRepository) Delete(ctx context.Context, id bson.ObjectID) (bool, error) {
	res, err := r.salesCol.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return false, err
	}
	if res.DeletedCount == 0 {
		return false, nil
	}

	_, err = r.saleItemsCol.DeleteMany(ctx, bson.M{"saleId": id})
	return true, err
}
