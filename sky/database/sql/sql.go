package sql

import (
	"database/sql"

	"github.com/WiFeng/go-sky/sky/database/sql/driver"
)

var (
	driverMap = make(map[string]*driver.Driver)
)

// Register register wrapper Driver and return the new wrapper DriverName
func Register(driverName string) string {
	dn := "sky" + driverName

	if _, ok := driverMap[dn]; ok {
		return dn
	}

	dr := &driver.Driver{
		BaseName: driverName,
	}
	driverMap[dn] = dr
	sql.Register(dn, dr)
	return dn
}
