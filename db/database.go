package db

import (
	"database/sql"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"

	// Using the blank identifier in order to solely
	// provide the side-effects of the package.
	// Eseentially the side effect is calling the `init()`
	// method of `lib/pq`:
	//	func init () {  sql.Register("postgres", &Driver{} }
	// which you can see at `github.com/lib/pq/conn.go`
	_ "github.com/lib/pq"
)

// Roach holds the connection pool to the database - created by a configuration
// object (`Config`).
type Roach struct {
	// Db holds a sql.DB pointer that represents a pool of zero or more
	// underlying connections - safe for concurrent use by multiple
	// goroutines -, with freeing/creation of new connections all managed
	// by `sql/database` package.
	Db  *sql.DB
	cfg Config
}

// Config holds the configuration used for instantiating a new Roach.
type Config struct {
	// Address that locates our postgres instance
	Host string
	// Port to connect to
	Port string
	// User that has access to the database
	User string
	// Password so that the user can login
	Password string
	// Database to connect to (must have been created priorly)
	Database string
}

// New returns a Roach with the sql.DB set with the postgres
// DB connection string in the configuration
func New(cfg Config) (roach Roach, err error) {
	if cfg.Host == "" || cfg.Port == "" || cfg.User == "" ||
		cfg.Password == "" || cfg.Database == "" {
		err = errors.Errorf(
			"All fields must be set (%s)",
			spew.Sdump(cfg))
		return
	}

	roach.cfg = cfg

	// The first argument corresponds to the driver name that the driver
	// (in this case, `lib/pq`) used to register itself in `database/sql`.
	// The next argument specifies the parameters to be used in the connection.
	// Details about this string can be seen at https://godoc.org/github.com/lib/pq
	db, err := sql.Open("postgres", fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		cfg.User, cfg.Password, cfg.Database, cfg.Host, cfg.Port))
	if err != nil {
		err = errors.Wrapf(err,
			"Couldn't open connection to postgre database (%s)",
			spew.Sdump(cfg))
		return
	}

	// Ping verifies if the connection to the database is alive or if a
	// new connection can be made.
	if err = db.Ping(); err != nil {
		err = errors.Wrapf(err,
			"Couldn't ping postgre database (%s)",
			spew.Sdump(cfg))
		return
	}

	roach.Db = db
	return
}

// Close performs the release of any resources that
// `sql/database` DB pool created. This is usually meant
// to be used in the exiting of a program or `panic`ing.
func (r *Roach) Close() (err error) {
	if r.Db == nil {
		return
	}

	if err = r.Db.Close(); err != nil {
		err = errors.Wrapf(err,
			"Errored closing database connection",
			spew.Sdump(r.cfg))
	}

	return
}
