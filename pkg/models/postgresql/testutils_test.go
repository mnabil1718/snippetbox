package postgresql

import (
	"database/sql"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func OpenDB(connString string) (*sql.DB, error) {
	db, err := sql.Open("pgx", connString)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func NewTestDB(t *testing.T) *sql.DB {
	connString := "postgres://mnabil:Cucibaju123@localhost:5432/test_snippetbox"
	db, err := OpenDB(connString)
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func TruncateTable(t *testing.T, db *sql.DB, tableName string) {
	_, err := db.Exec("TRUNCATE " + tableName)
	if err != nil {
		t.Fatal(err)
	}

	// Reset the associated sequence (assuming the column name is "id")
	_, err = db.Exec("ALTER SEQUENCE " + tableName + "_id_seq RESTART WITH 1")
	if err != nil {
		t.Fatal(err)
	}
}

func createAndInsertTable(t *testing.T, tableName string, statement string, args ...any) (*sql.DB, func()) {

	db := NewTestDB(t)
	TruncateTable(t, db, tableName) // make sure the table is empty initially
	_, err := db.Exec(statement, args...)
	if err != nil {
		t.Fatal(err)
	}

	return db, func() {
		TruncateTable(t, db, tableName) // make sure table is clear after use
		db.Close()
	}
}
