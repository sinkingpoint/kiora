using Go = import "/go.capnp";
@0xd9e4592b8a563262;

$Go.package("kioraproto");
$Go.import("internal/dto/kioraproto");

struct Silences {
    silences @0 :List(Silence);
}

struct Silence {

}