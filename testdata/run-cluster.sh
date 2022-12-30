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
./artifacts/kiora -c ./testdata/kiora.dot --raft.bootstrap --raft.data-dir artifacts/kiora-raft-data-1 --raft.local-id 1 --web.listen-url localhost:4278 --raft.listen-url localhost:4281 > artifacts/logs/kiora-1.log 2>&1 &
./artifacts/kiora -c ./testdata/kiora.dot --raft.data-dir artifacts/kiora-raft-data-2 --raft.local-id 2 --web.listen-url localhost:4279 --raft.listen-url localhost:4282 > artifacts/logs/kiora-2.log 2>&1 &
./artifacts/kiora -c ./testdata/kiora.dot --raft.data-dir artifacts/kiora-raft-data-3 --raft.local-id 3 --web.listen-url localhost:4280 --raft.listen-url localhost:4283 > artifacts/logs/kiora-3.log 2>&1 &

echo 'Establishing the Cluster'
sleep 3
curl -XPOST -d '{"id":"2","address":"localhost:4282"}' localhost:4278/admin/raft/add_member
curl -XPOST -d '{"id":"3","address":"localhost:4283"}' localhost:4278/admin/raft/add_member

curl localhost:4278/admin/raft/status

read -r -d '' _
