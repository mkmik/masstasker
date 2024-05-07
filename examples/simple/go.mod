module mkm.pub/masstasker/examples/simple

go 1.22

replace mkm.pub/masstasker => ../..

require (
	github.com/prometheus/client_golang v1.19.0
	google.golang.org/grpc v1.62.1
	mkm.pub/masstasker v0.1.13
)

require (
	github.com/alecthomas/kong v0.9.0
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/jonboulle/clockwork v0.4.0 // indirect
	github.com/prometheus/client_model v0.5.0 // indirect
	github.com/prometheus/common v0.48.0 // indirect
	github.com/prometheus/procfs v0.12.0 // indirect
	golang.org/x/net v0.23.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240308144416-29370a3891b7 // indirect
	google.golang.org/protobuf v1.32.0
)
