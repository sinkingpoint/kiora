#!/bin/bash

set -euo pipefail

make build

trap on_exit INT EXIT

function on_exit() {
    rm -rf artifacts/kiora-raft-data-*
    kill -- -$$
}

mkdir -p artifacts/logs

echo 'Starting Kiora'
./artifacts/kiora -c ./testdata/kiora.dot --web.listen-url localhost:4278 --cluster.listen-url localhost:4281 > artifacts/logs/kiora-1.log 2>&1 &
./artifacts/kiora -c ./testdata/kiora.dot --web.listen-url localhost:4279 --cluster.listen-url localhost:4282 --cluster.bootstrap-peers localhost:4281 > artifacts/logs/kiora-2.log 2>&1 &
./artifacts/kiora -c ./testdata/kiora.dot --web.listen-url localhost:4280 --cluster.listen-url localhost:4283 --cluster.bootstrap-peers localhost:4281 > artifacts/logs/kiora-3.log 2>&1 &

read -r -d '' _
