package database

import (
	"database/sql"
	"deller/internal/config"
	"fmt"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

var db *bun.DB

func InitDatabase() *bun.DB {
	Host, Port, User, Password, Database := config.GetSQL()
	sqldb := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithDSN(fmt.Sprintf(
			"postgres://%v:%v@%v:%v/%v?sslmode=disable",
			User, Password, Host, Port, Database))))

	db = bun.NewDB(sqldb, pgdialect.New())

	if err := db.Ping(); err != nil {
		panic("cannot connect to db " + Database)
	}
	return db
}

func GetDatabase() *bun.DB {
	return db
}
