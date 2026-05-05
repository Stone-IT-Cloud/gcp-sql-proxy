package proxy

import "fmt"

const (
	defaultPGHost = "127.0.0.1"
)

func BuildPostgresClientTemplate(host string, port int, dbName string, dbUser string) string {
	if host == "" {
		host = defaultPGHost
	}
	if dbName == "" {
		dbName = "<db_name>"
	}
	if dbUser == "" {
		dbUser = "<db_user>"
	}
	return fmt.Sprintf("psql \"host=%s port=%d dbname=%s user=%s sslmode=disable\"",
		host, port, dbName, dbUser)
}
