module github.com/jekiti/citydrive/admin

go 1.24.5

require (
	github.com/jackc/pgx/v5 v5.7.6
	github.com/jekiti/citydrive v0.0.0-00010101000000-000000000000
	github.com/joho/godotenv v1.5.1
	github.com/redis/go-redis/v9 v9.14.0
	google.golang.org/grpc v1.75.1
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	golang.org/x/crypto v0.41.0 // indirect
	golang.org/x/net v0.43.0 // indirect
	golang.org/x/sync v0.17.0 // indirect
	golang.org/x/sys v0.36.0 // indirect
	golang.org/x/text v0.29.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250707201910-8d1bb00bc6a7 // indirect
	google.golang.org/protobuf v1.36.9 // indirect
)

replace github.com/jekiti/citydrive => ..
