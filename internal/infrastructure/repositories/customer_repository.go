package repositories

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"

	"github.com/terminator791/Domain-Driven-Design-Architecture-GO/internal/domain/common"
	"github.com/terminator791/Domain-Driven-Design-Architecture-GO/internal/domain/customer"
	"github.com/terminator791/Domain-Driven-Design-Architecture-GO/internal/domain/order"
)

type PostgresCustomerRepository struct {
	db *sqlx.DB
}

func NewPostgresCustomerRepository(db *sqlx.DB) *PostgresCustomerRepository {
	return &PostgresCustomerRepository{db: db}
}

type customerRecord struct {
	ID    string `db:"id"`
	Name  string `db:"name"`
	Email string `db:"email"`
	Phone string `db:"phone"`
}

func (r *PostgresCustomerRepository) Save(ctx context.Context, c *customer.Customer) error {
	query := `
		INSERT INTO customers (id, name, email, phone)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			email = EXCLUDED.email,
			phone = EXCLUDED.phone,
			updated_at = CURRENT_TIMESTAMP`

	_, err := r.db.ExecContext(ctx, query,
		c.ID().String(),
		c.Name(),
		c.Email(),
		c.Phone(),
	)

	return err
}

func (r *PostgresCustomerRepository) FindByID(ctx context.Context, id order.CustomerID) (*customer.Customer, error) {
	var record customerRecord
	query := `SELECT id, name, email, phone FROM customers WHERE id = $1`

	err := r.db.GetContext(ctx, &record, query, id.String())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, common.ErrCustomerNotFound
		}
		return nil, err
	}

	return r.mapToCustomer(record)
}

func (r *PostgresCustomerRepository) FindByEmail(ctx context.Context, email string) (*customer.Customer, error) {
	var record customerRecord
	query := `SELECT id, name, email, phone FROM customers WHERE email = $1`

	err := r.db.GetContext(ctx, &record, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, common.ErrCustomerNotFound
		}
		return nil, err
	}

	return r.mapToCustomer(record)
}

func (r *PostgresCustomerRepository) FindAll(ctx context.Context) ([]*customer.Customer, error) {
	var records []customerRecord
	query := `SELECT id, name, email, phone FROM customers ORDER BY name`

	err := r.db.SelectContext(ctx, &records, query)
	if err != nil {
		return nil, err
	}

	var customers []*customer.Customer
	for _, record := range records {
		c, err := r.mapToCustomer(record)
		if err != nil {
			return nil, err
		}
		customers = append(customers, c)
	}

	return customers, nil
}

func (r *PostgresCustomerRepository) Delete(ctx context.Context, id order.CustomerID) error {
	query := `DELETE FROM customers WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id.String())
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return common.ErrCustomerNotFound
	}

	return nil
}

func (r *PostgresCustomerRepository) mapToCustomer(record customerRecord) (*customer.Customer, error) {
	c := customer.NewCustomer(record.Name, record.Email, record.Phone)
	return c, nil
}