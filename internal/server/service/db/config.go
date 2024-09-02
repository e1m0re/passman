package db

type DatabaseConfig struct {
	Driver                  string
	Url                     string
	ConnMaxLifetimeInMinute int
	MaxOpenConnections      int
	MaxIdleConnections      int
}
