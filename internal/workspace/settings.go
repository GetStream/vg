package workspace

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

type Settings struct {
	GlobalFallback bool                     `toml:"global-fallback"`
	LocalInstalls  map[string]*localInstall `toml:"local-install"`
}

type localInstall struct {
	Path       string `toml:"path"`
	Persistent bool   `toml:"persistent"`
	Successful bool   `toml:"successful"`
	Bindfs     bool   `toml:"bindfs"`
}

func NewSettings() *Settings {
	return &Settings{
		LocalInstalls: make(map[string]*localInstall),
	}
}

func DefaultSettings() *Settings {
	settings := NewSettings()
	settings.GlobalFallback = true
	return settings
}

const settingsFile = "virtualgo.toml"

// Settings returns the settings of the workspace
func (ws *Workspace) Settings() (*Settings, error) {
	if ws.settings != nil {
		return ws.settings, nil
	}
	return ws.LoadSettings()
}

func (ws *Workspace) LoadSettings() (*Settings, error) {
	ws.settings = DefaultSettings()
	settingsBytes, err := ioutil.ReadFile(ws.SettingsPath())
	if err != nil {
		if os.IsNotExist(err) {
			return ws.settings, nil
		}
		return nil, errors.Wrap(err, "Something went wrong loading the settings file")
	}
	err = toml.Unmarshal(settingsBytes, ws.settings)
	if err != nil {
		return nil, errors.Wrap(err, "Couldn't load settings from settings file")
	}

	return ws.settings, nil
}

func (ws *Workspace) SaveSettings(settings *Settings) error {
	ws.settings = settings

	file, err := os.OpenFile(ws.SettingsPath(), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return errors.WithStack(err)
	}

	err = toml.NewEncoder(file).Encode(settings)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (ws *Workspace) SettingsPath() string {
	return filepath.Join(ws.Path(), settingsFile)

}
