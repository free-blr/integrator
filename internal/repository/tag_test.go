package repository_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"free.blr/integrator/internal/model"
	"free.blr/integrator/internal/repository"
)

type TagSuite struct {
	DBSuite
	repo *repository.Tag
}

func TestTagRepository(t *testing.T) {
	suite.Run(t, new(TagSuite))
}

func (s *TagSuite) SetupTest() {
	s.repo = repository.NewTag(s.db)
}

func (s *TagSuite) TearDownTest() {

}

func (s *TagSuite) TestTagRepository_GetByOptions() {
	_, err := s.builder.
		Insert("tag").
		Columns("id", "name").
		Values(1, "Автопомощь").
		Values(2, "Хирургия").
		Exec()
	s.Require().NoError(err)

	s.Run("get tag without options", func() {
		expected := []*model.Tag{{
			ID:   1,
			Name: "Автопомощь",
		}, {
			ID:   2,
			Name: "Хирургия",
		},
		}

		rs, err := s.repo.GetByOptions(context.Background(), model.TagQueryOptions{})
		s.Require().NoError(err)
		s.Require().Len(rs, len(expected))

		for i, r := range rs {
			s.Equal(expected[i], r)
		}
	})

}
