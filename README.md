# Kiora

This is a prototype replacement of Prometheus Alertmanager, without the pathological fear of complexity. Alertmanager itself strives to exclude business logic as much as possible, relying on operators to handle basic things like ratelimiting, and data validation. Kiora attempts to buck that idea. Its goal is to provide as much customization, observability, and control over the whole alert lifecycle as possible.

## Features

As of when I wrote this, these are some of the notable features implemented in Kiora:

 - Alert notifications
 - Alert resolved notifications
 - Alert silencing
 - Alert acknowledging
 - Clustering/HA, with Hashicorp Serf

Here's what I want to work on:

 - Silence/Ack data validation
 - Alert Statistics
 - Multi-Tenancy and Rate limiting
 - Alert Histories
 - A UI

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

## Configuration

All Kiora configurations are also valid [Graphviz Dot](https://graphviz.org/doc/info/lang.html) files, allowing you to define flows for alerts, silences, and any other model as it passes through the system. See the examples/ folder for more concrete examples.

### Starting Point

The simplest configuration is an empty graph:

```
digraph config {}
```

Which does no validation and sends alerts nowhere. But that's not particularly interesting.

### Notifiers

Notifiers (things that actually tell people about alerts), are configured as nodes in the graph. For example:

```
digraph config {
    console [type="stdout"];
    alerts -> console;
}
```

This defines a notifier "console", that writes alerts to stdout. We then define a "link" from the alerts pseudo-node, to the console notifier. All alerts start at the "alerts" node, and make their way through the graph, notifying any notifiers as they go.

### Filtering

Sometime, you want to conditionally send alerts to places. You can do that with "filters" on the edges of your graph. For example, we could only send alerts that have a `destination` label containing "console" to the console notifier:

```
digraph config {
    console [type="stdout"];
    alerts -> console [type="regex" field="destination" regex="console"];
}
```
