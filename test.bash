#!/bin/bash
eval "$(vg eval --shell bash)"

rm -rf coverages

set -uex -o pipefail
go test -covermode=count -coverpkg="$(go list ./... | paste -sd ',' -)" -c github.com/GetStream/vg -o testbins/testvg

go build -i -o testbins/vg github.com/GetStream/vg/internal/testwrapper/vg


set +u
vg deactivate || true

go get github.com/pkg/errors

export PATH=$PWD/testbins:$PATH

echo PATH="$PATH"

bash -c 'which vg'

! bash -c 'vg activate'
! bash -c 'vg deactivate'
! bash -c 'vg cdpackages'

set +x
eval "$(vg eval --shell bash)"
set -x

vg version
vg status

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
vg destroy
