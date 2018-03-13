package utils

import (
	"syscall"
	"testing"

	"github.com/GetStream/vg/internal/testsuite"
)

type Suite struct {
	testsuite.TestSuite
}

func TestUtils(t *testing.T) {
	testsuite.Run(t, &Suite{})
}

func (s *Suite) home() string {
	home, found := syscall.Getenv("HOME")
	s.True(found)
	return home
}

func (s *Suite) TestReplaceHomeDir() {
	testcases := []struct{ in, out string }{
		{"/somewhere", "/somewhere"},
		{"/somewhere/~/else", "/somewhere/~/else"},
		{"~/.virtualgo", s.home() + "/.virtualgo"},
	}

	for _, tc := range testcases {
		s.Run(tc.in, func() {
			s.Equal(tc.out, ReplaceHomeDir(tc.in))
		})
	}
}

func (s *Suite) TestVirtualgoRoot() {
	expected := s.home() + "/.virtualgo"
	s.Equal(expected, VirtualgoRoot())

	// Make sure directory is created
	exists, err := DirExists(expected)
	s.NoError(err)
	s.True(exists)
}
