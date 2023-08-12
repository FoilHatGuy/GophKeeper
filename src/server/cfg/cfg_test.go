package cfg

import (
	"encoding/json"
	"os"
	"strconv"
	"testing"

	"github.com/mcuadros/go-defaults"
	"github.com/stretchr/testify/suite"
)

type ConfigTestSuite struct {
	suite.Suite
}

func (s *ConfigTestSuite) TestNew() {
	config1 := New(FromDefaults())

	config2 := &ConfigT{
		Server: &ServerT{},
		Data:   &DataStorageT{},
	}
	defaults.SetDefaults(config2.Server)
	defaults.SetDefaults(config2.Data)

	s.Assert().Equal(*config1.Server, *config2.Server)
	s.Assert().Equal(*config1.Data, *config2.Data)
}

func (s *ConfigTestSuite) TestWithBuild() {
	source := &BuildT{
		BuildVersion: "BuildVersion",
		BuildDate:    "BuildDate",
		BuildCommit:  "BuildCommit",
	}
	config1 := New(
		WithBuild(source),
	)

	s.Assert().Equal(config1.Build, source)
}

func (s *ConfigTestSuite) TestFromEnv() {
	t := s.T()

	const (
		address      = "ADDRESS"
		hTTPS        = true
		loggingLevel = "LOGGING_LEVEL"
		sessionLife  = 100
		fileSavePath = "FILE_SAVE_PATH"
		postgesDSN   = "POSTGES_DSN"
	)
	t.Setenv("SERVER_ADDRESS", address)
	t.Setenv("HTTPS", strconv.FormatBool(hTTPS))
	t.Setenv("LOGGING_LEVEL", loggingLevel)
	t.Setenv("SESSION_LIFE", strconv.Itoa(sessionLife))
	t.Setenv("FILE_SAVE_PATH", fileSavePath)
	t.Setenv("POSTGES_DSN", postgesDSN)

	config2 := New(FromDefaults(),
		FromEnv(),
	)
	s.Assert().Equal(address, config2.Server.Address)
	s.Assert().Equal(hTTPS, config2.Server.HTTPS)
	s.Assert().Equal(loggingLevel, config2.Server.LoggingLevel)
	s.Assert().Equal(sessionLife, config2.Server.SessionLife)
	s.Assert().Equal(fileSavePath, config2.Data.FileSavePath)
	s.Assert().Equal(postgesDSN, config2.Data.PostgesDSN)
}

func (s *ConfigTestSuite) TestFromJSONFile() {
	const filePath = "./test.json"
	origin := &ConfigT{
		Build: &BuildT{
			BuildVersion: "BuildVersion",
			BuildDate:    "BuildDate",
			BuildCommit:  "BuildCommit",
		},
		Server: &ServerT{
			Address:      "Address",
			HTTPS:        false,
			LoggingLevel: "LoggingLevel",
			SessionLife:  230,
		},
		Data: &DataStorageT{
			FileSavePath: "FileSavePath",
			PostgesDSN:   "PostgesDSN",
		},
	}

	file, _ := json.MarshalIndent(origin, "", "\t")

	_ = os.WriteFile(filePath, file, 0o600)
	defer func() {
		err := os.Remove(filePath)
		s.Assert().NoError(err)
	}()

	data, err := os.ReadFile(filePath)
	s.Assert().NoError(err)
	s.Assert().NotEmpty(data)

	t := s.T()
	t.Setenv("CONFIG", filePath)

	config1 := New(
		FromJSON(),
	)

	s.Assert().Equal(*config1.Server, *origin.Server)
	s.Assert().Equal(*config1.Data, *origin.Data)

	// broken file
	file, _ = json.MarshalIndent(origin, "\"", "\"")
	//nolint:gosec
	_ = os.WriteFile(filePath, file, 0o300)
	New(
		FromJSON(),
	)

	t.Setenv("CONFIG", "")
	New(
		FromJSON(),
	) // cause an error
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}
