module github.com/jekiti/citydrive/auth

go 1.24.5

require (
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/jackc/pgx/v5 v5.7.6
	github.com/jekiti/citydrive v0.0.0
	github.com/joho/godotenv v1.5.1
	github.com/sethvargo/go-password v0.3.1
	golang.org/x/crypto v0.42.0
	google.golang.org/grpc v1.75.1
)

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/stretchr/testify v1.11.1 // indirect
	golang.org/x/net v0.43.0 // indirect
	golang.org/x/sync v0.17.0 // indirect
	golang.org/x/sys v0.36.0 // indirect
	golang.org/x/text v0.29.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250707201910-8d1bb00bc6a7 // indirect
	google.golang.org/protobuf v1.36.9 // indirect
)

replace github.com/jekiti/citydrive => ..
