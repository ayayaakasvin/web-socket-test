package postgresql

import (
	"database/sql"
	"fmt"

	goshutdownchannel "github.com/ayayaakasvin/go-shutdown-channel"
	"github.com/ayayaakasvin/web-socket-test/internal/config"

	_ "github.com/lib/pq"
)

const origin = "PostgreSQL"

type PostgreSQL struct {
	conn *sql.DB
}

func NewPostgreSQLConnection(dbConfig config.PostgreSQLConfig, s *goshutdownchannel.Shutdown) *PostgreSQL {
	psql := new(PostgreSQL)

	connection, err := sql.Open("postgres", dbConfig.URL)
	if err != nil {
		msg := fmt.Sprintf("failed to connect to db: %v\n", err)
		s.Send(origin, msg)
		return nil
	}

	psql.conn = connection

	if err := psql.conn.Ping(); err != nil {
		msg := fmt.Sprintf("failed to ping to db: %v\n", err)
		s.Send(origin, msg)
		return nil
	}

	return psql
}

func (p *PostgreSQL) Close() error {
	return p.Close()
}
