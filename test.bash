#!/bin/bash
eval "$(vg eval --shell bash)"
set -uex -o pipefail

rm -rf coverages

go test -covermode=count -coverpkg=./... -c github.com/GetStream/vg -o testbins/testvg

go build -i -o testbins/vg github.com/GetStream/vg/internal/testwrapper/vg


vg deactivate || true

go get github.com/pkg/errors

export PATH=$PWD/testbins:$PATH

export COVERDIR=$PWD/coverages

echo PATH="$PATH"
vg setup

bash -c 'which vg'

! bash -c 'vg activate'
! bash -c 'vg deactivate'
! bash -c 'vg cdpackages'
! bash -c 'vg init'

vg setup
vg setup
vg version
vg status


set +xu
eval "$(vg eval --shell bash)"
set -xu

vg activate testWS
vg status
vg deactivate testWS
vg destroy testWS

vg activate testWS
vg ensure -- -v
vg uninstall github.com/pkg/errors
vg localInstall github.com/pkg/errors
vg uninstall github.com/pkg/errors
vg destroy

vg activate testWS
vg globalExec dep ensure -v
vg moveVendor
vg destroy


vg activate testWS --full-isolation
vg destroy

vg activate testWS --global-fallback
vg destroy

cd testbins
vg init
vg link
vg unlink
vg destroy
