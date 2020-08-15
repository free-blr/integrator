package bot

import (
	"context"

	"free.blr/integrator/internal/model"
)

type tagService interface {
	GetByOptions(context.Context, model.TagQueryOptions) ([]*model.Tag, error)
}

type requestService interface {
	Insert(context.Context, ...*model.Request) error
}
