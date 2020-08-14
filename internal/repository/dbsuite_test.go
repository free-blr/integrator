package repository_test

import (
	"free.blr/integrator/testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/suite"

	_ "github.com/lib/pq"
)

type DBSuite struct {
	db      *testing.DB
	builder sq.StatementBuilderType
	suite.Suite
}

func (s *DBSuite) SetupSuite() {
	var err error

	s.db, err = testing.NewDB()
	s.Require().NoError(err)

	s.builder = sq.StatementBuilder.PlaceholderFormat(sq.Dollar).RunWith(s.db)

	testing.MigrateDownDb(s.db.DB)
	testing.MustMigrateDb(s.db.DB)
}

func (s *DBSuite) TearDownSuite() {
	err := s.db.Close()
	s.Require().NoError(err)
}
