#!/bin/bash
eval "$(vg eval --shell bash)"

rm -rf coverages

set -uex -o pipefail
go test -coverprofile=coverage.out -coverpkg=github.com/GetStream/vg,github.com/GetStream/vg/cmd -c main.go main_test.go -o testbins/testvg

go build -i -o testbins/vg github.com/GetStream/vg/internal/testwrapper/vg


set +u
vg deactivate || true

export PATH=testbins:$PATH

echo PATH="$PATH"

bash -c 'which vg'

! bash -c 'vg activate'
! bash -c 'vg deactivate'

set +x
eval "$(vg eval --shell bash)"
set -x

vg activate testWS
vg deactivate testWS

vg destroy testWS

vg activate testWS
vg destroy
