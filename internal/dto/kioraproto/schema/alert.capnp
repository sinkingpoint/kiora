using Go = import "/go.capnp";
@0xbca48e54c9238692;

$Go.package("kioraproto");
$Go.import("internal/dto/kioraproto");

# An alert is the raw data in Kiora - something went wrong,
# and we might need to tell people about it.
struct Alert {
    labels @0 :Map(Text, Text);
    annotations @1 :Map(Text, Text);
    status @2 :AlertStatus;
    startTime @3 :Int64;
    endTime @4 :Int64;
}

# AlertStatus is the status of the alert coming in.
enum AlertStatus {
    # Firing alerts are active, and we should tell someone (maybe).
    firing @0;

    # Resolved alerts are no longer active, and Kiora should purge them.
    resolved @1;
}

struct Map(Key, Value) {
  entries @0 :List(Entry);
  struct Entry {
    key @0 :Key;
    value @1 :Value;
  }
}
