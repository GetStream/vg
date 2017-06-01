package cmd

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

func replaceHomeDir(path string) (string, error) {
	if path[:2] != "~/" {
		return path, nil
	}

	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	homeDir := usr.HomeDir
	return filepath.Join(homeDir, path[2:]), nil
}

var virtualgoDir string

func init() {
	var err error
	virtualgoDir, err = replaceHomeDir("~/.virtualgo")
	if err != nil {
		panic(fmt.Sprintf("Couldn't get path of virtualgo directory: %v", err))
	}

	err = os.MkdirAll(virtualgoDir, 0755)
	if err != nil {
		panic(fmt.Sprintf("Couldn't create virtualgo directory: %v", err))
	}
}
