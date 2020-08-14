package repository

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"free.blr/integrator/internal/model"
)

type Request struct {
	pgRepository
}

func NewRequest(db sqlx.ExtContext) *Request {
	return &Request{pgRepository: newPgRepository(db)}
}

func (r *Request) GetByOptions(ctx context.Context, opts model.RequestQueryOptions) (_ []*model.Request, err error) {
	qb := r.baseQuery().OrderBy("b.created_at DESC")

	if len(opts.Type) != 0 {
		qb = qb.Where(sq.Eq{"r.type": opts.Type})
	}
	if len(opts.TagID) != 0 {
		qb = qb.Where(sq.Eq{"r.tag_id": opts.TagID})
	}

	return r.selects(ctx, qb)
}

func (r *Request) selects(ctx context.Context, builder sq.SelectBuilder) ([]*model.Request, error) {
	query, args, err := builder.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "could not build query")
	}

	var rs []*model.Request
	if err := sqlx.SelectContext(ctx, r.queryer(ctx), &rs, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "could not select requests")
	}
	return rs, nil
}

func (r *Request) baseQuery() sq.SelectBuilder {
	return r.
		queryBuilder().
		Select(
			`r.id`,
			`r.type`,
			`r.tg_user_id`,
			`t.id as tag.id`,
			`t.name as tag.name`,
		).
		From("request AS r").
		Join("tag t on t.id = r.tag_id")
}
