# virtualgo

Virtualgo (or `vg` for short) is a tool which provides workspace based
development for Go. The goals of the project are as follows:

1. Must be extremely easy to use
2. Shouldn't interfere with other go tools
3. Must support full isolation, for both imports and installed executables

It doesn't do dependency resolution or version pinning itself, but it
integrates well with `dep` and other dependency management tools.
For people coming from Python it's very similar to `virtualenv`, except that
it's much easier to use.

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
   `GOPATH` causing confusion?
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
[`required`](https://github.com/golang/dep/blob/master/FAQ.md#when-should-i-use-constraint-override-required-or-ignored-in-gopkgtoml)
list from `dep`.
However, if you don't want to install
all packages in the `required` list you can achieve that by putting the
following in `Gopkg.toml`:

```
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



## License

MIT

## Careers @ Stream

Would you like to work on cool projects like this? We are currently hiring for talented Gophers in Amsterdam and Boulder, get in touch with us if you are interested! tommaso@getstream.io

