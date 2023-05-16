package storage

import (
	"context"
	"fmt"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/iho/bitly/conf"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

type postgres struct{ pool *pgxpool.Pool }

func New(c *conf.Config) Storage {
	pool, err := pgxpool.New(
		context.Background(),
		fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s",
			c.Host, c.Port, c.User, c.Password, c.DbName,
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	return &postgres{pool: pool}
}

func (p postgres) Save(ctx context.Context, key, url string) error {

	sql, args, err := psql.Insert("urls").Columns("id", "url").Values(key, url).ToSql()
	if err != nil {
		return err
	}

	p.pool.QueryRow(ctx, sql, args...)
	return nil
}

type URL struct {
	ID  string `db:"id"`
	URL string `db:"url"`
}

func (p postgres) Load(ctx context.Context, key string) (string, error) {
	sql, args, err := psql.Select("*").From("urls").Where(sq.Eq{"id": key}).Limit(1).ToSql()
	if err != nil {
		return "", err
	}

	row := p.pool.QueryRow(ctx, sql, args...)
	if err != nil {
		return "", err
	}

	res := new(URL)
	err = row.Scan(&res.ID, &res.URL)
	if err != nil {
		return "", err
	}

	return res.URL, nil
}

// Module for go fx
var Module = fx.Options(fx.Provide(New))
