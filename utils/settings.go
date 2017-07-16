package utils

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

type WorkspaceSettings struct {
	GlobalFallback bool                    `toml:"global-fallback"`
	LocalInstalls  map[string]LocalInstall `toml:"local-install"`
}

type LocalInstall struct {
	Path string
}

func NewWorkspaceSettings() *WorkspaceSettings {
	return &WorkspaceSettings{
		LocalInstalls: make(map[string]LocalInstall),
	}
}

func DefaultWorkspaceSettings() *WorkspaceSettings {
	settings := NewWorkspaceSettings()
	settings.GlobalFallback = true
	return settings
}

const settingsFile = "virtualgo.toml"

// Settings returns the settings for a specific workspace.
func Settings(workspace string) (*WorkspaceSettings, error) {
	settings := DefaultWorkspaceSettings()
	settingsBytes, err := ioutil.ReadFile(SettingsPath(workspace))
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

func CurrentSettings() (*WorkspaceSettings, error) {
	workspace, err := CurrentWorkspace()
	if err != nil {
		return nil, err
	}

	return Settings(workspace)
}

func SettingsPath(workspace string) string {
	return filepath.Join(WorkspaceDir(workspace), settingsFile)

}

func CurrentSettingsPath() (string, error) {
	dir, err := CurrentWorkspaceDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, settingsFile), nil
}
