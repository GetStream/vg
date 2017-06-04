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

### Automatic shell configuration

You can run the following command to configure all supported shells
automatically:

```sh
vg setup
```

After this you have to reload (`source`) your config:

```sh
source ~/.bashrc                   # for bash
source ~/.zshrc                    # for zsh
source ~/.config/fish/config.fish  # for fish
```

### Manual shell configuration

You can also edit your shell config file manually. Afterwards you still have to
`source` the file like explained above.

For bash put this in your `.bashrc`:

```bash
eval "$(vg eval --shell bash)"
```

Or for zsh, put his in your `.zshrc`:

```zsh
eval "$(vg eval --shell zsh)"
```

Or for fish, put this in your `config.fish`:

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


## Using a virtualgo workspace with an IDE (e.g. Gogland)

Because virtualgo is just a usability wrapper around changing your `GOPATH` for
a specific project it is usually quite easy to use it in combination with an
IDE. For Gogland you can set multiple `GOPATH`s in the preferences window on a
per project basis. To find out which `GOPATH`s, activate your desired workspace
and run:

```bash
$ echo $GOPATH
/home/stream/.virtualgo/myworkspace:/home/stream/go
```

As you can see there's two path separated by a semicolon. If you can set this
string directly that is fine. For Gogland you have to add the first one (with
`.virtualgo` in it) first and then the second one.


## Possible future additions

- [x] `vg globalExec <command>`, run command outside current active workspace
- [x] `vg uninistall <package>`, uninstall a package from the workspace
  (pkg including cache)
- [x] `vg upgrade` update vg to the latest version and re-eval it


## License

MIT
