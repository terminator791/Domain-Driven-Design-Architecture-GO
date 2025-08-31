package repositories

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"

	"github.com/terminator791/Domain-Driven-Design-Architecture-GO/internal/domain/common"
	"github.com/terminator791/Domain-Driven-Design-Architecture-GO/internal/domain/order"
	"github.com/terminator791/Domain-Driven-Design-Architecture-GO/internal/domain/product"
)

type PostgresProductRepository struct {
	db *sqlx.DB
}

func NewPostgresProductRepository(db *sqlx.DB) *PostgresProductRepository {
	return &PostgresProductRepository{db: db}
}

type productRecord struct {
	ID            string  `db:"id"`
	Name          string  `db:"name"`
	Description   string  `db:"description"`
	PriceAmount   float64 `db:"price_amount"`
	PriceCurrency string  `db:"price_currency"`
	StockLevel    int     `db:"stock_level"`
	IsAvailable   bool    `db:"is_available"`
}

func (r *PostgresProductRepository) Save(ctx context.Context, p *product.Product) error {
	query := `
		INSERT INTO products (id, name, description, price_amount, price_currency, stock_level, is_available)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			description = EXCLUDED.description,
			price_amount = EXCLUDED.price_amount,
			price_currency = EXCLUDED.price_currency,
			stock_level = EXCLUDED.stock_level,
			is_available = EXCLUDED.is_available,
			updated_at = CURRENT_TIMESTAMP`

	_, err := r.db.ExecContext(ctx, query,
		p.ID().String(),
		p.Name(),
		p.Description(),
		p.Price().Amount(),
		p.Price().Currency(),
		p.StockLevel(),
		p.IsAvailable(),
	)

	return err
}

func (r *PostgresProductRepository) FindByID(ctx context.Context, id order.ProductID) (*product.Product, error) {
	var record productRecord
	query := `SELECT id, name, description, price_amount, price_currency, stock_level, is_available FROM products WHERE id = $1`

	err := r.db.GetContext(ctx, &record, query, id.String())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, common.ErrProductNotFound
		}
		return nil, err
	}

	return r.mapToProduct(record)
}

func (r *PostgresProductRepository) FindAll(ctx context.Context) ([]*product.Product, error) {
	var records []productRecord
	query := `SELECT id, name, description, price_amount, price_currency, stock_level, is_available FROM products ORDER BY name`

	err := r.db.SelectContext(ctx, &records, query)
	if err != nil {
		return nil, err
	}

	var products []*product.Product
	for _, record := range records {
		p, err := r.mapToProduct(record)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

func (r *PostgresProductRepository) Delete(ctx context.Context, id order.ProductID) error {
	query := `DELETE FROM products WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id.String())
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return common.ErrProductNotFound
	}

	return nil
}

func (r *PostgresProductRepository) mapToProduct(record productRecord) (*product.Product, error) {
	price, err := order.NewMoney(record.PriceAmount, record.PriceCurrency)
	if err != nil {
		return nil, err
	}

	p := product.NewProduct(record.Name, record.Description, price, record.StockLevel)
	
	// Set availability
	if !record.IsAvailable {
		p.UpdateStock(0) // This will set availability to false
	}

	return p, nil
}