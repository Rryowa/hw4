package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"homework/internal/models"
	"homework/internal/storage"
	"homework/internal/util"
	"log"
)

type Repository struct {
	pool *pgxpool.Pool
	ctx  context.Context
}

func NewSQLRepository(ctx context.Context, cfg *models.Config) storage.Storage {
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	var pool *pgxpool.Pool
	var err error

	err = util.DoWithTries(func() error {
		ctxTimeout, cancel := context.WithTimeout(ctx, cfg.Timeout)
		defer cancel()

		pool, err = pgxpool.New(ctxTimeout, connStr)
		if err != nil {
			log.Fatal(err, "db connection error")
		}

		return nil
	}, cfg.Attempts, cfg.Timeout)

	if err != nil {
		log.Fatal(err, "DoWithTries error")
	}
	log.Println("Connected to db")

	return &Repository{
		pool: pool,
		ctx:  ctx,
	}
}

func (r *Repository) Insert(order models.Order) error {
	query := `
		INSERT INTO orders (id, user_id, storage_until, issued, issued_at, returned, order_price, weight, package_type, package_price, hash) 
	    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	    `

	_, err := r.pool.Exec(r.ctx, query, order.ID, order.UserID, order.StorageUntil, order.Issued, order.IssuedAt, order.Returned, order.OrderPrice, order.Weight, order.PackageType, order.PackagePrice, order.Hash)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			log.Println(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s", pgErr.Code, pgErr.Detail, pgErr.Where))
		}
		return err
	}
	return nil
}

func (r *Repository) Update(order models.Order) error {
	query := `
		UPDATE orders SET returned=$1
        WHERE id=$2
        `

	_, err := r.pool.Exec(r.ctx, query, order.Returned, order.ID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			log.Println(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s", pgErr.Code, pgErr.Detail, pgErr.Where))
		}
		return err
	}
	return nil
}

func (r *Repository) IssueUpdate(orders []models.Order) error {
	tx, err := r.pool.BeginTx(r.ctx, pgx.TxOptions{
		IsoLevel:   pgx.RepeatableRead,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return err
	}
	defer tx.Rollback(r.ctx)

	query := `
		UPDATE orders SET issued=$1, issued_at=$2
        WHERE id=$3
        `
	batch := &pgx.Batch{}
	for _, order := range orders {
		batch.Queue(query, order.Issued, order.IssuedAt, order.ID)
		log.Printf("Order with id:%s issued\n", order.ID)
	}

	br := tx.SendBatch(r.ctx, batch)
	for _, i := range orders {
		_, err := br.Exec()
		if err != nil {
			br.Close()
			return fmt.Errorf("error executing batch at order index %d: %w", i, err)
		}
	}
	err = br.Close()

	return tx.Commit(r.ctx)
}

func (r *Repository) Delete(id string) error {
	query := `
		DELETE FROM orders WHERE id=$1
		`

	_, err := r.pool.Exec(r.ctx, query, id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			log.Println(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s", pgErr.Code, pgErr.Detail, pgErr.Where))
		}
		return err
	}

	return nil
}

func (r *Repository) Get(id string) (models.Order, error) {
	var order models.Order
	query := `
		SELECT id, user_id, storage_until, issued, issued_at, returned, order_price, weight, package_type, package_price, hash FROM orders
		WHERE id=$1
		`
	if err := pgxscan.Get(r.ctx, r.pool, &order, query, id); err != nil {
		return models.Order{}, err
	}
	return order, nil
}

func (r *Repository) GetReturns(offset, limit int) ([]models.Order, error) {
	query := `
        SELECT id, user_id, storage_until, issued, issued_at, returned, order_price, weight, package_type, package_price, hash
        FROM orders
        WHERE returned = TRUE
        ORDER BY id
        OFFSET $1
 		FETCH NEXT $2 ROWS ONLY
    `

	rows, err := r.pool.Query(r.ctx, query, offset, limit)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			log.Println(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s", pgErr.Code, pgErr.Detail, pgErr.Where))
		}
		return nil, err
	}

	defer rows.Close()

	var returns []models.Order
	if err := pgxscan.ScanAll(&returns, rows); err != nil {
		return nil, err
	}
	return returns, nil
}

func (r *Repository) GetOrders(userId string, offset, limit int) ([]models.Order, error) {
	query := `
		SELECT id, user_id, storage_until, issued, issued_at, returned, order_price, weight, package_type, package_price, hash
		FROM orders
		WHERE user_id = $1 AND issued = FALSE
		ORDER BY storage_until
		OFFSET $2
		FETCH NEXT $3 ROWS ONLY
	`

	rows, err := r.pool.Query(r.ctx, query, userId, offset, limit)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			log.Println(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s", pgErr.Code, pgErr.Detail, pgErr.Where))
		}
		return nil, err
	}
	defer rows.Close()

	var userOrders []models.Order
	if err := pgxscan.ScanAll(&userOrders, rows); err != nil {
		return nil, err
	}
	return userOrders, err
}
