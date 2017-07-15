package utils

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

type WorkspaceSettings struct {
	FullyIsolated bool
}

const settingsFile = "virtualgo.toml"

// Settings returns the settings for a specific workspace.
func Settings(workspace string) (*WorkspaceSettings, error) {
	settings := &WorkspaceSettings{FullyIsolated: false}
	settingsBytes, err := ioutil.ReadFile(filepath.Join(WorkspaceDir(workspace), settingsFile))
	if err != nil {
		if os.IsNotExist(err) {
			return settings, nil
		}
		return nil, errors.Wrap(err, "Something went wrong loading the settings file")
	}
	err = toml.Unmarshal(settingsBytes, settings)
	if err != nil {
		return nil, errors.Wrap(err, "Couldn't load settings from settings file")
	}

	return settings, nil
}
