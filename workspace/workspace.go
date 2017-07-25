package workspace

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/GetStream/vg/utils"
)

type Workspace struct {
	name     string
	path     string
	settings *Settings
}

func New(name string) *Workspace {
	return &Workspace{
		name: name,
		path: filepath.Join(utils.VirtualgoRoot(), name),
	}
}

func Current() (*Workspace, error) {
	name := os.Getenv("VIRTUALGO")
	if name == "" {
		return nil, errors.New("A virtualgo workspace should be active first by using `vg activate [workspaceName]`")
	}

	return New(name), nil

}

func (ws *Workspace) Name() string {
	return ws.name
}

func (ws *Workspace) Path() string {
	return ws.path
}

func (ws *Workspace) Src() string {
	return filepath.Join(ws.Path(), "src")
}

func (ws *Workspace) Pkg() string {
	return filepath.Join(ws.Path(), "pkg")
}
