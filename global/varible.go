package global

import (
	"github.com/jmoiron/sqlx"
)

// Db Global db connection
var Db *sqlx.DB
