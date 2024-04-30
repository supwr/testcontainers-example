package database

import (
	"context"
	"fmt"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/gorm"
	"log"
	"log/slog"
	"os"
	"time"
)

const (
	databaseHost     = "postgres.players.integration"
	databaseUsername = "test_user"
	databasePassword = "test_password"
	databaseDBName   = "test_db"
	databasePort     = "5432"
	databaseSchema   = "sc_players"
	networkName      = "players_network"
)

type TestDatabase struct {
	Conn      *gorm.DB
	migration *Migration
	container testcontainers.Container
}

func SetupTestDatabase() *TestDatabase {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	container, host, port, err := createContainer(ctx)
	if err != nil {
		log.Fatal("failed to setup test", err)
	}

	config := Config{
		Environment:      "DEV",
		DatabaseHost:     host,
		DatabasePort:     port,
		DatabaseDBName:   databaseDBName,
		DatabaseSchema:   databaseSchema,
		DatabaseUsername: databaseUsername,
		DatabasePassword: databasePassword,
		MigrationsDir:    "/app/migrations/",
	}

	conn, err := NewConnection(config)

	if err != nil {
		log.Fatal("error connecting to database", err)
	}

	migration := NewMigration(conn, config, slog.New(slog.NewTextHandler(os.Stderr, nil)))
	migration.CreateSchema()
	migration.Migrate()

	return &TestDatabase{
		Conn:      conn,
		migration: migration,
		container: container,
	}
}

func (tdb *TestDatabase) CleanUp() {
	tdb.migration.DropSchema()
	tdb.migration.CreateSchema()
	tdb.migration.Migrate()
}

func (tdb *TestDatabase) TearDown() {
	_ = tdb.container.Terminate(context.Background())
}

func createContainer(ctx context.Context) (testcontainers.Container, string, string, error) {
	var port = fmt.Sprintf("%s/tcp", databasePort)
	var env = map[string]string{
		"POSTGRES_PASSWORD": databasePassword,
		"POSTGRES_USER":     databaseUsername,
		"POSTGRES_DB":       databaseDBName,
	}

	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:          "postgres:14-alpine",
			ExposedPorts:   []string{port},
			Networks:       []string{networkName},
			NetworkAliases: map[string][]string{networkName: {databaseHost}},
			Env:            env,
			WaitingFor: wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5 * time.Second),
		},
		Started: true,
	}

	container, err := testcontainers.GenericContainer(ctx, req)
	if err != nil {
		return container, "", "", err
	}

	host, err := container.Host(ctx)
	if err != nil {
		return container, "", "", err
	}

	p, err := container.MappedPort(ctx, databasePort)
	if err != nil {
		return container, "", "", err
	}

	log.Println("postgres container ready and running at port: ", p.Port())

	return container, host, p.Port(), nil
}
