SHELL=/bin/bash
GO_FILES = $(shell find . -name "*.go" | grep -v "^./vendor/" |grep -v "_test.go$$" |  xargs)

CURRENT_VERSION_MAJOR = 0
CURRENT_VERSION_MINOR = 5
CURRENT_VERSION_BUG = 0

BINDATA = cmd/bindata.go

ifndef VIRTUALGO
    $(error No virtualgo workspace is not active)
endif

.PHONY: install
.PHONY: publish publish-major publish-minor publish-bug update-master

install: $(GO_FILES) $(BINDATA) .installed-deps
	go install
	@# install vg executable globally as well
	cp $(GOBIN)/vg $(_VIRTUALGO_OLDGOBIN)/vg

bindata: $(BINDATA) .installed-deps
$(BINDATA): data/*
	go-bindata -o cmd/bindata.go -pkg cmd data/*

.installed-deps: Gopkg.lock Gopkg.toml
	vg ensure -- -v
	touch .installed-deps

publish: $(BINDATA)
	@if [ "$(VERSION)" = "" ]; then echo You should define the version like so: make publish VERSION=x.y.z; exit 1; fi
	@git diff --exit-code --cached || { git status; echo You have changes that are staged but not committed ; false ; };
	@git diff --exit-code || { git status; echo You have changes that are not committed ; false ; };
	$(eval dots = $(subst ., ,$(VERSION)))
	$(eval new_major = $(word 1, $(dots)))
	$(eval new_minor = $(word 2, $(dots)))
	$(eval new_bug = $(word 3, $(dots)))
	sed -i.bak -e 's/\(\tVersion string = \).*/\1"$(VERSION)"/g' cmd/version.go
	sed -i.bak -e 's/^\(CURRENT_VERSION_MAJOR = \).*/\1$(new_major)/g' Makefile
	sed -i.bak -e 's/^\(CURRENT_VERSION_MINOR = \).*/\1$(new_minor)/g' Makefile
	sed -i.bak -e 's/^\(CURRENT_VERSION_BUG = \).*/\1$(new_bug)/g' Makefile
	rm Makefile.bak cmd/version.go.bak

	git commit -am 'Bump version to v$(VERSION)'
	git tag -m '' v$(VERSION)
	git push --follow-tags

update-master:
	git checkout master
	git pull

publish-major: update-master
	make publish VERSION=$$(($(CURRENT_VERSION_MAJOR) + 1)).0.0
publish-minor: update-master
	make publish VERSION=$(CURRENT_VERSION_MAJOR).$$(($(CURRENT_VERSION_MINOR) + 1)).0
publish-bug: update-master
	make publish VERSION=$(CURRENT_VERSION_MAJOR).$(CURRENT_VERSION_MINOR).$$(($(CURRENT_VERSION_BUG) + 1))
