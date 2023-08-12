//go:build integration

package grpcclient

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type DatabaseIntegrationTestSuite struct {
	suite.Suite
}

func (s *DatabaseIntegrationTestSuite) SetupSuite() {
}

func (s *DatabaseIntegrationTestSuite) TestNew() {
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(DatabaseIntegrationTestSuite))
}
