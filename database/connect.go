package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/polarisjrex0406/federico-app/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	dsnWithoutDb := connectionString(cfg, false)
	if dsnWithoutDb == "" {
		log.Fatal("incorrect database configuration")
	}

	// Open a connection to the PostgreSQL server (without specifying a database)
	conn, err := sql.Open("postgres", dsnWithoutDb)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	dbName := cfg.Postgres.DBName
	if err := checkDatabaseExists(conn, dbName); err != nil {
		if err := createDatabase(conn, dbName); err != nil {
			log.Fatal(err)
		}
	}

	dsn := connectionString(cfg, true)
	// Open the database connection
	DB, err = gorm.Open(
		postgres.Open(dsn),
		&gorm.Config{
			TranslateError: true,
		},
	)
	if err != nil {
		log.Fatal(err)
	}
}

func connectionString(cfg *config.Config, dbSpecified bool) string {
	// Construct the connection string
	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%d sslmode=%s",
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.SSLMode,
	)

	if dbSpecified {
		dsn = fmt.Sprintf("%s dbname=%s", dsn, cfg.Postgres.DBName)
	}
	return dsn
}

// checkDatabaseExists checks if the specified database exists.
func checkDatabaseExists(conn *sql.DB, dbName string) error {
	var exists bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname='%s')", dbName)
	err := conn.QueryRow(query).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("database %s does not exist", dbName)
	}
	return nil
}

// createDatabase creates the specified database.
func createDatabase(conn *sql.DB, dbName string) error {
	query := fmt.Sprintf("CREATE DATABASE %s", dbName)
	_, err := conn.Exec(query)
	if err != nil {
		return err
	}
	log.Printf("Database %s created successfully", dbName)
	return nil
}
