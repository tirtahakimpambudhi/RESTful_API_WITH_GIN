package model

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose"
	_ "github.com/ziutek/mymysql/godrv"
)

const SupportDriver = "mysql postgres sqlite3 redshift"

type Connection struct {
	Driver        string
	DBString      string
	FileMigration string
	Dir           string
	Version       int64
}

func NewConnection(version int64, driver, dbstring, file, dir string) *Connection {
	return &Connection{
		Driver:        driver,
		DBString:      dbstring,
		FileMigration: file,
		Dir:           dir,
		Version:       version,
	}
}
func (c *Connection) IsDriverSupported() bool {
	return strings.Contains(SupportDriver, c.Driver)
}

func (c *Connection) OpenDB() (*sql.DB, error) {
	if !c.IsDriverSupported() {
		return nil, fmt.Errorf("Driver Not Supported")
	}

	if err := goose.SetDialect(c.Driver); err != nil {
		return nil, err
	}

	db, err := sql.Open(c.Driver, c.DBString)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (c *Connection) Run(command string) error {
	db, err := c.OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()

	switch command {
	case "create":
		args := strings.Split(c.FileMigration, ".")
		if err := goose.Run(command, db, c.Dir, args...); err != nil {
			return err
		}
	case "up-to":
		if err := goose.UpTo(db, c.Dir, c.Version); err != nil {
			return err
		}
	case "down-to":
		if err := goose.DownTo(db, c.Dir, c.Version); err != nil {
			return err
		}
	default:
		if err := goose.Run(command, db, c.Dir); err != nil {
			return err
		}
	}

	return nil
}
