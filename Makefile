SHELL=/bin/bash
GO_FILES = $(shell find . -name "*.go" | grep -v "^./vendor/" |grep -v "_test.go$$" |  xargs)

BINDATA = cmd/bindata.go

ifndef VIRTUALGO
    $(error No virtualgo workspace is not active)
endif

.installed-deps:
	vg ensure
	touch .installed-deps

install: $(GO_FILES) $(BINDATA) .installed-deps
	go install
	@# install vg executable globally as well
	cp $(GOBIN)/vg $(_VIRTUALGO_OLDGOBIN)/vg

bindata: $(BINDATA) .installed-deps
$(BINDATA): data/*
	go-bindata -o cmd/bindata.go -pkg cmd data/*
