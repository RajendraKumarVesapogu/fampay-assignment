package connections

import (
	"context"
	"fmt"
	"net/url"

	"fampay-assignment/config"
	"fampay-assignment/logger"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/newrelic/go-agent/v3/integrations/nrpgx5"
)

var (
	postgresConnection *pgxpool.Pool
)

type PostgresCreds struct {
	Host        string
	Port        uint16
	Db          string
	
	Username string `json:"username"`
	Password string `json:"password"`
}

func GetPostgresDb() (db *pgxpool.Pool, ok bool) {
	if postgresConnection == nil {
		return nil, false
	}
	return postgresConnection, true
}

func getPostgresConnectionString(creds *PostgresCreds) string {
	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s",
		creds.Username,
		url.QueryEscape(creds.Password),
		creds.Host,
		creds.Port,
		creds.Db,
	)
}

func ConnectPostgres(creds *PostgresCreds) *pgxpool.Pool {
	connectionString := getPostgresConnectionString(creds)
	connectionConfig, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		logger.Log.WithFields(logger.Fields{
			"err": err,
		}).Error("failed to parse postgres connection string")
		return nil
	}

	connectionConfig.BeforeConnect = func(_ context.Context, config *pgx.ConnConfig) error {
		config.Tracer = nrpgx5.NewTracer(nrpgx5.WithQueryParameters(false))
		return nil
	}

	pool, err := pgxpool.NewWithConfig(ctx, connectionConfig)
	if err != nil {
		logger.Log.WithFields(logger.Fields{
			"host": creds.Host,
			"user": creds.Username,
			"port": creds.Port,
			"db":   creds.Db,
			"err":  err,
		}).Fatal("failed to connect to postgres")
	}
	logger.Log.WithFields(logger.Fields{
		"host": creds.Host,
		"user": creds.Username,
		"port": creds.Port,
		"db":   creds.Db,
	}).Info("connected to postgres")

	return pool
}

func init() {
	pgCreds:= &PostgresCreds{
		Host:        config.DataDbHost,
		Port:        uint16(config.DataDbPort),
		Db:          "postgres",
		Username:    config.DataDbUser,
		Password:    config.DataDbPassword,
	}
	ConnectPostgres(pgCreds)
}
