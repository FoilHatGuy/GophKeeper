//go:build unit

package passwords

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

type HashTestSuite struct {
	suite.Suite
}

func (s *HashTestSuite) TestHash() {
	pw := "TestHash"
	fmt.Printf("pw: %q\n", pw)
	hashed, err := HashPassword(pw)
	s.Assert().NoError(err)
	fmt.Printf("hashed pw: %q\n", hashed)

	ok, err := ComparePasswordHash(hashed, pw)
	s.Assert().NoError(err)
	s.Assert().True(ok)
}

func (s *HashTestSuite) TestWrongPW() {
	pw := "TestWrongPW"
	fmt.Printf("pw: %q\n", pw)
	hashed, err := HashPassword(pw)
	s.Assert().NoError(err)
	fmt.Printf("hashed pw: %q\n", hashed)

	pw2 := "other password"
	ok, err := ComparePasswordHash(hashed, pw2)
	s.Assert().NoError(err)
	s.Assert().False(ok)
}

func (s *HashTestSuite) TestFaultyEncoding() {
	pw := "TestFaultyEncoding"
	fmt.Printf("pw: %q\n", pw)
	hashed, err := HashPassword(pw)
	s.Assert().NoError(err)
	fmt.Printf("hashed pw: %q\n", hashed)

	incorrectlyHashed := "incorrectly hashed password"
	ok, err := ComparePasswordHash(incorrectlyHashed, pw)
	s.Assert().Error(err)
	s.Assert().False(ok)
}

func (s *HashTestSuite) TestBadArgonVer() {
	pw := "TestBadArgonVer"
	hashed := "$argon2id$v=17$m=65536,t=2,p=2$n+QasX1OlriHZY9FtzWtsw$vxMgsOXDiLQo28r69qYqNYJmxMTBcr6IzJGz2C+7cJM"
	//                  ^v=19 - correct version
	fmt.Printf("pw: %q\n", pw)
	fmt.Printf("hashed pw: %q\n", hashed)

	ok, err := ComparePasswordHash(hashed, pw)
	s.Assert().ErrorIs(err, ErrArgonVersion)
	s.Assert().False(ok)

	hashed2 := "$argon2id$a=19$m=65536,t=2,p=2$n+QasX1OlriHZY9FtzWtsw$vxMgsOXDiLQo28r69qYqNYJmxMTBcr6IzJGz2C+7cJM"
	//                   ^v=19 - correct version
	fmt.Printf("pw: %q\n", pw)
	fmt.Printf("hashed pw: %q\n", hashed2)

	ok, err = ComparePasswordHash(hashed2, pw)
	s.Assert().ErrorContains(err, "input does not match format")
	s.Assert().False(ok)
}

func (s *HashTestSuite) TestBadHashStructure() {
	// nolint:gosec
	pw := "TestBadHashStructure"
	hashed := "$argon2id$v=19$s=65536,v=2,a=2$n+QasX1OlriHZY9FtzWtsw$vxMgsOXDiLQo28r69qYqNYJmxMTBcr6IzJGz2C+7cJM"
	//                      ^$m=65536,t=2,p=2 - correct config
	fmt.Printf("pw: %q\n", pw)
	fmt.Printf("hashed pw: %q\n", hashed)

	ok, err := ComparePasswordHash(hashed, pw)
	s.Assert().ErrorContains(err, "input does not match format")
	s.Assert().False(ok)
}

func (s *HashTestSuite) TestBadBase64Encode() {
	// nolint:gosec
	pw := "TestBadBase64Encode"
	hashed := "$argon2id$v=19$m=65536,t=2,p=2$%SignIsNotAllowedInEnc$vxMgsOXDiLQo28r69qYqNYJmxMTBcr6IzJGz2C+7cJM"
	//                correct salt base64 -> ^n+QasX1OlriHZY9FtzWtsw
	fmt.Printf("pw: %q\n", pw)
	fmt.Printf("hashed pw: %q\n", hashed)

	ok, err := ComparePasswordHash(hashed, pw)
	s.Assert().ErrorContains(err, "illegal base64 data at input byte")
	s.Assert().False(ok)

	hashed2 := "$argon2id$v=19$m=65536,t=2,p=2$n+QasX1OlriHZY9FtzWtsw$%SignIsNotAllowedInBase64EncodingAsWellAs//"
	//                                    correct password base64 -> ^vxMgsOXDiLQo28r69qYqNYJmxMTBcr6IzJGz2C+7cJM
	fmt.Printf("pw: %q\n", pw)
	fmt.Printf("hashed pw: %q\n", hashed2)

	ok, err = ComparePasswordHash(hashed2, pw)
	s.Assert().ErrorContains(err, "illegal base64 data at input byte")
	s.Assert().False(ok)
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(HashTestSuite))
}
