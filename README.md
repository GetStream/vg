# virtualgo

Virtualgo (or `vg` for short) is a tool to add workspace based development to go
projects. Its goal is to improve on the currently ubiquitous `vendor` directory
based approach. It's extremely easy to use and optionally integrates with
[`dep`](https://github.com/golang/dep) to allow for version pinning of
dependencies.

Virtualgo solves many problems with the `vendor` directory by using an extra
`GOPATH` for each project instead. One way this improves upon `vendor` is by
allowing specific versions of an executable to be installed in a workspace, such
as linters and codegen tools. Another advantage of this approach is that all the
`go` commands can be used in the normal way, but they will affect the workspace
instead of the regular `GOPATH`. If you were already using the `vendor`
directory and you had any of the following problems this package will solve
them:

- When running `go test ./...` the tests in your `vendor` directory are executed.
- You want to depend on a specific version of an executable package, such as a
  linter or a codegen tool. You find out that the `vendor` directory only works
  for libraries.
- You work on two projects, A and B. Both of them contain a `vendor` directory.
  You want to use project A from your `GOPATH` when compiling B. To do this you
  remove A from the `vendor` directory of B, so it will fall back to import A
  from `GOPATH`. Suddenly you get a lot of weird import errors.
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
# The most used command is to create and activate a workspace named after the
# current direcory
vg init
# In the future each time you cd to this directory it will be activated
# automatically.

# All go commands in this shell are now executed from within your workspace. The
# following will install github.com/pkg/errors inside the workspace
go get github.com/pkg/errors
# (See below in the README on how to use the workspace from an IDE)

# It's also possible to only activate (and create) a workspace and not link it
# to the current directory.
vg activate myProject

# You can then link the currently active workspace to the current directory
vg link

# You can also uninstall a package from your workspace
vg uninstall github.com/pkg/errors
# If the removed package is installed in your normal GOPATH as well imports will
# now use that version instead. This can be usefull when patching one of
# your projects dependencies. This way you can use test the patches directly
# inside your workspace.

# Deactivate the current workspace
vg deactivate
# Activating a new workspace automatically deactivates the previous one as well

# For a full overview of all commands
vg help

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

### Integration with other dependency management tools (e.g glide)

Even though `dep` is the main tool that virtualgo integrates with. It's also possible
to use other dependency management tools instead, as long as they create a
`vendor` directory. Installing executable dependencies is not supported though.

To use `vg` with `glide` works like this:

```bash
# Install dependencies into vendor with glide
glide install

# Move these dependencies into the workspace
vg moveVendor
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

As you can see there's two path separated by a colon. If you can set this
string directly that is fine. For Gogland you have to add the first one (with
`.virtualgo` in it) first and then the second one.


## Comparison to similar tools

The main difference between virtualgo and other similar tools is that it's just
an easy wrapper around a feature that is already built into `go` itself, having
multiple `GOPATH`s. Because of this all `go` commands simply keep working as
they normally do.

- [`gb`](https://github.com/constabulary/gb) requires to use the gb command for
  everything.
- [`wgo`](https://github.com/skelterjohn/wgo) uses the
  `vendor` directory and thus has all the same issues mentioned at the start of the
  README (e.g, no version pinning of executables).


## License

MIT

## Careers @ Stream

Would you like to work on cool projects like this? We are currently hiring for talented Gophers in Amsterdam and Boulder, get in touch with us if you are interested! tommaso@getstream.io

