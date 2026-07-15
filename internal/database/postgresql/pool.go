package postgresql

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"
	"unicode"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const actorSetting = "atlas.actor_id"

// Config contains non-secret PostgreSQL pool behavior and the externally
// delivered connection string. The URL must come from runtime secret delivery,
// not a committed configuration file.
type Config struct {
	URL               string
	ApplicationName   string
	MaxConnections    int32
	MinConnections    int32
	ConnectTimeout    time.Duration
	MaxConnectionAge  time.Duration
	MaxConnectionIdle time.Duration
	HealthInterval    time.Duration
	StatementTimeout  time.Duration
	LockTimeout       time.Duration
	IdleInTxTimeout   time.Duration
}

func (c Config) validate() error {
	if strings.TrimSpace(c.URL) == "" {
		return errors.New("database URL is required")
	}
	decodedURL, decodeErr := url.QueryUnescape(c.URL)
	if strings.Contains(strings.ToLower(c.URL), actorSetting) ||
		(decodeErr == nil && strings.Contains(strings.ToLower(decodedURL), actorSetting)) {
		return errors.New("database URL must not set atlas.actor_id")
	}
	if strings.TrimSpace(c.ApplicationName) == "" {
		return errors.New("database application name is required")
	}
	if c.MaxConnections < 1 {
		return errors.New("maximum database connections must be at least one")
	}
	if c.MinConnections < 0 || c.MinConnections > c.MaxConnections {
		return errors.New("minimum database connections must be between zero and the maximum")
	}
	for name, value := range map[string]time.Duration{
		"connect timeout":              c.ConnectTimeout,
		"maximum connection age":       c.MaxConnectionAge,
		"maximum connection idle time": c.MaxConnectionIdle,
		"health interval":              c.HealthInterval,
		"statement timeout":            c.StatementTimeout,
		"lock timeout":                 c.LockTimeout,
		"idle-in-transaction timeout":  c.IdleInTxTimeout,
	} {
		if value <= 0 {
			return fmt.Errorf("%s must be positive", name)
		}
	}
	return nil
}

// Pool owns the least-privileged application connection pool.
type Pool struct {
	pool *pgxpool.Pool
}

func Open(ctx context.Context, cfg Config) (*Pool, error) {
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	poolCfg, err := pgxpool.ParseConfig(cfg.URL)
	if err != nil {
		return nil, errors.New("parse database configuration: invalid connection string")
	}
	if err := rejectActorRuntimeParameters(poolCfg.ConnConfig.RuntimeParams); err != nil {
		return nil, err
	}
	poolCfg.MaxConns = cfg.MaxConnections
	poolCfg.MinConns = cfg.MinConnections
	poolCfg.MaxConnLifetime = cfg.MaxConnectionAge
	poolCfg.MaxConnIdleTime = cfg.MaxConnectionIdle
	poolCfg.HealthCheckPeriod = cfg.HealthInterval
	poolCfg.ConnConfig.ConnectTimeout = cfg.ConnectTimeout

	if poolCfg.ConnConfig.RuntimeParams == nil {
		poolCfg.ConnConfig.RuntimeParams = make(map[string]string)
	}
	if err := rejectActorRuntimeParameters(poolCfg.ConnConfig.RuntimeParams); err != nil {
		return nil, err
	}
	poolCfg.ConnConfig.RuntimeParams["application_name"] = cfg.ApplicationName
	poolCfg.ConnConfig.RuntimeParams["statement_timeout"] = durationMilliseconds(cfg.StatementTimeout)
	poolCfg.ConnConfig.RuntimeParams["lock_timeout"] = durationMilliseconds(cfg.LockTimeout)
	poolCfg.ConnConfig.RuntimeParams["idle_in_transaction_session_timeout"] = durationMilliseconds(cfg.IdleInTxTimeout)

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("create database pool: %w", err)
	}
	result := &Pool{pool: pool}
	if err := result.Check(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	return result, nil
}

func rejectActorRuntimeParameters(params map[string]string) error {
	for name, value := range params {
		candidate := strings.ToLower(name + "=" + value)
		if strings.Contains(candidate, actorSetting) {
			return errors.New("database configuration must not provide session actor context")
		}
	}
	return nil
}

func durationMilliseconds(value time.Duration) string {
	if value <= 0 {
		return "0"
	}
	return fmt.Sprintf("%dms", value.Milliseconds())
}

func (p *Pool) Name() string { return "postgresql" }

func (p *Pool) Check(ctx context.Context) error {
	if p == nil || p.pool == nil {
		return errors.New("database pool is not initialized")
	}
	if err := p.pool.Ping(ctx); err != nil {
		return fmt.Errorf("database dependency unavailable: %w", err)
	}
	return nil
}

func (p *Pool) Close() {
	if p != nil && p.pool != nil {
		p.pool.Close()
	}
}

// WithActor creates a transaction, binds the authenticated actor with a
// transaction-local PostgreSQL setting, executes the callback, and commits.
// The actor setting is never applied at session scope and therefore cannot be
// intentionally retained on a pooled connection.
func (p *Pool) WithActor(
	ctx context.Context,
	actorID string,
	fn func(context.Context, pgx.Tx) error,
) error {
	if p == nil || p.pool == nil {
		return errors.New("database pool is not initialized")
	}
	actorID = strings.TrimSpace(actorID)
	if err := validateActorID(actorID); err != nil {
		return err
	}
	if fn == nil {
		return errors.New("transaction callback is required")
	}

	tx, err := p.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return fmt.Errorf("begin governed transaction: %w", err)
	}
	defer func() {
		rollbackCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = tx.Rollback(rollbackCtx)
	}()

	var boundActor string
	if err := tx.QueryRow(
		ctx,
		"SELECT set_config($1, $2, true)",
		actorSetting,
		actorID,
	).Scan(&boundActor); err != nil {
		return fmt.Errorf("bind transaction actor: %w", err)
	}
	if boundActor != actorID {
		return errors.New("database did not bind the requested transaction actor")
	}

	if err := fn(ctx, tx); err != nil {
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit governed transaction: %w", err)
	}
	return nil
}

func validateActorID(actorID string) error {
	if actorID == "" {
		return errors.New("authenticated actor is required")
	}
	if len(actorID) > 256 {
		return errors.New("authenticated actor identifier is too long")
	}
	if strings.IndexFunc(actorID, unicode.IsControl) >= 0 {
		return errors.New("authenticated actor identifier contains control characters")
	}
	return nil
}

func (p *Pool) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return p.pool.Query(ctx, sql, args...)
}

func (p *Pool) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return p.pool.QueryRow(ctx, sql, args...)
}
