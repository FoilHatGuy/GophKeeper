package database

import (
	"context"
	"errors"
	"gophKeeper/src/server/cfg"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	ErrConflict = errors.New("this data is already stored")
	ErrNotFound = errors.New("query returned empty result")
)

type CategoryHead []*struct {
	DataID   string
	Metadata string
}

type StorageController interface {
	Initialise(ctx context.Context, config *cfg.ConfigT) (err error)
	AddUser(ctx context.Context, login, password string) (err error)
	GetPassword(ctx context.Context, login string) (password string, err error)
	AddSession(ctx context.Context, login, sid string) (err error)
	GetSession(ctx context.Context, login string) (sid string, err error)
	UpdateSession(ctx context.Context, login, sid string) (err error)
	RefreshSession(ctx context.Context, sid string) (uid string, ok bool, err error)
	AddCredentials(ctx context.Context, dataID, metadata string, data []byte) (err error)
	GetCredentials(ctx context.Context, dataID string) (metadata string, data []byte, err error)
	GetCredentialsHead(ctx context.Context) (head CategoryHead, err error)
}

func New(ctx context.Context, config *cfg.ConfigT) (ctrl StorageController) {
	ctrl = &storageWrapper{}
	err := ctrl.Initialise(ctx, config)
	if err != nil {
		logrus.Fatalf("database was not initialised: %v", err)
		return nil
	}
	return ctrl
}

type storageWrapper struct {
	PSQL *gorm.DB
}

func (s *storageWrapper) Initialise(ctx context.Context, config *cfg.ConfigT) (err error) {
	db, err := gorm.Open(postgres.Open(config.Data.PostgesDSN), &gorm.Config{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(
		&Session{},
		&User{},
		&SecureCredential{},
		&SecureText{},
		&SecureCard{},
		&SecureFile{},
	)
	if err != nil {
		return err
	}
	s.PSQL = db
	return nil
}

func (s *storageWrapper) AddUser(ctx context.Context, login, password string) (err error) {
	logrus.Debug("PSQL added login password", login, password)
	return nil
}

func (s *storageWrapper) GetPassword(ctx context.Context, login string) (password string, err error) {
	logrus.Debug("PSQL loaded login password", login)
	if login == "" {
		return "", ErrNotFound
	}
	return "password", nil
}

func (s *storageWrapper) AddSession(ctx context.Context, login, sid string) (err error) {
	if len(login) > 0 {
		logrus.Debug("PSQL loaded session for user", login, sid)
		return nil
	}
	logrus.Debug("PSQL loaded session for user", login)
	return ErrConflict
}

func (s *storageWrapper) GetSession(ctx context.Context, login string) (sid string, err error) {
	logrus.Debug("PSQL added new session for user", login, sid)
	return "err", nil
}

func (s *storageWrapper) UpdateSession(ctx context.Context, login, sid string) (err error) {
	logrus.Debug("PSQL updated session for user", login, sid)
	return nil
}

func (s *storageWrapper) RefreshSession(ctx context.Context, sid string) (uid string, ok bool, err error) {
	logrus.Debug("PSQL refreshed session", sid)
	return "uid", true, nil
}

func (s *storageWrapper) GetCredentialsHead(ctx context.Context) (head CategoryHead, err error) {
	logrus.Debug("PSQL loaded data for login pass pair")
	return CategoryHead{}, nil
}

func (s *storageWrapper) AddCredentials(ctx context.Context, dataID, metadata string, data []byte) (err error) {
	logrus.Debug("PSQL added data for login pass pair", dataID, metadata)
	return nil
}

func (s *storageWrapper) GetCredentials(ctx context.Context, dataID string) (metadata string, data []byte, err error) {
	logrus.Debug("PSQL loaded data for login pass pair", dataID)
	return "meta", []byte("password"), nil
}
