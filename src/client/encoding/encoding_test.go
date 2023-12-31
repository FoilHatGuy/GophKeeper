package encoding

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

type HashTestSuite struct {
	suite.Suite

	encoder *Encoder
}

func (s *HashTestSuite) SetupSuite() {
	s.encoder = New("new secret")
}

func (s *HashTestSuite) TestHash() {
	in := "TestHash"
	fmt.Printf("input: %q\n", in)
	result := s.encoder.Encode(in)
	fmt.Printf("result: % X\n", result)

	out, err := s.encoder.Decode(result)
	s.Assert().NoError(err)
	s.Assert().Equal(in, out)
}

func (s *HashTestSuite) TestWrongDecoder() {
	in := "TestWrongDecoder"
	fmt.Printf("input: %q\n", in)
	result := s.encoder.Encode(in)
	fmt.Printf("result: % x\n", result)

	wrongDecoder := New("other secret")
	out, err := wrongDecoder.Decode(result)
	fmt.Printf("output: %q\n", out)
	s.Assert().ErrorIs(err, ErrWrongKey)
}

func TestBlockEncodingUnit(t *testing.T) {
	suite.Run(t, new(HashTestSuite))
}
