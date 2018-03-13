// Copyright © 2017 Stream
//

package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/BurntSushi/toml"
	"github.com/GetStream/vg/internal/utils"
	"github.com/GetStream/vg/internal/workspace"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type depConfig struct {
	Required []string
	Metadata struct {
		InstallRequiredDeprecated *bool `toml:"install_required"`
		InstallRequired           bool  `toml:"install-required"`
		Install                   []string
	}
}

// ensureCmd represents the ensure command
var ensureCmd = &cobra.Command{
	Use:   "ensure [-- [arguments to dep ensure]]",
	Short: "A wrapper for dep that installs the dependencies in the virtualgo workspace instead of vendor",
	Long: `To simlpy install the dependencies in Gopkg.lock you can run:

	vg ensure

It's also possible to pass arguments to dep ensure, such as:

	vg ensure -- -update
	vg ensure -- github.com/pkg/errors

This command also adds an extra feature to Gopkg.toml. You can install certain
packages (such as binaries) in the virtualgo workspace. This uses the
metadata section in the Gopkg.toml

	[metadata]
	# install all packges in the root required list
	install_required = true
	# install these specific packages
	install = [
	    'github.com/golang/dep/cmd/dep',
	    'github.com/golang/mock/...',
	]


This command requires that dep is installed in $PATH. `,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error

		ws, err := workspace.Current()
		if err != nil {
			return err
		}

		err = os.RemoveAll("vendor")
		if err != nil {
			return errors.Wrap(err, "Couldn't remove the current vendor directory")
		}

		if false {
			// TODO: This is causing some errors, packages are not actually
			// installed. Not sure why, maybe bug in go dep.
			err = os.Rename(ws.Src(), "vendor")
			if err != nil {
				err = err.(*os.LinkError).Err
				if err != syscall.ENOENT {
					// If src doesn't exist it doesn't have to be moved
					return errors.Wrap(err, "Couldn't move the the sources of the active workspace to vendor")
				}
			}
		}

		gopath := os.Getenv("GOPATH")
		err = os.Setenv("GOPATH", os.Getenv("_VIRTUALGO_OLDGOPATH"))
		if err != nil {
			return errors.WithStack(err)
		}
		depCmd := exec.Command("dep", append([]string{"ensure"}, args...)...)
		depCmd.Stderr = os.Stderr
		depCmd.Stdout = os.Stdout

		argsString := ""
		if len(args) > 0 {
			argsString = " " + strings.Join(args, " ")
		}

		fmt.Printf("Running %q\n", "dep ensure"+argsString)
		exitSignals := make(chan os.Signal, 1)
		signal.Notify(exitSignals, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
		go func() {
			sig := <-exitSignals
			p := depCmd.Process
			if p != nil {
				p.Signal(sig)
				p.Wait()
			}
			os.Exit(1)
		}()

		err = depCmd.Run()
		if err != nil {

			// Try to revert move after insuccessful dep
			// TODO: Uncomment when fixing above todo
			// _ = os.Rename("vendor", ws.Src())
			return errors.Wrap(err, "dep failed to run")
		}

		err = os.Setenv("GOPATH", gopath)
		if err != nil {
			return errors.WithStack(err)
		}

		exists, err := utils.VendorExists()
		if err != nil {
			return err
		}

		if !exists {
			fmt.Fprintln(os.Stderr, "ERROR: The vendor directory was not created by the dep ensure command.")
			fmt.Fprintln(os.Stderr, "ERROR: This means the workspace content won't be changed")
			os.Exit(1)
		}

		err = ws.ClearSrc()
		if err != nil {
			return errors.Wrap(err, "Couldn't clear the src path of the active workspace")
		}

		err = os.Rename("vendor", ws.Src())
		if err != nil {
			return errors.Wrap(err, "Couldn't move the vendor directory to the active workspace")
		}

		gopkgData, err := ioutil.ReadFile("Gopkg.toml")
		if err != nil {
			return errors.Wrap(err, "Couldn't read Gopkg.toml")
		}
		config := depConfig{}
		config.Metadata.InstallRequired = true // Good default

		_, err = toml.Decode(string(gopkgData), &config)
		if err != nil {
			return errors.Wrap(err, "Couldn't unmarshal Gopkg.toml")
		}

		err = ws.InstallSavedLocalPackages()
		if err != nil {
			return err
		}

		installRequired := config.Metadata.InstallRequired

		if config.Metadata.InstallRequiredDeprecated != nil {
			installRequired = *config.Metadata.InstallRequiredDeprecated
		}

		if installRequired {
			err := installPackages(ws.Src(), config.Required)
			if err != nil {
				return err
			}
		}

		err = installPackages(ws.Src(), config.Metadata.Install)
		if err != nil {
			return err
		}

		return ws.UpdateEnsureMarker()
	},
}

func installPackages(srcPath string, packages []string) error {
	for _, pkg := range packages {
		var recursive bool
		var installCmd *exec.Cmd
		fmt.Printf("Installing %q\n", pkg)
		pkgComponents := strings.Split(pkg, "/")
		if pkgComponents[len(pkgComponents)-1] == "..." {
			recursive = true
			pkgComponents = pkgComponents[:len(pkgComponents)-2]
		}

		pkgPath := filepath.Join(append([]string{srcPath}, pkgComponents...)...)
		if !recursive {
			installCmd = exec.Command("go", "install", ".")
		} else {
			installCmd = exec.Command("go", "install", "./...")
		}

		installCmd.Dir = pkgPath
		installCmd.Stderr = os.Stderr
		installCmd.Stdout = os.Stdout
		err := installCmd.Run()
		if err != nil {
			return errors.Wrapf(err, "Installation of %s failed", pkg)
		}

	}

	return nil
}

func init() {
	RootCmd.AddCommand(ensureCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ensureCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ensureCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
