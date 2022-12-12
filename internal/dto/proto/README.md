# Proto

Cap'n Proto / Protobuf is hard, so this folder structure attempts to keep it somewhat contained. Here, we define structures that are used for both Prometheus (or other alerting clients) -> Kiora, and Kiora -> Kiora (via raft) communication. As such, it doesn't contain any _operational state_, e.g. Alerts don't have whether or not they are silenced in them - that's something that Kiora works out for itself.

## Generation

All of the code in this directory is automagically generated from the definitions in the schema/ directory. To recompile them, you'll need [capnproto](https://capnproto.org/), and [go-capnpc](https://github.com/capnproto/go-capnproto2) installed. With that, you can run the `make generate` command in the root of this repository which will dump out the compiled structs into this dir. Because these are generated, you should not attempt to modify anything in this directory.
