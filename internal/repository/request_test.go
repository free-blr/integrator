package repository_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"free.blr/integrator/internal/model"
	"free.blr/integrator/internal/repository"
)

type RequestSuite struct {
	DBSuite
	repo *repository.Request
}

func TestRequestRepository(t *testing.T) {
	suite.Run(t, new(RequestSuite))
}

func (s *RequestSuite) SetupTest() {
	s.repo = repository.NewRequest(s.db)
	_, err := s.builder.
		Insert("tag").
		Columns("id", "name").
		Values(1, "Автопомощь").
		Values(2, "Хирургия").
		Exec()
	s.NoError(err)
}

func (s *RequestSuite) TearDownTest() {
	_ = s.db.MustExec(`TRUNCATE TABLE "tag" CASCADE;`)
}

func (s *RequestSuite) TestRequestRepository_GetByOptions() {
	_, err := s.builder.
		Insert("request").
		Columns("id", "type", "tg_user_id", "tag_id").
		Values(1, model.RequestTypeIn, "11", 1).
		Values(2, model.RequestTypeIn, "12", 2).
		Values(3, model.RequestTypeOut, "13", 1).
		Values(4, model.RequestTypeOut, "14", 2).
		Values(5, model.RequestTypeOut, "15", 2).
		Exec()
	s.Require().NoError(err)

	s.Run("get requests without options", func() {
		expected := []int{1, 2, 3, 4, 5}

		rs, err := s.repo.GetByOptions(context.Background(), model.RequestQueryOptions{})
		s.Require().NoError(err)
		s.Require().Len(rs, len(expected))

		for i, r := range rs {
			s.Equal(expected[i], r.ID)
		}
	})

	s.Run("get requests type", func() {
		expected := []*model.Request{{
			ID:       1,
			Type:     model.RequestTypeIn,
			TgUserID: "11",
			Tag: model.Tag{
				ID:   1,
				Name: "Автопомощь",
			},
		}, {
			ID:       2,
			Type:     model.RequestTypeIn,
			TgUserID: "12",
			Tag: model.Tag{
				ID:   2,
				Name: "Хирургия",
			},
		}}

		rss, err := s.repo.GetByOptions(context.Background(), model.RequestQueryOptions{
			Type: []model.RequestType{model.RequestTypeIn},
		})
		s.Require().NoError(err)
		s.Require().Len(rss, len(expected))

		for i, rs := range rss {
			s.Equal(expected[i], rs)
		}
	})
}

func (s *RequestSuite) TestRequestRepository_Insert() {
	s.Run("insert", func() {
		expected := []*model.Request{{
			ID:       1,
			Type:     model.RequestTypeIn,
			TgUserID: "11",
			Tag: model.Tag{
				ID:   1,
				Name: "Автопомощь",
			},
			TagID: 1,
		}, {
			ID:       2,
			Type:     model.RequestTypeIn,
			TgUserID: "12",
			Tag: model.Tag{
				ID:   2,
				Name: "Хирургия",
			},
			TagID: 2,
		}}

		err := s.repo.Insert(context.Background(), expected...)
		s.Require().NoError(err)

		rss, err := s.repo.GetByOptions(context.Background(), model.RequestQueryOptions{})
		s.Require().NoError(err)
		s.Require().Len(rss, len(expected))

		for i, rs := range rss {
			expected[i].TagID = 0
			s.Equal(expected[i], rs)
		}
	})
}
