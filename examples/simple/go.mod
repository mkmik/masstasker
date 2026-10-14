module mkm.pub/masstasker/examples/simple

go 1.25.0

replace mkm.pub/masstasker => ../..

require (
	github.com/prometheus/client_golang v1.19.1
	google.golang.org/grpc v1.83.2
	mkm.pub/masstasker v0.1.13
)

require (
	github.com/alecthomas/kong v1.9.0
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/jonboulle/clockwork v0.5.0 // indirect
	github.com/prometheus/client_model v0.5.0 // indirect
	github.com/prometheus/common v0.48.0 // indirect
	github.com/prometheus/procfs v0.12.0 // indirect
	golang.org/x/net v0.58.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/text v0.41.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260526163538-3dc84a4a5aaa // indirect
	google.golang.org/protobuf v1.36.12
)
