//go:build unit

package cfg

import (
	"encoding/json"
	"flag"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/suite"
)

type CfgTestSuite struct {
	suite.Suite
}

func (s *CfgTestSuite) TestNew() {
	config := New()
	s.Assert().Nil(config.Build)
}

func (s *CfgTestSuite) TestNewFromDefaults() {
	config := New(FromDefaults())
	s.Assert().NotNil(config.Build)

	s.Assert().Equal(reflect.TypeOf(*config).Field(0).Tag.Get("default"), config.SecretPath)
	s.Assert().Equal(reflect.TypeOf(*config).Field(1).Tag.Get("default"), config.ServerAddress)
}

func (s *CfgTestSuite) TestNewFromJSON() {
	config := New(FromJSON())

	const filePath = "./test.json"
	err := flag.Set("c", configPath)
	s.Assert().NoError(err)

	origin := ConfigT{
		SecretPath:    "1",
		ServerAddress: "2",
	}

	file, _ := json.MarshalIndent(origin, "", "\t")

	_ = os.WriteFile(filePath, file, 0o600)
	defer func() {
		err := os.Remove(filePath)
		s.Assert().NoError(err)
	}()

	err = flag.Set("c", filePath)
	s.Assert().NoError(err)

	data, err := os.ReadFile(configPath)
	s.Assert().NoError(err)
	s.Assert().NotEmpty(data)

	config = New(FromJSON())
	s.Assert().Nil(config.Build)

	s.Assert().Equal(origin.ServerAddress, config.ServerAddress)
	s.Assert().Equal(origin.SecretPath, config.SecretPath)

	_ = os.WriteFile(filePath, []byte("}}}}}"), 0o600)

	config = New(FromJSON())
	s.Assert().Nil(config.Build)
	s.Assert().Equal("", config.ServerAddress)
	s.Assert().Equal("", config.SecretPath)
}

func (s *CfgTestSuite) TestNewFromEnv() {
	config := New()
	s.Assert().Nil(config.Build)
}

func (s *CfgTestSuite) TestNewFromFlags() {
	const (
		secretPath    = "secretPath"
		serverAddress = "serverAddress"
	)
	err := flag.Set("k", secretPath)
	s.Assert().NoError(err)
	err = flag.Set("a", serverAddress)
	s.Assert().NoError(err)

	config := New(FromFlags())
	s.Assert().Equal(config.SecretPath, secretPath)
	s.Assert().Equal(config.ServerAddress, serverAddress)

	// do stuff twice
	config = New(FromFlags())
	s.Assert().Equal(config.SecretPath, secretPath)
	s.Assert().Equal(config.ServerAddress, serverAddress)
}

func (s *CfgTestSuite) TestNewWithBuild() {
	build := &BuildT{
		BuildVersion: "BuildVersion",
		BuildDate:    "BuildDate",
		BuildCommit:  "BuildCommit",
	}
	config := New(WithBuild(build))
	s.Assert().Equal(build, config.Build)
}

func TestClientConfigUnit(t *testing.T) {
	suite.Run(t, new(CfgTestSuite))
}
