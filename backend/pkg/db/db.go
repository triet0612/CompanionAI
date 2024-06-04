package db

import (
	"CompanionBackend/pkg/config"
	"context"
	_ "embed"
	"log"

	"github.com/jackc/pgx/v5"
)

type DBHelper struct {
	*pgx.Conn
}

//go:embed schema.sql
var schemaFile string

//go:embed populate.sql
var populate string

func Init(cfg *config.Config) *DBHelper {
	ctx := context.Background()

	con, err := pgx.Connect(ctx, cfg.DB_URL)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := con.Exec(ctx, schemaFile); err != nil {
		log.Fatal(err)
	}
	if _, err := con.Exec(ctx, populate); err != nil {
		log.Fatal(err)
	}
	row, err := con.Query(context.Background(), "SELECT ConfigKey, ConfigValue FROM Config")
	if err != nil {
		log.Fatal(err)
	}
	for row.Next() {
		key, val := "", ""
		if err := row.Scan(&key, &val); err != nil {
			log.Fatal(err)
		}
		cfg.Dynamic[key] = val
	}
	return &DBHelper{Conn: con}
}
