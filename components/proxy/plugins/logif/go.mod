module github.com/gitpod-io/gitpod/proxy/plugins/logif

go 1.17

replace github.com/gitpod-io/gitpod/proxy/plugins/jsonselect => ../jsonselect

require (
	github.com/PaesslerAG/gval v1.1.1
	github.com/buger/jsonparser v1.1.1
	github.com/caddyserver/caddy/v2 v2.4.3
	github.com/gitpod-io/gitpod/proxy/plugins/jsonselect v0.0.0-00010101000000-000000000000
	go.uber.org/zap v1.18.1
	golang.org/x/term v0.0.0-20210615171337-6886f2dfbf5b
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/caddyserver/certmagic v0.14.0 // indirect
	github.com/cespare/xxhash/v2 v2.1.1 // indirect
	github.com/dustin/go-humanize v1.0.1-0.20200219035652-afde56e7acac // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/uuid v1.2.0 // indirect
	github.com/klauspost/cpuid/v2 v2.0.6 // indirect
	github.com/libdns/libdns v0.2.1 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/mholt/acmez v0.1.3 // indirect
	github.com/miekg/dns v1.1.42 // indirect
	github.com/prometheus/client_golang v1.10.1-0.20210603120351-253906201bda // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.26.0 // indirect
	github.com/prometheus/procfs v0.6.0 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	golang.org/x/crypto v0.0.0-20210513164829-c07d793c2f9a // indirect
	golang.org/x/net v0.0.0-20210525063256-abc453219eb5 // indirect
	golang.org/x/sys v0.0.0-20210615035016-665e8c7367d1 // indirect
	golang.org/x/text v0.3.6 // indirect
	google.golang.org/protobuf v1.26.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
)
