#!/bin/bash
eval "$(vg eval --shell bash)"

rm -rf coverages

set -uex -o pipefail
go test -covermode=count -coverpkg="$(go list ./... | paste -sd ',' -)" -c github.com/GetStream/vg -o testbins/testvg

go build -i -o testbins/vg github.com/GetStream/vg/internal/testwrapper/vg


set +u
vg deactivate || true

export PATH=$PWD/testbins:$PATH

echo PATH="$PATH"

bash -c 'which vg'

! bash -c 'vg activate'
! bash -c 'vg deactivate'
! bash -c 'vg cdpackages'

set +x
eval "$(vg eval --shell bash)"
set -x

vg activate testWS
vg deactivate testWS

vg destroy testWS

vg activate testWS
vg destroy

vg version
