# Kiora

This is a prototype replacement of Prometheus Alertmanager, without the pathological fear of complexity. Alertmanager itself strives to exclude business logic as much as possible, relying on operators to handle basic things like ratelimiting, and data validation. Kiora attempts to buck that idea. It's goal is to provide as much customization, observability, and control over the whole alert lifecycle as possible.

## Usage

```
Usage: kiora

An experimental Alertmanager

Flags:
  -h, --help                                                   Show context-sensitive help.
      --web.listen-url="localhost:4278"                        the address to listen on
  -c, --config.file="./kiora.dot"                              the config file to load config from
      --cluster.node-name=STRING                               the name to join the cluster with
      --cluster.listen-url="localhost:4279"                    the address to run cluster activities on
      --cluster.bootstrap-peers=CLUSTER.BOOTSTRAP-PEERS,...    the peers to bootstrap with
```
