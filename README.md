# virtualgo [![Build Status](https://travis-ci.org/GetStream/vg.svg?branch=master)](https://travis-ci.org/GetStream/vg)

Virtualgo (or `vg` for short) is a tool which provides workspace based
development for Go. Its main feature set that makes it better than other
solutions is as follows:

1. Extreme ease of use
2. No interference with other go tools
3. Version pinning for imports
4. Version pinning for executables, such as linters (e.g. [`errcheck`](https://github.com/kisielk/errcheck)) and codegen tools (e.g. [`protoc-gen-go`](https://github.com/golang/protobuf))
5. Importing a dependency that's locally checked out outside of the workspace
   (also called multi project workflow)
6. Optional full isolation for imports, see the section on [import
   modes](#workspace-import-modes) for details.

Virtualgo doesn't do dependency resolution or version pinning itself, because
this is a hard problem that's already being solved by other tools. Its approach
is to build on top of these tools, such as
[`dep`](https://github.com/golang/dep), to provide the features features listed
above.
For people coming from Python `vg` is very similar to `virtualenv`, with `dep`
being respective to `pip`. The main difference is that `vg` is much easier to
use than `virtualenv`, because there's almost no mental overhead in using `vg`.

## Example usage

Below is an example showing some basic usage of `vg`. See further down and `vg help`
for more information and examples.

```bash
$ cd $GOPATH/src/github.com/GetStream/example
$ vg init  # initial creation of workspace

# Now all commands will be executed from within the example workspace
(example) $ go get github.com/pkg/errors # package only present in workspace
(example) $ vg ensure  # installs the dependencies of the example project using dep
(example) $ vg deactivate

$ cd ~
$ cd $GOPATH/src/github.com/GetStream/example
(example) $ # The workspace is now activated automatically after cd-ing to the project directory
```

## Advantages over existing solutions

The obvious question is: Why should you use `vg`? What advantages does it
bring over what you're using now? This obviously depends on what you're using
now:

### Advantages over `vendor` directory

1. You can pin versions of executable dependencies, such as linting and code
   generation tools.
2. No more issues with `go test ./...` running tests in the `vendor` directory
   when using `go` 1.8 and below.
3. You can easily use a dependency from your global `GOPATH` inside your
   workspace, without running into confusing import errors.
4. It has optional [full isolation](#workspace-import-modes). If enabled there's
   no accidental fallbacks to regular `GOPATH` causing confusion about what
   version of a package you're using.
5. When using full isolation, tools such as IDEs can spend much less time on
   indexing. This is simply because they don't have to index the packages
   outside the workspace.
6. You don't have problems when using plugins: https://github.com/akutz/gpd

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

Although not required, it is recommended to install
[`bindfs`](http://bindfs.org/) as well. This gives the best experience when
using [full isolation](#workspace-import-modes) and when using
`vg localInstall`. If you do this, DON'T remove things manually from
`~/.virtualgo`. Only use `vg destroy`/`vg uninstall`, otherwise you can very
well lose data.

```bash
# OSX
brew install bindfs
# Ubuntu
apt install bindfs
# Arch Linux
pacaur -S bindfs  # or yaourt or whatever tool you use for AUR
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

You can also edit your shell configuration file manually. Afterwards you still
have to `source` the file like explained above.

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
$ cd $GOPATH/src/github.com/GetStream/example
$ vg init
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

# It's also possible to only activate a workspace and not link it to the
# current directory. If the workspace doesn't exist it will also be
# created on the fly. Activating a new workspace automatically deactivates
# a previous one:
(example) $ vg activate example2
(example2) $ cobra
bash: cobra: command not found

# To deactivate the workspace simply run:
(example2) $ vg deactivate
$ vg activate
(example) $

# When a workspace is active, a go compilation will try to import packages
# installed from the workspace first. In some cases you might want to use the
# version of a package that is installed in your global GOPATH though. For
# instance when you are fixing a bug in a dependency and want to test the fix.
# In these cases you can easily install a package from your global GOPATH
# into the workspace:
(example) $ vg localInstall github.com/GetStream/utils
# You can even install a package from a specific path:
(example) $ vg localInstall github.com/GetStream/utils ~/weird/path/utils

# You can also uninstall a package from your workspace again
(example) $ vg uninstall github.com/spf13/cobra
# NOTE: At the moment this only removes the sources and static libs in pkg/, not
# executables. So the cobra command is still available.

# See the following sections for integration with dependency management tools.
# And for a full overview of all commands just run:
(example) $ vg help
# For detailed help of a specific command run:
(example) $ vg help <command>

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

## Workspace import modes

A workspace can be set up in two different import modes, global fallback or full
isolation.
The import mode of a workspace determines how imports from code behave and it is
chosen when the workspace is created.

### Global fallback
In global fallback mode, packages are imported from the original `GOPATH` when
they are not found in the workspace.
This is the default import mode for newly created workspaces, as this interferes
the least with existing go tools.

### Full isolation

In full isolation mode, package imports will only search in the packages that
are installed inside the workspace.
This has some advantages:

1. Tools such as IDE's don't have to search the global GOPATH for imports, which
   can result in a significant speedup for operations such as indexing.
2. You always know the location of an imported package.
3. It's not possible to accidentally import of a package that is not managed by
   your vendoring tool of choice.

However, there's also some downsides to full isolation of a workspace. These are
all caused by the fact that the project you're actually working on is not inside
your `GOPATH` anymore. So normally go would not be able to find any imports
to it. This is partially worked around by locally installing the project into
your workspace, but it does not fix all issues.

In the sections below the remaining issues are described and you can decide for
yourself if the above advantages are worth the disadvantages. If you want to try
out full isolation you can create a new workspace using the `--full-isolation`
flag:

```bash
$ vg init --full-isolation
# To change an existing workspace, you have to destroy and recreate it
$ vg destroy example
$ vg activate example --full-isolation
```

This will cause the workspace to use full isolation import mode each time it is
activated in the future. So there's no need to specify the
`--full-isolation` flag on each activation afterwards.

#### With `bindfs` installed

If you have [`bindfs`](http://bindfs.org/) installed the issues you will run
into are only a slight inconvenience, for which easy workarounds exist. However,
it is important that you know about them, because they will probably cause
confusion otherwise. If you run into any other issues than the ones mentioned
here, [please report them](https://github.com/GetStream/vg/issues/new).

##### Relative packages in commands

The first set of issues happen when using relative reference to packages in
commands. For instance `go list ./...` will return weirdly formatted paths, such
as `_/home/stream/go/src/github.com/GetStream/vg`. Also, running
`go test ./...`, might cause an `init` function to be executed twice. This can
all easily be worked around by using absolute package paths for these commands.
So for the `vg` repo you would use the following alternatives:

```bash
# go list ./...
go list github.com/GetStream/vg/...
# go test ./...
go test github.com/GetStream/vg/...`
```

##### `dep` commands

Another issue that pops up is that `dep` doesn't allow it's commands to be
executed outside of the `GOPATH`. This is not a problem for `dep ensure`, since
you usually use `vg ensure`, which handles this automatically. However, this is
an issue for other commands, such as `dep status` and `dep init`. Luckily
there's an easy workaround for this as well. You can simply use `vg globalExec`,
to execute commands from your regular `GOPATH`, which fixes the issue:

```bash
vg globalExec -- vg init
vg globalExec -- vg status
```

#### Without `bindfs` installed

If `bindfs` is not installed, symbolic links will be used to do the local
install.
This has the same issues as described for `bindfs`, but there's also some extra
ones that cannot be worked around as easily.
The reason for this is that go tooling does not like symbolic links in `GOPATH`
([golang/go#15507](https://github.com/golang/go/issues/15507), [golang/go#17451](https://github.com/golang/go/issues/17451)).

Compiling will still work, but `go list github.com/...` will not list your
package. Other than that there are also issues when using `delve`
([#11](https://github.com/GetStream/vg/issues/11)). Because of these issues it
is NOT RECOMMENDED to use virtualgo in full isolation mode without `bindfs`
installed.

## Using a virtualgo workspace with an IDE (e.g. Gogland)

Because virtualgo is just a usability wrapper around changing your `GOPATH` for
a specific project it is usually quite easy to use it in combination with an
IDE. Just check out your `GOPATH` after activating a workspace and configure the
IDE accordingly. Usually if you show your `GOPATH` you will see two paths
separated by a colon:

```bash
$ echo $GOPATH
/home/stream/.virtualgo/myworkspace:/home/stream/go
```

If you can set this full string directly that is fine. For
[Gogland](https://www.jetbrains.com/go/) you have to add the first one first and
then the second one.

When using a workspace in full isolation mode it's even easier to set up as
there's only one `GOPATH` set.

```bash
$ echo $GOPATH
/home/stream/.virtualgo/myworkspace
```

## License

MIT

## Careers @ Stream

Would you like to work on cool projects like this? We are currently hiring for
talented Gophers in Amsterdam and Boulder, get in touch with us if you are
interested! tommaso@getstream.io
