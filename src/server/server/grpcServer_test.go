//go:build integration

package app

import (
	"testing"

	"github.com/stretchr/testify/suite"

	pb "gophKeeper/src/pb"
)

type DatabaseIntegrationTestSuite struct {
	suite.Suite
	serverKeep pb.GophKeeperServer
	serverAuth pb.AuthServer
}

func (s *DatabaseIntegrationTestSuite) SetupSuite() {
}

func (s *DatabaseIntegrationTestSuite) TestNew() {
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(DatabaseIntegrationTestSuite))
}
