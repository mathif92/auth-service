package testutil

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"

	_ "github.com/go-sql-driver/mysql"
)

var (
	db *sqlx.DB
)

func TestMain(m *testing.M) {
	// Setup the database before running tests
	setupDatabase()

	// Run the tests
	exitCode := m.Run()

	// Tear down the database after running the tests
	teardownDatabase()

	// Exit with the appropriate exit code
	os.Exit(exitCode)
}

func setupDatabase() {
	q := make(url.Values)
	q.Set("tls", "false")
	q.Set("loc", "UTC")
	q.Set("parseTime", "true")
	q.Set("timeout", "1s")
	q.Set("readTimeout", "2s")
	q.Set("writeTimeout", "2s")

	connStr := fmt.Sprintf("%s:%s@tcp(%s)/%s?%s", "auth", "auth", "0.0.0.0:3307", "auth", q.Encode())

	conn, err := sqlx.Open("mysql", connStr)
	if err != nil {
		log.Fatalf("Could not connect to the DB: %s", err)
	}

	db = conn
}

func teardownDatabase() {
	db.Close()
}
