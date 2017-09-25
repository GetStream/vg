package utils

import (
	"fmt"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/pkg/errors"
)

func ReplaceHomeDir(path string) string {
	if path[:2] != "~/" {
		return path
	}

	usr, err := user.Current()
	if err != nil {
		panic(fmt.Sprintf("Couldn't get the current user: %v", err))
	}
	homeDir := usr.HomeDir
	return filepath.Join(homeDir, path[2:])
}

func VirtualgoRoot() string {
	dir := os.Getenv("VIRTUALGO_ROOT")
	if dir == "" {
		dir = ReplaceHomeDir("~/.virtualgo")
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		panic(fmt.Sprintf("Couldn't create virtualgo directory: %v", err))
	}
	return dir
}

func PkgToDir(pkg string) string {
	return filepath.Join(strings.Split(pkg, "/")...)

}

func DirToPkg(dir string) string {
	return path.Join(strings.Split(dir, string(os.PathSeparator))...)
}

func OriginalGopath() string {
	gopath := os.Getenv("_VIRTUALGO_OLDGOPATH")
	if gopath == "" {
		return CurrentGopath()
	}

	return gopath
}

// exists returns whether the given file or directory exists or not
func DirExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err == nil {
		if !info.IsDir() {
			return false, errors.Errorf("%q is not a directory", path)
		}
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, errors.Wrapf(err, "error occured when checking if directory %q exists", path)
}

func VendorExists() (bool, error) {
	return DirExists("vendor")
}

func CurrentGopath() string {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		return defaultGOPATH()
	}

	return gopath
}

// Taken from https://github.com/golang/go/blob/go1.8/src/go/build/build.go#L260-L277
func defaultGOPATH() string {
	env := "HOME"
	if runtime.GOOS == "windows" {
		env = "USERPROFILE"
	} else if runtime.GOOS == "plan9" {
		env = "home"
	}
	if home := os.Getenv(env); home != "" {
		def := filepath.Join(home, "go")
		if filepath.Clean(def) == filepath.Clean(runtime.GOROOT()) {
			// Don't set the default GOPATH to GOROOT,
			// as that will trigger warnings from the go tool.
			return ""
		}
		return def
	}
	return ""
}
