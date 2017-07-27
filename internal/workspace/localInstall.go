package workspace

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/GetStream/vg/internal/utils"
	"github.com/pkg/errors"
)

func (ws *Workspace) InstallLocalPackage(pkg string, localPath string) error {
	err := ws.Uninstall(pkg, os.Stderr)
	if err != nil {
		return err
	}

	_, _ = fmt.Fprintf(os.Stderr, "Installing local sources at %q in workspace as %q\n", localPath, pkg)
	pkgDir := filepath.Join(path.Split(pkg))
	linkName := filepath.Join(ws.Src(), pkgDir)

	err = os.MkdirAll(filepath.Dir(linkName), 0755)
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

func (ws *Workspace) Uninstall(pkg string, logWriter io.Writer) error {
	pkgDir := utils.PkgToDir(pkg)
	err := os.RemoveAll(filepath.Join(ws.Src(), pkgDir))
	if err != nil {
		return errors.Wrapf(err, "Couldn't remove package src '%s'", ws.Name())
	}

	pkgInstalledDirs, err := filepath.Glob(filepath.Join(ws.Pkg(), "*", pkgDir))
	if err != nil {
		return errors.Wrapf(err, "Something went wrong when globbing files for '%s'", pkg)
	}

	for _, path := range pkgInstalledDirs {
		fmt.Fprintln(logWriter, "Removing", path)

		err = os.RemoveAll(path)
		if err != nil {
			return errors.Wrapf(err, "Couldn't remove installed package files for '%s'", pkg)
		}
	}

	pkgInstalledFiles, err := filepath.Glob(filepath.Join(ws.Pkg(), "*", pkgDir+".a"))
	if err != nil {
		return errors.Wrapf(err, "Something went wrong when globbing files for '%s'", pkg)
	}

	for _, path := range pkgInstalledFiles {
		fmt.Fprintln(logWriter, "Removing", path)

		err = os.RemoveAll(path)
		if err != nil {
			return errors.Wrapf(err, "Couldn't remove installed package files for '%s'", pkg)
		}
	}

	settings, err := ws.Settings()
	if err != nil {
		return err
	}

	if _, ok := settings.LocalInstalls[pkg]; ok {
		fmt.Fprintf(logWriter, "Removing %q from persistent local installs\n", pkg)
		delete(settings.LocalInstalls, pkg)

		err = ws.SaveSettings(settings)
		if err != nil {
			return err
		}
	}

	return nil
}
