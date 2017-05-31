
GO_FILES = $(shell find . -name "*.go" | grep -v "^./vendor/" |grep -v "_test.go$$" |  xargs)

BINDATA = cmd/bindata.go

install: $(GO_FILES) $(BINDATA)
	go install

bindata: $(BINDATA)
$(BINDATA): data/*
	go-bindata -o cmd/bindata.go -pkg cmd data/*
