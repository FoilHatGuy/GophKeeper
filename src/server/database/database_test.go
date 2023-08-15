//go:build integration

package database

import (
	"context"
	"os"
	"testing"
	"time"

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
	config.Server.SessionLife = 10
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
	wrapper := New(context.Background(), config)
	s.Assert().IsType(&storageWrapper{}, wrapper)
	s.Assert().Implements((*StorageController)(nil), wrapper)
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

	// get nonexistent session
	otherLogin := uuid.NewString()
	_, _, err = s.wrapper.GetUserData(s.ctx, otherLogin)
	s.Assert().Error(err)
}

func (s *DatabaseIntegrationTestSuite) TestSessions() {
	uid := uuid.NewString()
	login := uuid.NewString()
	password := uuid.NewString()
	err := s.wrapper.AddUser(s.ctx, uid, login, password)
	s.Assert().NoError(err)

	time.Sleep(1 * time.Second)
	sid := uuid.NewString()
	err = s.wrapper.AddSession(s.ctx, uid, sid)
	s.Assert().NoError(err)
	err = s.wrapper.AddSession(s.ctx, uid, sid)
	s.Assert().ErrorIs(err, ErrConflict)

	time.Sleep(1 * time.Second)
	uid2, ok, err := s.wrapper.RefreshSession(s.ctx, sid)
	s.Assert().NoError(err)
	s.Assert().True(ok)
	s.Assert().Equal(uid, uid2)

	time.Sleep(1 * time.Second)
	sid3 := uuid.NewString()
	_, ok, err = s.wrapper.RefreshSession(s.ctx, sid3)
	s.Assert().Error(err)
	s.Assert().False(ok)

	time.Sleep(1 * time.Second)
	sid2 := uuid.NewString()
	err = s.wrapper.UpdateSession(s.ctx, uid, sid2)
	s.Assert().NoError(err)

	// add twice
	time.Sleep(1 * time.Second)
	uidOther := uuid.NewString()
	err = s.wrapper.AddSession(s.ctx, uidOther, sid)
	s.Assert().Error(err)

	time.Sleep(1 * time.Second)
	err = s.wrapper.UpdateSession(s.ctx, uidOther, sid2)
	s.Assert().NoError(err)

	time.Sleep(11 * time.Second)
	_, ok, err = s.wrapper.RefreshSession(s.ctx, sid2)
	s.Assert().ErrorIs(err, ErrSessionStale)
	s.Assert().False(ok)
}

func (s *DatabaseIntegrationTestSuite) TestCredentials() {
	uid := uuid.NewString()
	login := uuid.NewString()
	password := uuid.NewString()
	err := s.wrapper.AddUser(s.ctx, uid, login, password)
	s.Assert().NoError(err)

	time.Sleep(1 * time.Second)
	sid := uuid.NewString()
	err = s.wrapper.AddSession(s.ctx, uid, sid)
	s.Assert().NoError(err)

	time.Sleep(1 * time.Second)
	dataID := uuid.NewString()
	data := []byte(uuid.NewString())
	metadata := uuid.NewString()

	err = s.wrapper.AddCredentials(s.ctx, uid, dataID, metadata, data)
	s.Assert().NoError(err)
	// add twice
	err = s.wrapper.AddCredentials(s.ctx, uid, dataID, metadata, data)
	s.Assert().Error(err)

	metadata2, data2, err := s.wrapper.GetCredentials(s.ctx, uid, dataID)
	s.Assert().NoError(err)
	s.Assert().Equal(data, data2)
	s.Assert().Equal(metadata, metadata2)

	// get nonexistent
	dataID2 := uuid.NewString()
	_, _, err = s.wrapper.GetCredentials(s.ctx, uid, dataID2)
	s.Assert().Error(err)

	head, err := s.wrapper.GetCredentialsHead(s.ctx, uid)
	s.Assert().NoError(err)
	s.Assert().Equal(1, len(head))
	s.Assert().Equal(metadata, head[0].Metadata)
	s.Assert().Equal(dataID, head[0].ID)
}

func (s *DatabaseIntegrationTestSuite) TestCard() {
	uid := uuid.NewString()
	login := uuid.NewString()
	password := uuid.NewString()
	err := s.wrapper.AddUser(s.ctx, uid, login, password)
	s.Assert().NoError(err)

	time.Sleep(1 * time.Second)
	sid := uuid.NewString()
	err = s.wrapper.AddSession(s.ctx, uid, sid)
	s.Assert().NoError(err)

	time.Sleep(1 * time.Second)
	dataID := uuid.NewString()
	data := []byte(uuid.NewString())
	metadata := uuid.NewString()

	err = s.wrapper.AddCard(s.ctx, uid, dataID, metadata, data)
	s.Assert().NoError(err)
	// add twice
	err = s.wrapper.AddCard(s.ctx, uid, dataID, metadata, data)
	s.Assert().Error(err)

	metadata2, data2, err := s.wrapper.GetCard(s.ctx, uid, dataID)
	s.Assert().NoError(err)
	s.Assert().Equal(data, data2)
	s.Assert().Equal(metadata, metadata2)

	// get nonexistent
	dataID2 := uuid.NewString()
	_, _, err = s.wrapper.GetCard(s.ctx, uid, dataID2)
	s.Assert().Error(err)

	head, err := s.wrapper.GetCardHead(s.ctx, uid)
	s.Assert().NoError(err)
	s.Assert().Equal(1, len(head))
	s.Assert().Equal(metadata, head[0].Metadata)
	s.Assert().Equal(dataID, head[0].ID)
}

func (s *DatabaseIntegrationTestSuite) TestText() {
	uid := uuid.NewString()
	login := uuid.NewString()
	password := uuid.NewString()
	err := s.wrapper.AddUser(s.ctx, uid, login, password)
	s.Assert().NoError(err)

	time.Sleep(1 * time.Second)
	sid := uuid.NewString()
	err = s.wrapper.AddSession(s.ctx, uid, sid)
	s.Assert().NoError(err)

	time.Sleep(1 * time.Second)
	dataID := uuid.NewString()
	data := []byte(uuid.NewString())
	metadata := uuid.NewString()

	err = s.wrapper.AddText(s.ctx, uid, dataID, metadata, data)
	s.Assert().NoError(err)
	// add twice
	err = s.wrapper.AddText(s.ctx, uid, dataID, metadata, data)
	s.Assert().Error(err)

	metadata2, data2, err := s.wrapper.GetText(s.ctx, uid, dataID)
	s.Assert().NoError(err)
	s.Assert().Equal(data, data2)
	s.Assert().Equal(metadata, metadata2)

	// get nonexistent
	dataID2 := uuid.NewString()
	_, _, err = s.wrapper.GetText(s.ctx, uid, dataID2)
	s.Assert().Error(err)

	head, err := s.wrapper.GetTextHead(s.ctx, uid)
	s.Assert().NoError(err)
	s.Assert().Equal(1, len(head))
	s.Assert().Equal(metadata, head[0].Metadata)
	s.Assert().Equal(dataID, head[0].ID)
}

func TestDatabaseIntegration(t *testing.T) {
	suite.Run(t, new(DatabaseIntegrationTestSuite))
}
