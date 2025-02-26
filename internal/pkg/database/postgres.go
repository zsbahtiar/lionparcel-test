package database

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgres struct {
	*pgxpool.Pool
}

type Postgres interface {
	Acquire(ctx context.Context) (*pgxpool.Conn, error)
	AcquireAllIdle(ctx context.Context) []*pgxpool.Conn
	AcquireFunc(ctx context.Context, f func(*pgxpool.Conn) error) error
	Begin(ctx context.Context) (pgx.Tx, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
	Close()
	Config() *pgxpool.Config
	CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Ping(ctx context.Context) error
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Reset()
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
	Stat() *pgxpool.Stat

	Select(ctx context.Context, dest any, sql string, args ...any) error
	SelectOne(ctx context.Context, dest any, sql string, args ...any) error
	Rebind(query string) string
}

func NewPostgres(dbUser, dbPassword, dbHost, dbPort, dbName string) Postgres {
	dsn := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(dbUser, dbPassword),
		Host:   fmt.Sprintf("%s:%s", dbHost, dbPort),
		Path:   dbName,
	}

	q := dsn.Query()
	q.Add("sslmode", "disable")

	dsn.RawQuery = q.Encode()

	var config *pgxpool.Config
	config, err := pgxpool.ParseConfig(dsn.String())
	if err != nil {
		panic(err)
	}
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		panic(err)
	}

	return &postgres{
		Pool: pool,
	}
}

func (p postgres) Select(ctx context.Context, dest any, sql string, args ...any) error {
	rows, err := p.Pool.Query(ctx, sql, args...)
	if err != nil {
		return err
	}

	return pgxscan.ScanAll(dest, rows)
}

func (p postgres) SelectOne(ctx context.Context, dest any, sql string, args ...any) error {
	rows, err := p.Pool.Query(ctx, sql, args...)
	if err != nil {
		return err
	}
	if err := pgxscan.ScanOne(dest, rows); err != nil {
		if pgxscan.NotFound(err) {
			return pgx.ErrNoRows
		}
		return err
	}
	return nil
}

func (p postgres) Rebind(query string) string {
	rqb := make([]byte, 0, len(query)+10)

	var i, j int

	for i = strings.Index(query, "?"); i != -1; i = strings.Index(query, "?") {
		rqb = append(rqb, query[:i]...)

		rqb = append(rqb, '$')

		j++
		rqb = strconv.AppendInt(rqb, int64(j), 10)

		query = query[i+1:]
	}

	return string(append(rqb, query...))
}
