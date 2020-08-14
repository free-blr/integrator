package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type pgRepository struct {
	db sqlx.ExtContext
}

func newPgRepository(db sqlx.ExtContext) pgRepository {
	return pgRepository{db: db}
}

func (pg pgRepository) queryBuilder() sq.StatementBuilderType {
	return sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
}

func (pg pgRepository) queryer(_ context.Context) sqlx.ExtContext {
	return pg.db
}
