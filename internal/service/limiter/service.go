package limiter

import (
	"RateBalancer/internal/adapter/dbs/postgres/entity"
	"RateBalancer/pkg/hash"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

type Service struct {
	db               *sqlx.DB
	hasher           hash.Hasher
	defaultCapacity  int64
	defaultPerSecond int64
}

func NewServiceLimiter(db *sqlx.DB, hasher hash.Hasher, c *Config) *Service {
	return &Service{
		db:               db,
		hasher:           hasher,
		defaultCapacity:  c.Capacity,
		defaultPerSecond: c.PerSecond,
	}
}

func (s *Service) ConsumeTokens(ctx context.Context, apiKey string) (bool, error) {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return false, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	var tokenRequested int64 = 1

	apiKeyHash, err := s.hasher.Hash(apiKey)
	if err != nil {
		return false, fmt.Errorf(`hash apikey %w`, err)
	}

	query := ` SELECT *
       		   FROM clients
               WHERE api_key = $1
               FOR UPDATE`

	client := &entity.Client{}

	err = tx.GetContext(ctx, client, query, apiKeyHash)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf(`get client by api_key %w`, InvalidAPIKey)
		}
		return false, fmt.Errorf(`select for update: %w`, err)
	}

	// подставляем дефолтные значения если в базе нет настроек для клиента
	capacity, perSecond := s.defaultCapacity, s.defaultPerSecond
	if client.Capacity != nil {
		capacity = int64(*client.Capacity)
		perSecond = int64(*client.PerSecond)
	}

	tokens, lastRefill := s.refill(capacity, perSecond, int64(client.Tokens), client.LastRefill)

	// Проверка на достаточность токенов
	if tokens < tokenRequested {
		if err := tx.Commit(); err != nil {
			return false, fmt.Errorf("commit (no tokens): %w", err)
		}
		return false, nil
	}

	// Списываем и обновляем токены
	tokens -= tokenRequested

	query = ` UPDATE clients
        	  SET tokens = $1, last_refill = $2
        	  WHERE api_key = $3`
	_, err = tx.ExecContext(ctx, query, tokens, lastRefill, apiKeyHash)
	if err != nil {
		return false, fmt.Errorf("update token_bucket: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return false, fmt.Errorf("commit: %w", err)
	}
	return true, nil
}

// refill
// nanosToGenerateOneToken - время в наносекундах для генерации одного токена.
// duration - количество наносекунд которое прошло с последнего обновления.
// tokensSinceLastRefill - количество токенов которые мы можем добавить.
// при целочисленном делении теряется остаток, чтобы учитывать это
// lastRefill двигается не на текущее время (now), а ровно на использованное - кол-во начисленных токенов * скорость генерации одного токена
// но если достигнут capacity - двигаем время на текущее
func (s *Service) refill(capacity int64, perSecond int64, tokens int64, lastRefill time.Time) (int64, time.Time) {
	nanosToGenerateOneToken := (time.Second / time.Duration(perSecond)).Nanoseconds()
	now := time.Now().UTC()

	duration := now.Sub(lastRefill).Nanoseconds()
	tokensSinceLastRefill := duration / nanosToGenerateOneToken

	if tokensSinceLastRefill > 0 {
		tokens = min(capacity, tokens+tokensSinceLastRefill)
		lastRefill = lastRefill.Add(time.Duration(tokensSinceLastRefill * nanosToGenerateOneToken))
		if tokens == capacity {
			lastRefill = now
		}
	}

	return tokens, lastRefill
}
