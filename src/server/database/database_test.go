//go:build integration

package database

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/sakirsensoy/genv/dotenv"
	"github.com/stretchr/testify/suite"

	"gophKeeper/src/server/cfg"
)

type DatabaseIntegrationTestSuite struct {
	suite.Suite
	ctx     context.Context
	wrapper StorageController
}

func (s *DatabaseIntegrationTestSuite) SetupSuite() {
	err := dotenv.Load("../../../.env")
	if err != nil {
		panic("error in config")
	}

	config := cfg.New(
		cfg.FromDefaults(),
		cfg.FromEnv(),
	)
	config.Data.PostgesDSN = os.ExpandEnv("host=localhost user=${POSTGRES_USER} " +
		"password=${POSTGRES_PASSWORD} dbname=postgres port=${PGPORT} sslmode=disable")
	config.Server.SessionLife = 200
	s.ctx = context.Background()
	s.wrapper = New(s.ctx, config)
}

func (s *DatabaseIntegrationTestSuite) TestNew() {
	config := cfg.New(
		cfg.FromDefaults(),
		cfg.FromEnv(),
	)
	config.Data.PostgesDSN = os.ExpandEnv("host=wrong user=${POSTGRES_USER} " +
		"password=${POSTGRES_PASSWORD} dbname=dont_exist port=${PGPORT} sslmode=disable")

	s.Panics(func() {
		s.wrapper = New(s.ctx, config)
	})

	config.Data.PostgesDSN = os.ExpandEnv("host=localhost user=${POSTGRES_USER} " +
		"password=${POSTGRES_PASSWORD} dbname=postgres port=${PGPORT} sslmode=disable")
	s.wrapper = New(context.Background(), config)
}

func (s *DatabaseIntegrationTestSuite) TestUsers() {
	uid := uuid.NewString()
	login := uuid.NewString()
	password := uuid.NewString()
	err := s.wrapper.AddUser(s.ctx, uid, login, password)
	s.Assert().NoError(err)

	// add user twice
	err = s.wrapper.AddUser(s.ctx, uid, login, password)
	s.Assert().Error(err)

	_, pw2, err := s.wrapper.GetUserData(s.ctx, login)
	s.Assert().NoError(err)
	s.Assert().Equal(password, pw2)
}

func (s *DatabaseIntegrationTestSuite) TestSessions() {
	uid := uuid.NewString()
	login := uuid.NewString()
	password := uuid.NewString()
	err := s.wrapper.AddUser(s.ctx, uid, login, password)
	s.Assert().NoError(err)

	sid := uuid.NewString()
	err = s.wrapper.AddSession(s.ctx, uid, sid)
	s.Assert().NoError(err)

	// add user twice
	uid2, ok, err := s.wrapper.RefreshSession(s.ctx, sid)
	s.Assert().NoError(err)
	s.Assert().True(ok)
	s.Assert().Equal(uid, uid2)

	sid2 := uuid.NewString()
	err = s.wrapper.UpdateSession(s.ctx, uid, sid2)
	s.Assert().NoError(err)
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(DatabaseIntegrationTestSuite))
}
