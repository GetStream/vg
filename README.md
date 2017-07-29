# virtualgo [![Build Status](https://travis-ci.org/GetStream/vg.svg?branch=master)](https://travis-ci.org/GetStream/vg)

Virtualgo (or `vg` for short) is a tool which provides workspace based
development for Go. Its main feature set that makes it better than other
solutions is as follows:

1. Extreme easy of use
2. No interference with other go tools
3. Version pinning for imports
4. Version pinning for executables, such as linters (e.g. [`errcheck`](https://github.com/kisielk/errcheck)) and codegen tools (e.g. [`protoc-gen-go`](https://github.com/golang/protobuf))
5. Importing a dependency that's locally checked out outside of the workspace
   (also called multi project workflow)
6. Optional full isolation, for both imports and installed executables

Virtualgo doesn't do dependency resolution or version pinning itself, because
this is a hard problem that's already being solved by other tools. Its approach
is to build on top of these tools, such as
[`dep`](https://github.com/golang/dep), to provide the features features listed
above.
For people coming from Python `vg` is very similar to `virtualenv`, with `dep`
being respective to `pip`. The main difference is that `vg` is much easier to
use than `virtualenv`.

## Example usage

Below is an example showing some basic usage of `vg`. See further down and `vg help`
for more information and examples.

```bash
$ cd $GOPATH/src/github.com/Getstream/example
$ vg init  # initial creation of workspace

# Now all commands will be executed from within the example workspace
(example) $ go get github.com/pkg/errors # package only present in workspace
(example) $ vg ensure  # installs the dependencies of the example project using dep
(example) $ vg deactivate

$ cd ~
$ cd $GOPATH/src/github.com/Getstream/example
(example) $ # The workspace is now activated automatically after cd-ing to the project directory
```

## Advantages over existing solutions

The obvious question is: Why should you use `vg`? What advantages does it
bring over what you're using now? This obviously depends on what you're using
now:

### Advantages over `vendor` directory

1. You can pin versions of executable dependencies, such as linting and code
   generation tools.
2. It has full isolation by default, so no accidental fallbacks to regular
   `GOPATH` causing confusion about what version you're using.
3. No more issues with `go test ./...` running tests in the vendor directory.
4. You can easily use a dependency from your global `GOPATH` inside your
   workspace, without running into confusing import errors.
5. You don't have problems when using plugins: https://github.com/akutz/gpd

### Advantages over manually managing multiple `GOPATH`s

1. Automatic activation of a `GOPATH` when you `cd` into a directory.
2. Integration with version management tools such as `dep` and `glide` allow for
   reproducible builds.
3. Useful commands to manage installed packages. For instance for uninstalling
   a package or installing a local package from another `GOPATH`.


## Installation

First install the package:

```bash
go get -u github.com/GetStream/vg
```

### Automatic shell configuration

You can run the following command to configure all supported shells
automatically:

```sh
vg setup
```

After this you have to reload (`source`) your shell configuration file:

```sh
source ~/.bashrc                   # for bash
source ~/.zshrc                    # for zsh
source ~/.config/fish/config.fish  # for fish
```

### Manual shell configuration

You can also edit your shell config file manually. Afterwards you still have to
`source` the file like explained above.

For bash put this in your `~/.bashrc` file:

```bash
command -v vg >/dev/null 2>&1 && eval "$(vg eval --shell bash)"
```

Or for zsh, put his in your `~/.zshrc` file:

```zsh
command -v vg >/dev/null 2>&1 && eval "$(vg eval --shell zsh)"
```

Or for fish, put this in your `~/.config/fish/config.fish` file:

```fish
command -v vg >/dev/null 2>&1; and vg eval --shell fish | source
```

## Usage

The following commands are the main commands to use `vg`:

```bash
# The first command to use is the one to create and activate a workspace named
# after the current direcory
$ cd $GOPATH/src/github.com/Getstream/example
(example) $
# This command also links the current directory to the created workspace. This
# way the next time you cd to this directory the workspace will be activated
# automatically.
# (See below in the README on how to use the workspace from an IDE)

# All go commands in this shell are now executed from within your workspace. The
# following will install the most recent version of the cobra command and
# library inside the workspace
(example) $ go get -u github.com/spf13/cobra/cobra
(example) $ cobra
Cobra is a CLI library for Go that empowers applications.
......

# It's also possible to only activate (and create) a workspace and not link it
# to the current directory. Activating a new workspace automatically deactivates
# a previous one:
(example) $ vg activate example2
(example2) $ cobra
bash: cobra: command not found

# To deactivate the workspace simply run:
(example2) $ vg deactivate
$ vg activate
(example) $

# When a workspace is active go builds cannot import packages from your
# normal GOPATH (you can still use executables though). This is good for
# isolation as you can not accidentally import something outside of the
# workspace. However, you can easily install a package from your global GOPATH
# into the workspace.
(example) $ vg localInstall github.com/GetStream/utils
# You can even install a package from a specific path
(example) $ vg localInstall github.com/GetStream/utils ~/weird/path/utils


# You can also uninstall a package from your workspace again
(example) $ vg uninstall github.com/spf13/cobra
# NOTE: At the moment this only removes the sources and static libs in pkg/, not
# executables. So the cobra command is still available.

# See the following sections for integration with dependency management tools.
# And for a full overview of all commands just run:
(example) $ vg help

```



### `dep` integration

`vg` integrates well with `dep` (https://github.com/golang/dep):

```bash
# Install the dependencies from Gopkg.lock into your workspace instead of the
# vendor directory
vg ensure

# Pass options to `dep ensure`
vg ensure -- -v -update
```

It also extends `dep` with a way to install executable dependencies. The `vg`
repo itself uses it to install the `go-bindata` and `cobra` command. It does
this by adding the following in `Gopkg.toml`:

```toml
required = [
    'github.com/jteeuwen/go-bindata/go-bindata',
    'github.com/spf13/cobra/cobra'
]
```
Running `vg ensure` after adding this will install the `go-bindata` and `cobra`
command in the `GOBIN` of the current workspace.

As you just saw `vg` reuses the
[`required`](https://github.com/golang/dep/blob/master/docs/Gopkg.toml.md#required)
list from `dep`.
However, if you don't want to install
all packages in the `required` list you can achieve that by putting the
following in `Gopkg.toml`:

```toml
[metadata]
install-required = false
```

You can also specify which packages to install without the `required` list:
```toml
[metadata]
install = [
    'github.com/jteeuwen/go-bindata/go-bindata',
    'github.com/golang/mock/...', # supports pkg/... syntax
]
```

### Integration with other dependency management tools (e.g glide)

Even though `dep` is the main tool that virtualgo integrates with. It's also possible
to use other dependency management tools instead, as long as they create a
`vendor` directory. Installing executable dependencies is not supported though
(PRs for this are welcome).

To use `vg` with `glide` works like this:

```bash
# Install dependencies into vendor with glide
glide install

# Move these dependencies into the workspace
vg moveVendor
```

## Workspaces with global `GOPATH` fallback

There's one downside to full isolation of a workspace. The project you're
actually working on is not inside your `GOPATH` by default, so normally go
would not be able to find any imports to it. This is worked around by locally
installing the project into your workspace. This works quite well, but the
current implementation of this uses symbolic links and these are not supported
incredibly well by the go tooling. Commands such as `go list` will not list your
package because of it.

Also you can get some weird behaviour if you run a command such as `go test
./...`. The go package names that are shown will contain the full path of the
package with an underscore in front. This is because `.` is not in your
`GOPATH`. Usually this is not an issue  and everything compiles fine, but there
are cases where an `init` function will be executed twice.

These problems can all be solved by using `vg` in "global fallback" mode. In
this mode it's possible to import packages from your global `GOPATH`.
This is not the recommended way to use `vg`, because you won't have full
isolation anymore. However, if you are running into these problems this should
be an easy fix. You can create workspace in "global fallback" mode by running:

```bash
$ vg init --global-fallback
# To change an existing workspace, you should destroy and recreate it
$ vg destroy example
$ vg init example --global-fallback
```

If you create a workspace this way, any imports you do first search in the
workspace. If a package cannot be found there it will try the original `GOPATH`
you had before activating the workspace.

## Using a virtualgo workspace with an IDE (e.g. Gogland)

Because virtualgo is just a usability wrapper around changing your `GOPATH` for
a specific project it is usually quite easy to use it in combination with an
IDE. Just check out your `GOPATH` after activating a workspace and configure the
IDE accordingly.

```bash
$ echo $gopath
/home/stream/.virtualgo/myworkspace
```

When using a workspace with global `GOPATH` fallback, it's only a little harder
to configure your `GOPATH`. If you show your `GOPATH` you will see two paths
separated by a colon:

```bash
$ echo $gopath
/home/stream/.virtualgo/myworkspace:/home/stream/go
```

If you can set this full string directly that is fine. For Gogland you have to
add the first one first and then the second one.


## License

MIT

## Careers @ Stream

Would you like to work on cool projects like this? We are currently hiring for talented Gophers in Amsterdam and Boulder, get in touch with us if you are interested! tommaso@getstream.io

