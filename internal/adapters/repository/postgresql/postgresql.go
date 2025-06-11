package postgresql

import (
	stdsql "database/sql"
	// "frog-go/internal/adapters/repository/postgresql/hooks"
	"fmt"
	"frog-go/internal/core/ports/outbound/repository"
	"frog-go/internal/ent"
	"frog-go/internal/utils/logger"
	"net/url"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/lib/pq"
)

type PostgreSQL struct {
	log    *logger.Logger
	Client *ent.Client
	db     *stdsql.DB
}

func NewPostgreSQL(user, password, host, port, database, SeedPath string) (repository.Repository, error) {
	log := logger.NewLogger("PostgreSQL")

	dbURI := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		url.QueryEscape(user),
		url.QueryEscape(password),
		host,
		port,
		database,
	)

	drv, err := entsql.Open(dialect.Postgres, dbURI)
	if err != nil {
		return nil, err
	}

	sqlDB := drv.DB()
	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	client := ent.NewClient(ent.Driver(drv))
	// categorizer, err := hooks.NewCategorizer(SeedPath)
	// if err != nil {
	// 	return nil, err
	// }

	// client.Debt.Use(
	// 	hooks.UpdateInvoiceAmountHook(client),
	// 	hooks.SetCategoryFromTitleHook(client, categorizer),
	// )

	log.Start("Host: %s:%s | User: %s | DB: %s", host, port, user, database)

	return &PostgreSQL{Client: client, log: log, db: sqlDB}, nil
}

func (d *PostgreSQL) Close() {
	if err := d.Client.Close(); err != nil {
		d.log.Error("%v", err)
	} else {
		d.log.Info("Database connection closed.")
	}
}
