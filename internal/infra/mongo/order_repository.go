package mongo

import (
	"context"
	"errors"
	"time"

	"github.com/renan/go-api-worker/internal/application/port"
	"github.com/renan/go-api-worker/internal/domain/order"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type orderRepository struct {
	coll *mongo.Collection
}

func NewOrderRepository(db *mongo.Database) (*orderRepository, error) {
	coll := db.Collection("orders")
	r := &orderRepository{coll: coll}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := coll.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "order_id", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (r *orderRepository) Create(ctx context.Context, o order.Order) error {
	_, err := r.coll.InsertOne(ctx, o)
	return err
}

func (r *orderRepository) GetByID(ctx context.Context, orderID string) (order.Order, error) {
	var o order.Order
	err := r.coll.FindOne(ctx, bson.D{{Key: "order_id", Value: orderID}}).Decode(&o)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return order.Order{}, port.ErrNotFound{}
		}
		return order.Order{}, err
	}
	return o, nil
}

func (r *orderRepository) UpdateStatus(ctx context.Context, orderID string, status order.Status) error {
	res, err := r.coll.UpdateOne(
		ctx,
		bson.D{{Key: "order_id", Value: orderID}},
		bson.D{{Key: "$set", Value: bson.D{{Key: "status", Value: status}, {Key: "updated_at", Value: time.Now()}}}},
	)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return port.ErrNotFound{}
	}
	return nil
}

func IsNotFound(err error) bool {
	return errors.Is(err, port.ErrNotFound{})
}
