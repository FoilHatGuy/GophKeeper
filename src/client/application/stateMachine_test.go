//go:build unit

package application

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type StateMachineTestSuite struct {
	suite.Suite
}

func TestStateMachine(t *testing.T) {
	suite.Run(t, new(StateMachineTestSuite))
}
