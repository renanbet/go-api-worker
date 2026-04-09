package order

import (
	"errors"
	"fmt"
	"time"
)

type Status string

const (
	StatusCreated    Status = "CREATED"
	StatusProcessing Status = "PROCESSING"
	StatusConcluded  Status = "CONCLUDED"
)

func (s Status) Validate() error {
	switch s {
	case StatusCreated, StatusProcessing, StatusConcluded:
		return nil
	}
	return fmt.Errorf("invalid status: %s", s)
}

func (o Order) Validate() error {
	if o.ProductName == "" {
		return errors.New("product name is required")
	}
	if o.Quantity <= 0 {
		return errors.New("quantity must be greater than 0")
	}
	if err := o.Status.Validate(); err != nil {
		return err
	}
	return nil
}

type Order struct {
	OrderID     string    `bson:"order_id" json:"order_id"`
	ProductName string    `bson:"product_name" json:"product_name"`
	Quantity    int       `bson:"quantity" json:"quantity"`
	Status      Status    `bson:"status" json:"status"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
}
