//go:build unit

package application

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/suite"
)

type UtilsTestSuite struct {
	suite.Suite
}

func (s *UtilsTestSuite) TestIncludes() {
	list := []string{"one", "two"}

	ok := includes(list, "one")
	s.Assert().True(ok)

	ok = includes(list, "not_ok")
	s.Assert().False(ok)
}

func (s *UtilsTestSuite) TestFirstN() {
	const inputString = "12345678901234567890123456789012345678901234567890"

	newLen := rand.Intn(len(inputString))

	newString1 := firstN(inputString, newLen)
	s.Assert().Equal(newLen, len(newString1))

	newString2 := firstN(inputString, len(inputString)+10)
	s.Assert().Equal(len(inputString), len(newString2))
}

func TestUtils(t *testing.T) {
	suite.Run(t, new(UtilsTestSuite))
}
