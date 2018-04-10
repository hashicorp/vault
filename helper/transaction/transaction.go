package transaction

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

type Config struct {
	Ctx        context.Context
	Name       string
	Username   string
	Password   string
	Expiration string
	Tx         *sql.Tx
	DB         *sql.DB
}

func Execute(c *Config, query string) error {

	if err := validate(c); err != nil {
		return err
	}

	stmt, err := statement(c, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	return execute(c, stmt)
}

func validate(c *Config) error {
	if c.DB == nil && c.Tx == nil {
		return errors.New("a *sql.Tx or *sql.DB must be provided to prepare the statement")
	}
	if c.DB != nil && c.Tx != nil {
		return errors.New("cannot provide both a *sql.Tx and a *sql.DB because only one must be used to prepare the statement")
	}
	return nil
}

func statement(c *Config, query string) (*sql.Stmt, error) {

	q := parseQuery(c, query)

	if c.Tx != nil {
		if c.Ctx != nil {
			return c.Tx.PrepareContext(c.Ctx, q)
		}
		return c.Tx.Prepare(q)
	}

	if c.Ctx != nil {
		return c.DB.PrepareContext(c.Ctx, q)
	}
	return c.DB.Prepare(q)
}

func execute(c *Config, stmt *sql.Stmt) error {

	if c.Ctx != nil {
		if _, err := stmt.ExecContext(c.Ctx); err != nil {
			return err
		}
		return nil
	}

	if _, err := stmt.Exec(); err != nil {
		return err
	}
	return nil
}

func parseQuery(c *Config, tpl string) string {

	if c.Name == "" && c.Username == "" && c.Password == "" && c.Expiration == "" {
		return tpl
	}

	data := make(map[string]string)
	if c.Name != "" {
		data["name"] = c.Name
	}
	if c.Username != "" {
		data["username"] = c.Username
	}
	if c.Password != "" {
		data["password"] = c.Password
	}
	if c.Expiration != "" {
		data["expiration"] = c.Expiration
	}

	for k, v := range data {
		tpl = strings.Replace(tpl, fmt.Sprintf("{{%s}}", k), v, -1)
	}

	return tpl
}
