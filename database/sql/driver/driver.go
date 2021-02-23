package driver

import (
	"database/sql/driver"

	"github.com/go-sql-driver/mysql"
)

// Driver ...
type Driver struct {
	BaseName string
}

// Open ...
func (d Driver) Open(dns string) (driver.Conn, error) {
	dr := d.getBaseDriver()
	baseconn, err := dr.Open(dns)
	if err != nil {
		return nil, err
	}

	c := &conn{
		base: baseconn,
	}
	return c, nil
}

func (d Driver) getBaseDriver() driver.Driver {
	switch d.BaseName {
	case "mysql":
		return &mysql.MySQLDriver{}
	}
	return nil
}
