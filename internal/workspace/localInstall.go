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

func (ws *Workspace) InstallLocalPackagePersistently(pkg string, localPath string) error {
	err := ws.InstallLocalPackage(pkg, localPath)
	if err != nil {
		return err
	}

	settings, err := ws.Settings()
	if err != nil {
		return err
	}

	settings.LocalInstalls[pkg].Persistent = true

	fmt.Fprintf(os.Stderr, "Persisting the local install for %q\n", pkg)
	return ws.SaveSettings(settings)
}

func (ws *Workspace) InstallLocalPackage(pkg string, localPath string) error {
	pkgDir := filepath.Join(path.Split(pkg))
	target := filepath.Join(ws.Src(), pkgDir)

	err := ws.Uninstall(pkg, os.Stderr)
	if err != nil {
		return err
	}

	settings, err := ws.Settings()
	if err != nil {
		return err
	}

	settings.LocalInstalls[pkg] = &localInstall{
		Path: localPath,
	}

	err = ws.SaveSettings(settings)
	if err != nil {
		return err
	}

	err = ws.installLocalPackageWithSymlink(pkg, localPath, target)
	if err != nil {
		return err
	}

	settings.LocalInstalls[pkg].Successful = true
	return ws.SaveSettings(settings)

}

func (ws *Workspace) installLocalPackageWithSymlink(pkg, src, target string) error {
	_, _ = fmt.Fprintf(os.Stderr, "Installing local sources at %q in workspace as %q\n", src, pkg)

	err := os.MkdirAll(filepath.Dir(target), 0755)
	if err != nil {
		return errors.WithStack(err)
	}

	err = os.RemoveAll(target)
	if err != nil {
		return errors.WithStack(err)
	}

	return errors.WithStack(os.Symlink(src, target))
}

func (ws *Workspace) InstallPersistentLocalPackages() error {
	settings, err := ws.Settings()
	if err != nil {
		return err
	}

	for pkg, install := range settings.LocalInstalls {
		if !install.Persistent {
			continue
		}
		err := ws.InstallLocalPackagePersistently(pkg, install.Path)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ws *Workspace) Uninstall(pkg string, logWriter io.Writer) error {
	pkgDir := utils.PkgToDir(pkg)
	pkgSrc := filepath.Join(ws.Src(), pkgDir)
	_, err := os.Stat(pkgSrc)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return errors.WithStack(err)
	}

	fmt.Fprintf(logWriter, "Uninstalling %q from workspace\n", pkg)
	fmt.Fprintf(logWriter, "  Removing sources at %q\n", pkgSrc)
	err = os.RemoveAll(pkgSrc)
	if err != nil {
		return errors.Wrapf(err, "Couldn't remove package src %q", ws.Name())
	}

	// Remove possible LocalInstall entry
	settings, err := ws.Settings()
	if err != nil {
		return err
	}
	install, ok := settings.LocalInstalls[pkg]
	if ok && !install.Persistent {
		fmt.Fprintf(logWriter, "  Removing %q from locally installed packages\n", pkg)
		delete(settings.LocalInstalls, pkg)
	}

	err = ws.SaveSettings(settings)
	if err != nil {
		return err
	}
	pkgInstalledDirs, err := filepath.Glob(filepath.Join(ws.Pkg(), "*", pkgDir))
	if err != nil {
		return errors.Wrapf(err, "Something went wrong when globbing files for %q", pkg)
	}

	for _, path := range pkgInstalledDirs {
		fmt.Fprintf(logWriter, "  Removing %q\n", path)

		err = os.RemoveAll(path)
		if err != nil {
			return errors.Wrapf(err, "Couldn't remove installed package files for %q", pkg)
		}
	}

	pkgInstalledFiles, err := filepath.Glob(filepath.Join(ws.Pkg(), "*", pkgDir+".a"))
	if err != nil {
		return errors.Wrapf(err, "Something went wrong when globbing files for %q", pkg)
	}

	for _, path := range pkgInstalledFiles {
		fmt.Fprintf(logWriter, "  Removing %q\n", path)

		err = os.RemoveAll(path)
		if err != nil {
			return errors.Wrapf(err, "Couldn't remove installed package files for %q", pkg)
		}
	}
	return nil
}

func (ws *Workspace) ClearSrc() error {
	settings, err := ws.Settings()
	if err != nil {
		return err
	}

	for pkg := range settings.LocalInstalls {
		err := ws.Uninstall(pkg, os.Stdout)
		if err != nil {
			return err
		}
	}

	return errors.WithStack(os.RemoveAll(ws.Src()))

}

func (ws *Workspace) UnpersistLocalInstall(pkg string) error {
	settings, err := ws.Settings()
	if err != nil {
		return err
	}

	if install, ok := settings.LocalInstalls[pkg]; ok {
		fmt.Printf("Removing %q from persistent local installs\n", pkg)
		install.Persistent = false

		err = ws.SaveSettings(settings)
		if err != nil {
			return err
		}
	}

	return nil
}
