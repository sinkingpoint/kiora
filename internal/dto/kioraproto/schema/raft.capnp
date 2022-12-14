using Alert = import "./alert.capnp";
using Silence = import "./silence.capnp";
using Go = import "/go.capnp";
@0xefede801da2b560c;

$Go.package("kioraproto");
$Go.import("internal/dto/kioraproto");

# A message that is replicated around the Raft cluster
struct RaftLog {
    log :union {
        alerts @0 :Alert.PostAlertsRequest;
        silences @1 :Silence.PostSilencesRequest;
    }
}
