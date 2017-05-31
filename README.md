# virtualgo
Virtualgo (or `vg` for short) is a virtualenv like solution to the go package isolation
problem with `dep` integration. If you have had any of the following problems
this package will solve them for you:

- You want to depend on a specific version of an executable dependency. You find
  out the `vendor` directory works only for libraries.
- You work on two projects with a vendor directory and want to your local
  version of one of those into the other. Some of the packages in these vendor
  directories overlap and suddenly get a lot of import errors.
- You want to `vendor` plug-ins and run into the issues: https://github.com/akutz/gpd

## Installation

First install the package:

```bash
go get -u github.com/getstream/vg
```

For POSIX compatible shells put the following in your .bashrc/.zshrc/.whateverrc:

```bash
eval "$(vg eval)"
```

Or for fish, put this in your config.fish:

```fish
vg eval --shell fish | source
```

## Usage

The following commands are the main commands to use `vg`:

```bash
# Activate (and create) a workspace
vg activate myProject

# All go commands are now executed from within your workspace
go get github.com/pkg/errors

# Bind the currently active workspace to the current directory
vg connect
# Everytime you cd into this directory the workspace will be activated
# automatically

# Deactivate the current workspace
vg deactivate
# Activating a new workspace automatically deactivates the previous one as well

# Without a workspace name it will activate a workspace named after the current
# directory
cd $GOPATH/src/github.com/pkg/errors
vg activate # activates a workspace called errors

```



### `dep` integration

`vg` integrates well with `dep` (https://github.com/golang/dep):

```bash
# Install the dependencies from Gopkg.lock into your workspace instead of the
# vendor directory
vg ensure

# Pass options to `dep ensure`
vg ensure -- -update
```

It also extends dep with a way to install executable dependencies. The `vg` repo
itself uses it to install the go-bindata command. It does this by having the
following in `Gopkg.toml`

```toml
required = [
    'github.com/jteeuwen/go-bindata/go-bindata'
]

[metadata]
install_required = true
```

Running `vg ensure` after adding this will install the `go-bindata` command in
the `GOBIN` of the current workspace. If you don't want to install all packages
in the required list (or install more packages) you can also provide a custom
list to install:

```toml
[metadata]
install = [
    'github.com/jteeuwen/go-bindata/go-bindata'
]
```


## How it works

All workspaces are fully isolated from each other. However, package and
executable resolution will fall back to your regular `GOPATH` and `GOBIN`.
This is done by using a very simple trick that is not well known: `GOPATH` can
contain multiple paths. So, `vg activate` simply prepends the workspace path to
the `GOPATH`. Furthermore it changes `GOBIN` and `PATH` to use `bin` directory
in the workspace. Running `vg deactivate` undoes these changes again.


## License

MIT
