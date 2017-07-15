package utils

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
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
		return "", errors.New("VIRTUALGO_PATH environment variable is not set")
	}
	return path, nil
}
