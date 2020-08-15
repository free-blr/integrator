package repository

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"free.blr/integrator/internal/model"
)

type Tag struct {
	pgRepository
}

func NewTag(db sqlx.ExtContext) *Tag {
	return &Tag{pgRepository: newPgRepository(db)}
}

func (r *Tag) GetByOptions(ctx context.Context, _ model.TagQueryOptions) (_ []*model.Tag, err error) {
	qb := r.baseQuery()
	return r.selects(ctx, qb)
}

func (r *Tag) selects(ctx context.Context, builder sq.SelectBuilder) ([]*model.Tag, error) {
	query, args, err := builder.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "could not build query")
	}

	var rs []*model.Tag
	if err := sqlx.SelectContext(ctx, r.queryer(ctx), &rs, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "could not select tags")
	}
	return rs, nil
}

func (r *Tag) baseQuery() sq.SelectBuilder {
	return r.
		queryBuilder().
		Select(
			`t.id`,
			`t.name`,
		).
		From("tag AS t")
}
