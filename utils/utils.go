package utils

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
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

func VirtualgoDir() string {
	var err error
	dir := ReplaceHomeDir("~/.virtualgo")

	err = os.MkdirAll(dir, 0755)
	if err != nil {
		panic(fmt.Sprintf("Couldn't create virtualgo directory: %v", err))
	}
	return dir
}

func WorkspaceDir(workspace string) string {
	return filepath.Join(VirtualgoDir(), workspace)
}

func CurrentWorkspaceDir() (string, error) {
	path := os.Getenv("VIRTUALGO_PATH")
	if path == "" {
		return "", errors.New("A virtualgo workspace should be active first by using `vg activate [workspaceName]`")
	}
	return path, nil
}

func SrcDir(workspace string) string {
	return filepath.Join(WorkspaceDir(workspace), "src")
}

func CurrentSrcDir() (string, error) {
	workspaceDir, err := CurrentWorkspaceDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(workspaceDir, "src"), nil
}

func CurrentWorkspace() (string, error) {
	workspace := os.Getenv("VIRTUALGO")
	if workspace == "" {
		return "", errors.New("A virtualgo workspace should be active first by using `vg activate [workspaceName]`")
	}
	return workspace, nil

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
