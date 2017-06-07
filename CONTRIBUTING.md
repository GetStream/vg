# How to contribute

The simplest way to develop virtualgo is by using virtualgo. The initial setup
is simple:

```
go get github.com/GetStream/vg
cd $GOPATH/src/github.com/GetStream/vg
vg init
```

Then just change something and run:

```
make
```


## Adding a new command

Virtualgo uses cobra under the hood. Adding a new command is done like so:

```
cobra add <commandName>
```

Then you can edit the behaviour for that command by editing
`cmd/<commandName>.go`.

## Changing shell specific behaviour

All shell specific behaviour is placed in the `data` directory:

- `data/sh` is used by zsh and bash
- `data/bash` is used by bash only
- `data/zsh` is used by zsh
- `data/fish` is used by fish
