package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/terminator791/Domain-Driven-Design-Architecture-GO/internal/domain/common"
	"github.com/terminator791/Domain-Driven-Design-Architecture-GO/internal/domain/order"
)

type PostgresOrderRepository struct {
	db *sqlx.DB
}

func NewPostgresOrderRepository(db *sqlx.DB) *PostgresOrderRepository {
	return &PostgresOrderRepository{db: db}
}

type orderRecord struct {
	ID            string    `db:"id"`
	CustomerID    string    `db:"customer_id"`
	Status        string    `db:"status"`
	TotalAmount   float64   `db:"total_amount"`
	TotalCurrency string    `db:"total_currency"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

type orderItemRecord struct {
	ID                string  `db:"id"`
	OrderID           string  `db:"order_id"`
	ProductID         string  `db:"product_id"`
	Quantity          int     `db:"quantity"`
	UnitPriceAmount   float64 `db:"unit_price_amount"`
	UnitPriceCurrency string  `db:"unit_price_currency"`
	TotalAmount       float64 `db:"total_amount"`
	TotalCurrency     string  `db:"total_currency"`
}

func (r *PostgresOrderRepository) Save(ctx context.Context, o *order.Order) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Save or update order
	query := `
		INSERT INTO orders (id, customer_id, status, total_amount, total_currency, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (id) DO UPDATE SET
			status = EXCLUDED.status,
			total_amount = EXCLUDED.total_amount,
			total_currency = EXCLUDED.total_currency,
			updated_at = EXCLUDED.updated_at`

	_, err = tx.ExecContext(ctx, query,
		o.ID().String(),
		o.CustomerID().String(),
		o.Status().String(),
		o.Total().Amount(),
		o.Total().Currency(),
		o.CreatedAt(),
		o.UpdatedAt(),
	)
	if err != nil {
		return err
	}

	// Delete existing order items
	_, err = tx.ExecContext(ctx, "DELETE FROM order_items WHERE order_id = $1", o.ID().String())
	if err != nil {
		return err
	}

	// Insert order items
	for _, item := range o.Items() {
		itemQuery := `
			INSERT INTO order_items (order_id, product_id, quantity, unit_price_amount, unit_price_currency, total_amount, total_currency)
			VALUES ($1, $2, $3, $4, $5, $6, $7)`

		_, err = tx.ExecContext(ctx, itemQuery,
			o.ID().String(),
			item.ProductID().String(),
			item.Quantity().Value(),
			item.UnitPrice().Amount(),
			item.UnitPrice().Currency(),
			item.TotalPrice().Amount(),
			item.TotalPrice().Currency(),
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *PostgresOrderRepository) FindByID(ctx context.Context, id order.OrderID) (*order.Order, error) {
	var record orderRecord
	query := `SELECT id, customer_id, status, total_amount, total_currency, created_at, updated_at FROM orders WHERE id = $1`

	err := r.db.GetContext(ctx, &record, query, id.String())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, common.ErrOrderNotFound
		}
		return nil, err
	}

	// Get order items
	var itemRecords []orderItemRecord
	itemQuery := `SELECT id, order_id, product_id, quantity, unit_price_amount, unit_price_currency, total_amount, total_currency FROM order_items WHERE order_id = $1`

	err = r.db.SelectContext(ctx, &itemRecords, itemQuery, id.String())
	if err != nil {
		return nil, err
	}

	return r.mapToOrder(record, itemRecords)
}

func (r *PostgresOrderRepository) FindByCustomerID(ctx context.Context, customerID order.CustomerID) ([]*order.Order, error) {
	var records []orderRecord
	query := `SELECT id, customer_id, status, total_amount, total_currency, created_at, updated_at FROM orders WHERE customer_id = $1 ORDER BY created_at DESC`

	err := r.db.SelectContext(ctx, &records, query, customerID.String())
	if err != nil {
		return nil, err
	}

	var orders []*order.Order
	for _, record := range records {
		// Get order items for each order
		var itemRecords []orderItemRecord
		itemQuery := `SELECT id, order_id, product_id, quantity, unit_price_amount, unit_price_currency, total_amount, total_currency FROM order_items WHERE order_id = $1`

		err = r.db.SelectContext(ctx, &itemRecords, itemQuery, record.ID)
		if err != nil {
			return nil, err
		}

		o, err := r.mapToOrder(record, itemRecords)
		if err != nil {
			return nil, err
		}

		orders = append(orders, o)
	}

	return orders, nil
}

func (r *PostgresOrderRepository) Delete(ctx context.Context, id order.OrderID) error {
	query := `DELETE FROM orders WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id.String())
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return common.ErrOrderNotFound
	}

	return nil
}

func (r *PostgresOrderRepository) mapToOrder(record orderRecord, itemRecords []orderItemRecord) (*order.Order, error) {
	// Create customer ID
	customerID, err := order.NewCustomerIDFromString(record.CustomerID)
	if err != nil {
		return nil, err
	}

	// Create order
	o := order.NewOrder(customerID)

	// Set order fields using reflection or direct access (if fields were public)
	// For this example, we'll recreate the order with the stored data
	// In a real implementation, you might want to add reconstruction methods

	// Add items
	for _, itemRecord := range itemRecords {
		productID, err := order.NewProductIDFromString(itemRecord.ProductID)
		if err != nil {
			return nil, err
		}

		quantity, err := order.NewQuantity(itemRecord.Quantity)
		if err != nil {
			return nil, err
		}

		unitPrice, err := order.NewMoney(itemRecord.UnitPriceAmount, itemRecord.UnitPriceCurrency)
		if err != nil {
			return nil, err
		}

		err = o.AddItem(productID, quantity, unitPrice)
		if err != nil {
			return nil, err
		}
	}

	// Update status if not pending
	if record.Status != "PENDING" {
		status := order.OrderStatus(record.Status)
		if status == order.OrderStatusConfirmed {
			err = o.Confirm()
			if err != nil {
				return nil, err
			}
		} else {
			err = o.UpdateStatus(status)
			if err != nil {
				return nil, err
			}
		}
	}

	// Clear events as this is a reconstituted object
	o.ClearEvents()

	return o, nil
}

// EventPublisher is a simple implementation for domain events
type SimpleEventPublisher struct{}

func NewSimpleEventPublisher() *SimpleEventPublisher {
	return &SimpleEventPublisher{}
}

func (p *SimpleEventPublisher) Publish(ctx context.Context, event interface{}) error {
	// In a real implementation, this would publish to a message queue or event bus
	// For now, we'll just log the event
	fmt.Printf("Publishing event: %+v\n", event)
	return nil
}