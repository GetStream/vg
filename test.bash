#!/bin/bash
eval "$(vg eval --shell bash)"
set -uex -o pipefail

rm -rf coverages

go test -covermode=count -coverpkg="$(go list ./... | paste -sd ',' -)" -c github.com/GetStream/vg -o testbins/testvg

go build -i -o testbins/vg github.com/GetStream/vg/internal/testwrapper/vg


vg deactivate || true

go get github.com/pkg/errors

export PATH=$PWD/testbins:$PATH

echo PATH="$PATH"

bash -c 'which vg'

! bash -c 'vg activate'
! bash -c 'vg deactivate'
! bash -c 'vg cdpackages'
! bash -c 'vg init'

set +xu
eval "$(vg eval --shell bash)"
set -xu

cd testbins
echo "PATH"
which vg
which testvg
vg init
which vg
which testvg
vg link
which vg
which testvg
vg unlink
which vg
which testvg
vg destroy

vg version
vg status

vg activate testWS
vg status
vg deactivate testWS
vg destroy testWS

cd ../
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

