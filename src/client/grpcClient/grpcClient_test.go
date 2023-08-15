//go:build unit

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

func TestGrPCClientUnit(t *testing.T) {
	suite.Run(t, new(DatabaseIntegrationTestSuite))
}
