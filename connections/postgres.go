package connections

import (
    "context"
    "fmt"
    "net/url"
    "time"

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
    Host     string
    Port     uint16
    Db       string
    Username string
    Password string
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
        url.QueryEscape(creds.Username),
        url.QueryEscape(creds.Password),
        creds.Host,
        creds.Port,
        creds.Db,
    )
}

func ConnectPostgres(creds *PostgresCreds) *pgxpool.Pool {
    // Create a context with timeout for connection
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    connectionString := getPostgresConnectionString(creds)
    connectionConfig, err := pgxpool.ParseConfig(connectionString)
    if err != nil {
        logger.Log.WithFields(logger.Fields{
            "err": err,
        }).Fatal("failed to parse postgres connection string")
        return nil
    }

    connectionConfig.MaxConns = 10
    connectionConfig.MinConns = 2
    connectionConfig.MaxConnLifetime = time.Hour
    connectionConfig.MaxConnIdleTime = 30 * time.Minute
    connectionConfig.HealthCheckPeriod = 1 * time.Minute

    connectionConfig.BeforeConnect = func(_ context.Context, config *pgx.ConnConfig) error {
        config.Tracer = nrpgx5.NewTracer(nrpgx5.WithQueryParameters(false))
        config.ConnectTimeout = 10 * time.Second
        return nil
    }

    var pool *pgxpool.Pool
    maxRetries := 3
    for attempt := 0; attempt < maxRetries; attempt++ {
        pool, err = pgxpool.NewWithConfig(ctx, connectionConfig)
        if err == nil {
            break
        }

        if attempt < maxRetries-1 {
            logger.Log.WithFields(logger.Fields{
                "attempt": attempt + 1,
                "err":     err,
            }).Warning("failed to connect to postgres, retrying...")
            time.Sleep(time.Duration(attempt+1) * time.Second)
            continue
        }

        logger.Log.WithFields(logger.Fields{
            "host": creds.Host,
            "user": creds.Username,
            "port": creds.Port,
            "db":   creds.Db,
            "err":  err,
        }).Fatal("failed to connect to postgres after all retries")
        return nil
    }

    if err := pool.Ping(ctx); err != nil {
        logger.Log.WithFields(logger.Fields{
            "err": err,
        }).Fatal("failed to ping postgres after connection")
        return nil
    }

    logger.Log.WithFields(logger.Fields{
        "host": creds.Host,
        "user": creds.Username,
        "port": creds.Port,
        "db":   creds.Db,
    }).Info("successfully connected to postgres")

    return pool
}

func init() {
    pgCreds := &PostgresCreds{
        Host:     config.DataDbHost,
        Port:     uint16(config.DataDbPort),
        Db:       "postgres",
        Username: config.DataDbUser,
        Password: config.DataDbPassword,
    }

    postgresConnection = ConnectPostgres(pgCreds)

    if postgresConnection == nil {
        logger.Log.Fatal("failed to initialize postgres connection")
    }
}