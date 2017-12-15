package workspace

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

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

	hasBindfs, err := utils.CommandExists("bindfs")
	if err != nil {
		return err
	}

	if hasBindfs {
		err = ws.installLocalPackageWithBindfs(pkg, localPath, target)
	} else {
		err = ws.installLocalPackageWithSymlink(pkg, localPath, target)
	}

	if err != nil {
		return err
	}

	settings.LocalInstalls[pkg].Successful = true
	return ws.SaveSettings(settings)

}

func (ws *Workspace) installLocalPackageWithBindfs(pkg, src, target string) error {
	_, _ = fmt.Fprintf(os.Stderr, "Installing local sources at %q in workspace as %q using bindfs\n", src, pkg)

	settings, err := ws.Settings()
	if err != nil {
		return err
	}

	settings.LocalInstalls[pkg].Bindfs = true

	err = ws.SaveSettings(settings)
	if err != nil {
		return err
	}

	err = os.MkdirAll(target, 0755)
	if err != nil {
		return errors.WithStack(err)
	}

	cmd := exec.Command("bindfs", "--no-allow-other", src, target)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stderr
	return errors.WithStack(cmd.Run())
}

func (ws *Workspace) installLocalPackageWithSymlink(pkg, src, target string) error {
	_, _ = fmt.Fprintf(os.Stderr, "Installing local sources at %q in workspace as %q using symbolic links\n", src, pkg)

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

func (ws *Workspace) InstallSavedLocalPackages() error {
	settings, err := ws.Settings()
	if err != nil {
		return err
	}

	for pkg, install := range settings.LocalInstalls {
		if install.Persistent {
			err = ws.InstallLocalPackagePersistently(pkg, install.Path)
		} else {
			err = ws.InstallLocalPackage(pkg, install.Path)
		}

		if err != nil {
			return err
		}
	}
	return nil
}

func (ws *Workspace) Uninstall(pkg string, logWriter io.Writer) error {
	return ws.uninstall(pkg, logWriter, "")
}

func (ws *Workspace) uninstall(pkg string, logWriter io.Writer, indent string) error {
	pkgDir := utils.PkgToDir(pkg)
	pkgSrc := filepath.Join(ws.Src(), pkgDir)
	_, err := os.Stat(pkgSrc)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return errors.WithStack(err)
	}

	fmt.Fprintf(logWriter, indent+"Uninstalling %q from workspace\n", pkg)
	// Check if locally installed
	settings, err := ws.Settings()
	if err != nil {
		return err
	}
	indent += "  "

	// Uninstall all locally installed subpackages
	for localPkg := range settings.LocalInstalls {
		if strings.HasPrefix(pkg, localPkg) && localPkg != pkg {
			return errors.Errorf("Cannot uninstall %q, because it's a subpackage of the locally installed %q", pkg, localPkg)
		}
	}

	install, localInstalled := settings.LocalInstalls[pkg]

	if localInstalled && install.Bindfs {
		fmt.Fprintf(logWriter, indent+"Unmounting bindfs mount at %q\n", pkgSrc)
		stderrBuff := &bytes.Buffer{}
		outputBuff := &bytes.Buffer{}

		var cmd *exec.Cmd
		var notMounted func(output, pkgSrc string) bool

		hasFusermount, err := utils.CommandExists("fusermount")
		if err != nil {
			return err
		}

		if hasFusermount {
			// Use fusermount if that exists
			cmd = exec.Command("fusermount", "-u", pkgSrc)
			notMounted = fusermountNotMounted
		} else {
			// Otherwise fallback to umount
			cmd = exec.Command("umount", pkgSrc)
			notMounted = umountNotMounted
		}

		cmd.Stderr = io.MultiWriter(stderrBuff, outputBuff)
		cmd.Stdout = outputBuff

		err = cmd.Run()

		if err != nil {
			if !notMounted(stderrBuff.String(), pkgSrc) {
				// We don't care if the write to stderr failed
				_, _ = io.Copy(os.Stderr, outputBuff)

				return errors.WithStack(err)
			}
		}

	}

	// Uninstall all locally installed subpackages
	for localPkg := range settings.LocalInstalls {
		if strings.HasPrefix(localPkg, pkg) && localPkg != pkg {
			err := ws.uninstall(localPkg, logWriter, indent)
			if err != nil {
				return err
			}
		}
	}

	fmt.Fprintf(logWriter, indent+"Removing sources at %q\n", pkgSrc)
	err = os.RemoveAll(pkgSrc)
	if err != nil {
		return errors.Wrapf(err, "Couldn't remove package src %q", ws.Name())
	}

	if localInstalled && !install.Persistent {
		fmt.Fprintf(logWriter, indent+"Removing %q from locally installed packages\n", pkg)
		delete(settings.LocalInstalls, pkg)
		err = ws.SaveSettings(settings)
		if err != nil {
			return err
		}
	}

	pkgInstalledDirs, err := filepath.Glob(filepath.Join(ws.Pkg(), "*", pkgDir))
	if err != nil {
		return errors.Wrapf(err, "Something went wrong when globbing files for %q", pkg)
	}

	for _, path := range pkgInstalledDirs {
		fmt.Fprintf(logWriter, indent+"Removing %q\n", path)

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
		fmt.Fprintf(logWriter, indent+"Removing %q\n", path)

		err = os.RemoveAll(path)
		if err != nil {
			return errors.Wrapf(err, "Couldn't remove installed package files for %q", pkg)
		}
	}
	return nil
}

func umountNotMounted(output, pkgSrc string) bool {
	return strings.HasPrefix(output, fmt.Sprintf("umount: %s: not currently mounted", pkgSrc)) || strings.HasPrefix(output, fmt.Sprintf("umount: %s: not mounted", pkgSrc))
}

func fusermountNotMounted(output, pkgSrc string) bool {
	return strings.HasPrefix(output, fmt.Sprintf("fusermount: entry for %s not found", pkgSrc))
}

func (ws *Workspace) ClearSrc() error {
	settings, err := ws.Settings()
	if err != nil {
		return err
	}

	err = os.RemoveAll(ws.ensureMarker())
	if err != nil {
		return errors.WithStack(err)
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
