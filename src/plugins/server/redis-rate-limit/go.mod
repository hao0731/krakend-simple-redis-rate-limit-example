module redis-rate-limit

go 1.22.7

require (
	github.com/go-redis/redis_rate/v10 v10.0.1
	github.com/redis/go-redis/v9 v9.0.2
)

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
)

replace github.com/cespare/xxhash/v2 => github.com/cespare/xxhash/v2 v2.3.0
