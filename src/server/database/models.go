package database

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	ID       string `gorm:"primaryKey"`
	Login    string
	Password string
}

type Session struct {
	gorm.Model
	ID      string
	UID     string `gorm:"uniqueIndex"`
	User    User   `gorm:"ForeignKey:UID;references:ID"`
	Expires time.Time
}

type SecureCredential struct {
	gorm.Model
	ID       string
	Data     string
	Metadata string
	UID      string
	User     User `gorm:"ForeignKey:UID;references:ID"`
}

type SecureText struct {
	gorm.Model
	ID       string
	Data     string
	Metadata string
	UID      string
	User     User `gorm:"ForeignKey:UID;references:ID"`
}

type SecureCard struct {
	gorm.Model
	ID       string
	Data     string
	Metadata string
	UID      string
	User     User `gorm:"ForeignKey:UID;references:ID"`
}

type SecureFile struct {
	gorm.Model
	ID       string
	Filename string
	Metadata string
	UID      string
	User     User `gorm:"ForeignKey:UID;references:ID"`
}
