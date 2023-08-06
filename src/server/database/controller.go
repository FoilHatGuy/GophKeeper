package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm/clause"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"gophKeeper/src/server/cfg"
)

var (
	ErrConflict     = errors.New("this data is already stored")
	ErrNotFound     = errors.New("query returned empty result")
	ErrSessionStale = errors.New("session is already expired")
)

// CategoryHead contains DataID's and metadata of all ds in category,
// excluding UID and actual data
type CategoryHead []*struct {
	DataID   string
	Metadata string
}

// StorageController is an interface for interaction with database.
// Can be implemented by other type and substituted in server
type StorageController interface {
	Initialise(ctx context.Context, config *cfg.ConfigT) (err error)

	AddUser(ctx context.Context, uid, login, password string) (err error)
	GetPassword(ctx context.Context, login string) (password string, err error)

	AddSession(ctx context.Context, uid, sid string) (err error)
	UpdateSession(ctx context.Context, uid, sid string) (err error)
	RefreshSession(ctx context.Context, sid string) (uid string, ok bool, err error)

	GetCredentialsHead(ctx context.Context, uid string) (head CategoryHead, err error)
	AddCredentials(ctx context.Context, uid, dataID, metadata string, data []byte) (err error)
	GetCredentials(ctx context.Context, uid, dataID string) (metadata string, data []byte, err error)

	GetTextHead(ctx context.Context, uid string) (head CategoryHead, err error)
	AddText(ctx context.Context, uid, dataID, metadata string, data []byte) (err error)
	GetText(ctx context.Context, uid, dataID string) (metadata string, data []byte, err error)

	GetCardHead(ctx context.Context, uid string) (head CategoryHead, err error)
	AddCard(ctx context.Context, uid, dataID, metadata string, data []byte) (err error)
	GetCard(ctx context.Context, uid, dataID string) (metadata string, data []byte, err error)
}

// New returns a new instance of database controller
func New(ctx context.Context, config *cfg.ConfigT) (ctrl StorageController) {
	ctrl = &storageWrapper{}
	err := ctrl.Initialise(ctx, config)
	if err != nil {
		logrus.Panicf("database was not initialised: %v", err)
		return nil
	}
	return ctrl
}

type storageWrapper struct {
	PSQL *gorm.DB
	conf *cfg.ConfigT
}

// Initialise operates with database using GORM
func (s *storageWrapper) Initialise(ctx context.Context, config *cfg.ConfigT) (err error) {
	db, err := gorm.Open(postgres.Open(config.Data.PostgesDSN), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("open connection: %w", err)
	}
	err = db.WithContext(ctx).AutoMigrate(
		&Session{},
		&User{},
		&SecureCredential{},
		&SecureText{},
		&SecureCard{},
		&SecureFile{},
	)
	if err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}
	s.PSQL = db
	s.conf = config
	return nil
}

// AddUser operates with database using GORM
func (s *storageWrapper) AddUser(ctx context.Context, uid, login, password string) (err error) {
	err = s.PSQL.WithContext(ctx).Create(&User{
		ID:       uid,
		Login:    login,
		Password: password,
	}).Error
	if err != nil {
		return fmt.Errorf("user adding failed: %w", err)
	}
	return nil
}

// GetPassword operates with database using GORM
func (s *storageWrapper) GetPassword(ctx context.Context, login string) (password string, err error) {
	err = s.PSQL.WithContext(ctx).Model(&User{}).Where("login = ?", login).Pluck("password", &password).Error
	if err != nil {
		return password, fmt.Errorf("user pw get: %w", err)
	}
	return password, nil
}

// AddSession operates with database using GORM
func (s *storageWrapper) AddSession(ctx context.Context, uid, sid string) (err error) {
	err = s.PSQL.WithContext(ctx).Create(&Session{
		ID:      sid,
		UID:     uid,
		Expires: time.Now().Add(time.Duration(s.conf.Server.SessionLife) * time.Second),
	}).Error
	if err != nil {
		return fmt.Errorf("session add: %w", err)
	}
	return nil
}

// UpdateSession operates with database using GORM
func (s *storageWrapper) UpdateSession(ctx context.Context, uid, sid string) (err error) {
	err = s.PSQL.WithContext(ctx).
		Where("id = ?", sid).
		Updates(Session{
			UID:     uid,
			Expires: time.Now().Add(time.Duration(s.conf.Server.SessionLife) * time.Second),
		}).
		Error
	if err != nil {
		return fmt.Errorf("session update: %w", err)
	}
	return nil
}

// RefreshSession operates with database using GORM
func (s *storageWrapper) RefreshSession(ctx context.Context, sid string) (uid string, ok bool, err error) {
	var currentSession Session
	var expTime time.Time

	op := s.PSQL.WithContext(ctx).
		Model(&currentSession).
		Clauses(clause.Returning{
			Columns: []clause.Column{
				{Name: "uid"},
				{Name: "expires"},
			},
		}).
		Where("id = ?", sid).
		Pluck("expires", &expTime).
		Where("expires < ?", time.Now()).
		Updates(Session{
			Expires: time.Now().Add(time.Duration(s.conf.Server.SessionLife) * time.Second),
		})
	ok = op.RowsAffected > 0
	err = op.Error
	if err != nil {
		return currentSession.UID, false, fmt.Errorf("session refresh: %w", err)
	}
	if currentSession.Expires.After(time.Now()) {
		return currentSession.UID, false, ErrSessionStale
	}
	return currentSession.UID, ok, nil
}

