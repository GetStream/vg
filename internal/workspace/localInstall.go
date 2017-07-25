package workspace

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/pkg/errors"
)

func (ws *Workspace) InstallLocalPackage(pkg string, localPath string) error {
	_, _ = fmt.Fprintf(os.Stderr, "Installing local sources at %q in workspace as %q\n", localPath, pkg)
	pkgDir := filepath.Join(path.Split(pkg))
	linkName := filepath.Join(ws.Src(), pkgDir)

	err := os.MkdirAll(filepath.Dir(linkName), 0755)
	if err != nil {
		return errors.WithStack(err)
	}

	err = os.RemoveAll(linkName)
	if err != nil {
		return errors.WithStack(err)
	}

	return errors.WithStack(os.Symlink(localPath, linkName))
}

func (ws *Workspace) InstallPersistentLocalPackages() error {
	settings, err := ws.Settings()
	if err != nil {
		return err
	}

	for pkg, install := range settings.LocalInstalls {
		err := ws.InstallLocalPackage(pkg, install.Path)
		if err != nil {
			return err
		}
	}
	return nil
}
