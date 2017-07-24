package utils

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/pkg/errors"
)

func InstallLocalPackage(workspace string, pkg string, localPath string) error {
	_, _ = fmt.Fprintf(os.Stderr, "Installing local sources at %q in workspace as %q\n", localPath, pkg)
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

	return errors.WithStack(os.Symlink(localPath, linkName))
}

func InstallCurrentLocalPackage(pkg string, localPath string) error {
	workspace, err := CurrentWorkspace()
	if err != nil {
		return err
	}
	return InstallLocalPackage(workspace, pkg, localPath)
}

func InstallPersistentLocalPackages(workspace string, settings *WorkspaceSettings) error {
	for pkg, install := range settings.LocalInstalls {
		err := InstallLocalPackage(workspace, pkg, install.Path)
		if err != nil {
			return err
		}
	}
	return nil
}

func InstallCurrentPersistentLocalPackages() error {
	workspace, err := CurrentWorkspace()
	if err != nil {
		return err
	}

	settings, err := CurrentSettings()
	if err != nil {
		return err
	}

	return InstallPersistentLocalPackages(workspace, settings)

}
