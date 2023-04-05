# Kiora

This is a prototype replacement of Prometheus Alertmanager, without the pathological fear of complexity. Alertmanager itself strives to exclude business logic as much as possible, relying on operators to handle basic things like ratelimiting, and data validation. Kiora attempts to buck that idea. Its goal is to provide as much customization, observability, and control over the whole alert lifecycle as possible.

## Features

As of when I wrote this, these are some of the notable features implemented in Kiora:

 - Alert firing notifications
 - Alert resolved notifications
 - Alert silencing
 - Alert acknowledging
 - Clustering/HA, with Hashicorp Serf
 - Silence/Ack data validation
 - Alert Grouping
 - A basic UI

Here's what I want to work on:

 - Alert Statistics
 - Multi-Tenancy and Rate limiting
 - Alert Histories

## Usage

```
Usage: kiora

An experimental Alertmanager

Flags:
  -h, --help                                                   Show context-sensitive help.
      --tracing.service-name=STRING
      --tracing.exporter-type=STRING
      --tracing.destination-url=STRING
      --web.listen-url="localhost:4278"                        the address to listen on
  -c, --config.file="./kiora.dot"                              the config file to load config from
      --cluster.node-name=STRING                               the name to join the cluster with
      --cluster.listen-url="localhost:4279"                    the address to run cluster activities on
      --cluster.shard-labels=CLUSTER.SHARD-LABELS,...          the labels that determine which node in a cluster will send a given alert
      --cluster.bootstrap-peers=CLUSTER.BOOTSTRAP-PEERS,...    the peers to bootstrap with
      --storage.backend="boltdb"                               the storage backend to use
      --storage.path="./kiora.db"                              the path to store data in
```

## Prometheus Configuration

Kiora provides a compatibility shim with the Prometheus Alertmanager API. Simply configure your Kiora instance as another Alertmanager, with the "api/prom-compat" path prefix:

```
alerting:
  alertmanagers:
    - path_prefix: api/prom-compat
      static_configs:
        - targets:
          - 0.0.0.0:4278
```

## Configuration

All Kiora configurations are also valid [Graphviz Dot](https://graphviz.org/doc/info/lang.html) files, allowing you to define flows for alerts, silences, and any other model as it passes through the system. See the [examples](examples) folder for more concrete examples.
Alerts
## Pseudo-Nodes

Kiora works with the concept of "pseudo nodes". These are nodes in the graph that act as either sources or sinks of data. We currently have three pseudo-nodes:

1. `alerts` - alerts that come into the system flow out of the `alerts` pseudo-node, following the graph and notifying any notifiers that they hit.
2. `silences` - silences that come into the system are only accepted if they have a valid path _into_ the `silences` psuedo-node. See [this example](examples/silence_validation.dot) for a more concrete example of this.
3. `acks` - similar to silences, alert acknowledgements that come into the system are only accepted if the have a valid path into the `acks` pseudo-node.

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

### Routing

Sometime, you want to conditionally send alerts to places. You can do that with "filters" on the edges of your graph. For example, we could only send alerts that have a `destination` label containing "console" to the console notifier:

```
digraph config {
    console [type="stdout"];
    alerts -> console [type="regex" field="destination" regex="console"];
}
```

## Data Validation

In order to enforce business rules on silences / alert acknowledgements, you can provide filters on links into the relevant pseudo-nodes. For example, to enforce that all acknowledgements contain an email in the creator field:

```
digraph config {
    test_email -> acks [type="regex" field="creator" regex=".+@example.com"]; // Check for an @example.com email
}
```

Note how this flow works - acknowledgments start at the leaf nodes of the tree, and work their way through the filters. If there's a path into the `acks` node for which the acknowledgement passes all the filters, then the acknowledgement is accepted, otherwise it is rejected.
