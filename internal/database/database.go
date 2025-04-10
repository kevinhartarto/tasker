package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/kevinhartarto/tasker/internal/logger"
	"github.com/kevinhartarto/tasker/internal/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Database interface {

	// Close terminates the database connection.
	// It returns an error if connection cannot be closed.
	Close()

	// Get active GORM connection
	// Return an active GORM connection
	Gorm() *gorm.DB
}

type database struct {
	connection *sql.DB
	gorm       *gorm.DB
}

var (
	gormService *database
	log         = logger.GetLogger()
)

func Start() Database {
	// Keep the connection alive
	if gormService != nil {
		return gormService
	}

	pgx, err := sql.Open("pgx", getDBConnection())
	if err != nil {
		log.Error("Failed to establish database connection, closing...", "message: ", err)
		pgx.Close()
		os.Exit(1)
	}

	// Initialize GORM
	taskerDialector := postgres.New(postgres.Config{
		Conn: pgx,
	})

	taskerConfig := gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "tasker.",
			SingularTable: true,
		},
	}

	gormDB, err := gorm.Open(taskerDialector, &taskerConfig)
	if err != nil {
		log.Error("Failed to initialize GORM, closing...", "message: ", err)
		pgx.Close()
		os.Exit(1)
	}

	gormService = &database{
		connection: pgx,
		gorm:       gormDB,
	}

	log.Info("Database connection established, GORM is running")
	return gormService
}

func (db *database) Gorm() *gorm.DB {
	return db.gorm
}

func (db *database) Close() {
	log.Info("Closing database connection.")
	if gormService != nil {
		if err := gormService.connection.Close(); err != nil {
			log.Info("Database close error", "message: ", err)
		} else {
			log.Info("Database connection closed.")
		}
	}
}

func getDBConnection() string {
	database := utils.GetEnvOrDefault("DB_DATABASE", "devstack")
	username := utils.GetEnvOrDefault("DB_USERNAME", "developer")
	password := utils.GetEnvOrDefault("DB_PASSWORD", "localDevstack01")

	port := utils.GetEnvOrDefault("DB_PORT", "5432")
	host := utils.GetEnvOrDefault("DB_HOST", "localhost")
	ssl := utils.GetEnvOrDefault("SSL_MODE", "disable")

	credential := fmt.Sprintf("%v:%v", username, password)
	databaseAddr := fmt.Sprintf("%v:%v/%v", host, port, database)
	connOptions := fmt.Sprintf("sslmode=%v", ssl)
	return fmt.Sprintf("postgres://%v@%v?%v", credential, databaseAddr, connOptions)
}
