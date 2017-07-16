package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

type WorkspaceSettings struct {
	GlobalFallback bool                    `toml:"global-fallback"`
	LocalInstalls  map[string]LocalInstall `toml:"local-install"`
}

type LocalInstall struct {
	Path string `toml:"path"`
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

func SaveSettings(workspace string, settings *WorkspaceSettings) error {
	settingsPath := SettingsPath(workspace)

	file, err := os.OpenFile(settingsPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return errors.WithStack(err)
	}

	err = toml.NewEncoder(file).Encode(settings)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func SaveCurrentSettings(settings *WorkspaceSettings) error {
	workspace, err := CurrentWorkspace()
	if err != nil {
		return err
	}

	return SaveSettings(workspace, settings)
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

func LinkLocalInstalls(workspace string, settings *WorkspaceSettings) error {
	for pkg, install := range settings.LocalInstalls {
		fmt.Printf("Linking %q sources locally to %q\n", pkg, install.Path)
		pkgDir := filepath.Join(path.Split(pkg))
		linkName := filepath.Join(SrcDir(workspace), pkgDir)

		err := os.MkdirAll(filepath.Dir(linkName), 0755)
		if err != nil {
			return errors.WithStack(err)
		}

		err = os.RemoveAll(linkName)
		if err != nil {
			return errors.WithStack(err)
		}
		err = os.Symlink(install.Path, linkName)
		if err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

func LinkCurrentLocalInstalls() error {
	workspace, err := CurrentWorkspace()
	if err != nil {
		return err
	}

	settings, err := CurrentSettings()
	if err != nil {
		return err
	}

	return LinkLocalInstalls(workspace, settings)

}
