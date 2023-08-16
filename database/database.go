package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"api-boilerplate/models/domain"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
)

type DB struct {
	db     *sqlx.DB
	ctx    context.Context
	cfg    domain.Cfg
	cancel func()
}

// Tx wraps the SQL Tx object to provide a timestamp at the start of the transaction.
type Tx struct {
	*sqlx.Tx
	db  *DB
	now time.Time
}

// NewDB initials the database
func NewDB(cfg domain.Cfg) *DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable timezone=CET",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.Dbname)

	db, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		log.Fatalln(err)
	}

	migrations := &migrate.FileMigrationSource{
		Dir: "database/migrations",
	}

	migrate.SetTable("migrations")

	n, err := migrate.Exec(db.DB, "postgres", migrations, migrate.Up)
	if err != nil {
		log.Println(err)
	}

	log.Printf("Applied %d migrations!\n", n)
	log.Println("Connected to database")

	database := &DB{}

	database.ctx, database.cancel = context.WithCancel(context.Background())
	database.db = db
	database.cfg = cfg
	return database
}

type Queryable interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Get(dest interface{}, query string, args ...interface{}) error
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Rebind(query string) string
	BindNamed(query string, arg interface{}) (string, []interface{}, error)
	Select(dest interface{}, query string, args ...interface{}) error
	DriverName() string
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	NamedExec(query string, arg interface{}) (sql.Result, error)
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)
	PrepareNamed(query string) (*sqlx.NamedStmt, error)
	PrepareNamedContext(ctx context.Context, query string) (*sqlx.NamedStmt, error)
	Preparex(query string) (*sqlx.Stmt, error)
	PreparexContext(ctx context.Context, query string) (*sqlx.Stmt, error)
	QueryRowx(query string, args ...interface{}) *sqlx.Row
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
	Queryx(query string, args ...interface{}) (*sqlx.Rows, error)
	QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
	MustExec(query string, args ...interface{}) sql.Result
	MustExecContext(ctx context.Context, query string, args ...interface{}) sql.Result
}

// BeginTx starts a transaction and returns a wrapper Tx type. This type
// provides a reference to the database and a fixed timestamp at the start of
// the transaction. The timestamp allows us to mock time during tests as well.
func (db *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	tx, err := db.db.BeginTxx(ctx, opts)
	if err != nil {
		return nil, err
	}

	// Return wrapper Tx that includes the transaction start time.
	return &Tx{
		Tx: tx,
		db: db,
	}, nil
}

//func (db *DB) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
//	return db.db.GetContext(ctx, dest, query, args...)
//}
