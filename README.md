# virtualgo
Virtualgo (or `vg` for short) is a [`dep`](https://github.com/golang/dep)
compatible solution to problems caused be managing dependencies in the `vendor`
directory. If you have had any of the following problems this package will solve
them for you:

- When running `go test ./...` the tests in your vendor directory are executed.
- You want to depend on a specific version of an executable package. You find
  out that the `vendor` directory only works for libraries.
- You work on two projects, A and B. Both of them contain a vendor directory.
  You want to use project A from your `GOPATH` when compiling B. To do this you
  remove A from the vendor directory of B, so it will fallback to import A from
  `GOPATH`. Suddenly you get a lot of weird import errors.
- You want to `vendor` plug-ins and run into issues: https://github.com/akutz/gpd

## Installation

First install the package:

```bash
go get -u github.com/getstream/vg
```


For bash put this in your .bashrc:

```bash
eval "$(vg eval)"
```


Or for zsh, put his in your .zshrc:

```zsh
eval "$(vg eval --shell zsh)"
```

Or for fish, put this in your config.fish:

```fish
vg eval --shell fish | source
```

## Usage

The following commands are the main commands to use `vg`:

```bash
# Create and activate a workspace named after the current direcory. Each time
# you cd to this directory it will be activated automatically.
vg init

# All go commands are now executed from within your workspace. The followinging
# will install github.com/pkg/errors inside the workspace
go get github.com/pkg/errors

# It's also possible to only activate Activate (and create) a workspace and not
# link it to the current directory.
vg activate myProject

# You can then link the currently active workspace to the current directory
vg link

# Deactivate the current workspace
vg deactivate
# Activating a new workspace automatically deactivates the previous one as well
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
