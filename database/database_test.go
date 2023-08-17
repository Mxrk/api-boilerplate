package database

import (
	"log"
	"os"
	"strconv"
	"testing"
	"time"

	"api-boilerplate/models/domain"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	db  *DB
	m   *migrate.Migration
	cfg domain.Cfg
}

var testDB *DB
var testCFG domain.Cfg

func TestMain(m *testing.M) {
	log.SetFlags(log.LstdFlags | log.Llongfile)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "latest",
		Env: []string{
			"POSTGRES_PASSWORD=secret",
			"POSTGRES_USER=user_name",
			"POSTGRES_DB=dbname",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}
	resource.Expire(120) // Tell docker to hard kill the container in 120 seconds

	pool.MaxWait = 120 * time.Second
	if err = pool.Retry(func() error {
		port, err := strconv.Atoi(resource.GetPort("5432/tcp"))
		if err != nil {
			log.Fatalf("Could not convert port: %s", err)
		}

		var test = domain.Cfg{}
		test.Server.Port = "8080"
		test.Database.Host = "localhost"
		test.Database.Port = port
		test.Database.User = "user_name"
		test.Database.Password = "secret"
		test.Database.Dbname = "dbname"

		log.Println("Connecting to database")
		db, err := NewDB(test, "migrations")
		if err != nil {
			return err
		}

		testDB = db
		testCFG = test

		return nil
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	code := m.Run()

	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestExampleTestSuite(t *testing.T) {
	log.Println("Running tests")
	suite.Run(t, new(TestSuite))
}

func (suite *TestSuite) SetupTest() {
	suite.db = testDB
	suite.cfg = testCFG
}
