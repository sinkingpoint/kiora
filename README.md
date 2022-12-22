# Kiora

This is a prototype replacement of Prometheus Alertmanager, without the pathological fear of complexity. Alertmanager itself strives to exclude business logic as much as possible, relying on operators to handle basic things like ratelimiting, and data validation. Kiora attempts to buck that idea.

Instead of Alertmanagers gossip protocol, which is notoriously opaque, Kiora uses a raft cluster to maintain alerts and silences. It's goal is to provide as much customization, observability, and control over the whole alert lifecycle as possible.

## Usage

```
Usage: kiora

Flags:
  -h, --help                                Show context-sensitive help.
      --web.listen-url="localhost:4278"     the address to listen on
  -c, --config.file="./kiora.dot"           the config file to load config from
      --raft.data-dir="./kiora/data"        the directory to put database state in
      --raft.bootstrap                      If set, bootstrap a new raft cluster
      --raft.local-id=""                    the name of this node in the raft cluster
      --raft.listen-url="localhost:4279"    the address for the raft node to listen on
```
