module github.com/skycoin/skycoin-service-discovery

go 1.17

require (
	github.com/go-chi/chi v4.1.2+incompatible
	github.com/go-chi/cors v1.1.1
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/lib/pq v1.10.7
	github.com/sirupsen/logrus v1.8.1
	github.com/skycoin/dmsg v1.3.0-rc1.0.20230105101327-c8f2541c3de8
	github.com/skycoin/skycoin v0.27.1
	github.com/skycoin/skywire v1.3.2-0.20230124214138-b3f0bcfb7e77
	github.com/skycoin/skywire-utilities v0.0.0-20230314131305-bdd8e27cbd46
	github.com/spf13/cobra v1.4.0
	github.com/stretchr/testify v1.8.0
	gorm.io/driver/postgres v1.3.8
	gorm.io/gorm v1.23.8
)

require (
	github.com/AudriusButkevicius/pfilter v0.0.0-20210515103320-4b4b86609d51 // indirect
	github.com/VictoriaMetrics/metrics v1.18.1 // indirect
	github.com/ccding/go-stun/stun v0.0.0-20200514191101-4dc67bcdb029 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/go-chi/chi/v5 v5.0.8-0.20220103230436-7dbe9a0bd10f // indirect
	github.com/go-redis/redis/v8 v8.11.5 // indirect
	github.com/gocarina/gocsv v0.0.0-20220927221512-ad3251f9fa25 // indirect
	github.com/google/uuid v1.1.2 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.12.1 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20200714003250-2b9c44734f2b // indirect
	github.com/jackc/pgtype v1.11.0 // indirect
	github.com/jackc/pgx/v4 v4.16.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid v1.2.4 // indirect
	github.com/klauspost/reedsolomon v1.9.9 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d // indirect
	github.com/mmcloughlin/avo v0.0.0-20200523190732-4439b6b2c061 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pires/go-proxyproto v0.6.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/skycoin/noise v0.0.0-20180327030543-2492fe189ae6 // indirect
	github.com/skycoin/yamux v0.0.0-20200803175205-571ceb89da9f // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stretchr/objx v0.4.0 // indirect
	github.com/templexxx/cpufeat v0.0.0-20180724012125-cef66df7f161 // indirect
	github.com/templexxx/xor v0.0.0-20191217153810-f85b25db303b // indirect
	github.com/tjfoc/gmsm v1.4.0 // indirect
	github.com/valyala/fastrand v1.1.0 // indirect
	github.com/valyala/histogram v1.2.0 // indirect
	github.com/xtaci/kcp-go v5.4.20+incompatible // indirect
	go.etcd.io/bbolt v1.3.6 // indirect
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519 // indirect
	golang.org/x/mod v0.6.0-dev.0.20220419223038-86c51ed26bb4 // indirect
	golang.org/x/net v0.0.0-20211020060615-d418f374d309 // indirect
	golang.org/x/sys v0.0.0-20220704084225-05e143d24a9e // indirect
	golang.org/x/term v0.0.0-20210927222741-03fcf44c2211 // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/tools v0.1.11 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// Uncomment it for tests with alternative branches and run `go mod vendor`
//replace github.com/skycoin/dmsg => ../dmsg

// replace github.com/skycoin/skywire => ../skywire

// replace github.com/skycoin/skywire-utilities => ../skywire-utilities
