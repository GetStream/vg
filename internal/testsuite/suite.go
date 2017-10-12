package testsuite

import "github.com/stretchr/testify/suite"

type TestSuite struct {
	suite.Suite
}

// Alias suite.Run so suite doesn't have to be imported
var Run = suite.Run

// NoError uses Require by default, because you almost always want an instant
// failure if an error occurs.
func (s *TestSuite) NoError(err error, msgAndArgs ...interface{}) {
	s.Require().NoError(err, msgAndArgs...)
}

// Len uses Require by default, because this is usually checked before indexing
// an array. Using Assert would result in a out of bounds panic, so Require is
// a better default.
func (s *TestSuite) Len(object interface{}, length int, msgAndArgs ...interface{}) {
	s.Require().Len(object, length, msgAndArgs...)
}

// Contains uses Require by default, because this is usually checked before
// indexing a map. Using Assert would result in a key not in bounds panic, so
// Require is a better default.
func (s *TestSuite) Contains(container interface{}, thingToFind interface{}, msgAndArgs ...interface{}) {
	s.Require().Contains(container, thingToFind, msgAndArgs...)
}

// NotNil uses Require by default, because this is usually checked before
// accessing attributes.
// Using Assert would then result in a nil dereference panic, so Require is
// a better default.
func (s *TestSuite) NotNil(object interface{}, msgAndArgs ...interface{}) {
	s.Require().NotNil(object, msgAndArgs...)
}
