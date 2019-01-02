SHELL=/bin/bash
GO_FILES = $(shell find . -name "*.go" | grep -v "^./vendor/" |grep -v "_test.go$$" |  xargs)

BIN ?= vg
BIN_PACKAGE = .

REPO=github.com/GetStream/vg

RELEASE_DIR = releases

CURRENT_VERSION_MAJOR = 0
CURRENT_VERSION_MINOR = 9
CURRENT_VERSION_BUG = 0

BINDATA = cmd/bindata.go

ifndef VIRTUALGO
    $(error No virtualgo workspace is active)
endif

.PHONY: install
.PHONY: publish publish-major publish-minor publish-bug update-master

LAST_ENSURE = $(VIRTUALGO_PATH)/last-ensure

DEPS = $(LAST_ENSURE) $(GO_FILES) $(BINDATA)

all: install
get-deps: $(LAST_ENSURE)

install: $(DEPS)
	go install $(REPO)
	@# install vg executable globally as well
	cp $(GOBIN)/vg $(_VIRTUALGO_OLDGOBIN)/vg

bindata: $(BINDATA) .installed-deps
$(BINDATA): data/*
	go-bindata -nometadata -o cmd/bindata.go -pkg cmd data/*

$(LAST_ENSURE): Gopkg.lock Gopkg.toml
	vg ensure -- -v

publish: $(BINDATA)
	@if [ "$(VERSION)" = "" ]; then echo You should define the version like so: make publish VERSION=x.y.z; exit 1; fi
	@git diff --exit-code --cached || { git status; echo You have changes that are staged but not committed ; false ; };
	@git diff --exit-code || { git status; echo You have changes that are not committed ; false ; };
	$(eval dots = $(subst ., ,$(VERSION)))
	$(eval new_major = $(word 1, $(dots)))
	$(eval new_minor = $(word 2, $(dots)))
	$(eval new_bug = $(word 3, $(dots)))
	sed -i.bak -e 's/\(\tVersion string = \).*/\1"$(VERSION)-dev"/g' cmd/version.go
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

publish-staging: $(DEPS)
	@if [ "$(SUFFIX)" = "" ]; then echo You should define the version like so: make publish SUFFIX=test-ratelimit; exit 1; fi
	@git diff --exit-code --cached || { git status; echo You have changes that are staged but not committed ; false ; };
	@git diff --exit-code || { git status; echo You have changes that are not committed ; false ; };
	$(eval VERSION := $(CURRENT_VERSION_MAJOR).$(CURRENT_VERSION_MINOR).$(CURRENT_VERSION_BUG)-$(SUFFIX))
	git tag -a -m "RELEASED BY: $$(git config user.name)" v$(VERSION)
	git push --follow-tags

test:
	go test $(REPO)/internal/...
	./test.bash

cover: $(DEPS)
	go test -cover -coverpkg=./... -covermode=count -coverprofile=coverage-unit-tests.out $(REPO)/internal/...
	./test.bash
	gocovmerge coverages/*.out coverage-unit-tests.out > coverage.out
	rm coverage-unit-tests.out
	rm -r coverages

release-all: $(DEPS)
	GOOS=linux GOARCH=amd64 make release
	GOOS=windows GOARCH=amd64 make release
	GOOS=darwin GOARCH=amd64 make release

release: $(DEPS)
	mkdir -p $(RELEASE_DIR)
	go build $(BUILD_FLAGS) -ldflags="-w -s -X github.com/GetStream/vg/cmd.Version=`git describe`" -o "$(RELEASE_DIR)/$(BIN)-`go env GOOS`-`go env GOARCH`" $(BIN_PACKAGE)

clean:
	rm $(BINDATA)