// credentials section

// GetCredentialsHead operates with database using GORM
func (s *storageWrapper) GetCredentialsHead(ctx context.Context, uid string) (head CategoryHead, err error) {
	op := s.PSQL.Model(&User{}).WithContext(ctx).Where("uid = ?", uid)
	err = op.Error
	logrus.Debug("PSQL loaded data for login pass pair")
	if err != nil {
		return CategoryHead{}, fmt.Errorf("credentials head get: %w", err)
	}
	return CategoryHead{}, nil
}

// AddCredentials operates with database using GORM
func (s *storageWrapper) AddCredentials(ctx context.Context, uid, dataID, metadata string, data []byte) (err error) {
	err = s.PSQL.WithContext(ctx).Create(SecureCredential{
		ID:       dataID,
		Data:     data,
		Metadata: metadata,
		UID:      uid,
	}).Error
	logrus.Debug("PSQL added data for login pass pair", dataID)
	if err != nil {
		return fmt.Errorf("credentials add: %w", err)
	}
	return nil
}

// GetCredentials operates with database using GORM
func (s *storageWrapper) GetCredentials(
	ctx context.Context,
	uid, dataID string,
) (metadata string, data []byte, err error) {
	err = s.PSQL.
		WithContext(ctx).
		Model(&SecureCredential{}).
		Where("uid = ?", uid).
		Where("id = ?", dataID).
		Pluck("data", &data).
		Pluck("metadata", &metadata).
		Error
	logrus.Debug("PSQL loaded data for login pass pair", dataID)
	if err != nil {
		return metadata, data, fmt.Errorf("credentials get: %w", err)
	}
	return metadata, data, nil
}

// Text section

// GetTextHead operates with database using GORM
func (s *storageWrapper) GetTextHead(ctx context.Context, uid string) (head CategoryHead, err error) {
	op := s.PSQL.Model(&User{}).WithContext(ctx).Where("uid = ?", uid)
	err = op.Error
	logrus.Debug("PSQL loaded data for login pass pair")
	if err != nil {
		return CategoryHead{}, fmt.Errorf("text head get: %w", err)
	}
	return CategoryHead{}, nil
}

// AddText operates with database using GORM
func (s *storageWrapper) AddText(ctx context.Context, uid, dataID, metadata string, data []byte) (err error) {
	err = s.PSQL.WithContext(ctx).Create(SecureText{
		ID:       dataID,
		Data:     data,
		Metadata: metadata,
		UID:      uid,
	}).Error
	logrus.Debug("PSQL added data for login pass pair", dataID)
	if err != nil {
		return fmt.Errorf("text add: %w", err)
	}
	return nil
}

// GetText operates with database using GORM
func (s *storageWrapper) GetText(
	ctx context.Context,
	uid, dataID string,
) (metadata string, data []byte, err error) {
	err = s.PSQL.
		WithContext(ctx).
		Model(&SecureText{}).
		Where("uid = ?", uid).
		Where("id = ?", dataID).
		Pluck("data", &data).
		Pluck("metadata", &metadata).
		Error
	logrus.Debug("PSQL loaded data for login pass pair", dataID)
	if err != nil {
		return metadata, data, fmt.Errorf("text get: %w", err)
	}
	return metadata, data, nil
}

// Card section

// GetCardHead operates with database using GORM
func (s *storageWrapper) GetCardHead(ctx context.Context, uid string) (head CategoryHead, err error) {
	op := s.PSQL.Model(&User{}).WithContext(ctx).Where("uid = ?", uid)
	err = op.Error
	logrus.Debug("PSQL loaded data for login pass pair")
	if err != nil {
		return CategoryHead{}, fmt.Errorf("card head get: %w", err)
	}
	return CategoryHead{}, nil
}

// AddCard operates with database using GORM
func (s *storageWrapper) AddCard(ctx context.Context, uid, dataID, metadata string, data []byte) (err error) {
	err = s.PSQL.WithContext(ctx).Create(SecureCard{
		ID:       dataID,
		Data:     data,
		Metadata: metadata,
		UID:      uid,
	}).Error
	logrus.Debug("PSQL added data for login pass pair", dataID)
	if err != nil {
		return fmt.Errorf("card add: %w", err)
	}
	return nil
}

// GetCard operates with database using GORM
func (s *storageWrapper) GetCard(
	ctx context.Context,
	uid, dataID string,
) (metadata string, data []byte, err error) {
	err = s.PSQL.
		WithContext(ctx).
		Model(&SecureCard{}).
		Where("uid = ?", uid).
		Where("id = ?", dataID).
		Pluck("data", &data).
		Pluck("metadata", &metadata).
		Error
	logrus.Debug("PSQL loaded data for login pass pair", dataID)
	if err != nil {
		return metadata, data, fmt.Errorf("card get: %w", err)
	}
	return metadata, data, nil
}
