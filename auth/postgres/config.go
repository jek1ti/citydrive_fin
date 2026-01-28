package postgres

type PostgresConfig struct {
	Host     string
	Database string
	User     string
	Password string
	MaxConn  int
}

const (
	defaultConnect = 2
)
