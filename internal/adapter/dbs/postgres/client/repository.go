package client

import (
	"RateBalancer/internal/adapter/dbs"
	postgres "RateBalancer/internal/adapter/dbs/postgres"
	"RateBalancer/internal/adapter/dbs/postgres/entity"
	"RateBalancer/internal/model"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(ctx context.Context, c *model.Client) error {
	query := ` INSERT INTO clients (id, api_key, tokens, last_refill, capacity, per_second) 
			   VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.db.ExecContext(ctx, query,
		c.Id, c.ApiKey, c.Tokens, c.LastRefill, c.Capacity, c.PerSecond)

	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf(`create client %w`, dbs.ErrorRecordAlreadyExists)
		}
		return fmt.Errorf(`create client %w`, err)
	}

	return nil
}

func (r *Repository) Get(ctx context.Context, id string) (*model.Client, error) {
	query := ` SELECT *
			   FROM clients
			   WHERE id = $1`

	client := &entity.Client{}
	err := r.db.GetContext(ctx, client, query, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf(`get client by id %w`, dbs.ErrorRecordNotFound)
		}
		return nil, fmt.Errorf(`get client by id %w`, err)
	}

	return postgres.ToClientServiceModel(client), nil
}

func (r *Repository) Update(ctx context.Context, c *model.Client) (*model.Client, error) {
	query := ` UPDATE clients
			   SET capacity = $1, per_second = $2
         	   WHERE id = $3
    		   RETURNING * `

	row := r.db.QueryRowxContext(ctx, query, c.Capacity, c.PerSecond, c.Id)

	client := &entity.Client{}
	if err := row.StructScan(client); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("update and fetch client %w", dbs.ErrorRecordNotFound)
		}
		return nil, fmt.Errorf("update  client %w", err)
	}

	return postgres.ToClientServiceModel(client), nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	const query = ` DELETE FROM clients
         			WHERE id = $1`

	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete client %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete client %w", err)
	}
	if n == 0 {
		return fmt.Errorf("delete client %w", dbs.ErrorRecordNotFound)
	}
	return nil
}
