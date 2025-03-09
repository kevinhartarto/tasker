package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Service interface {

	// Close terminates the database connection.
	// It returns an error if connection cannot be closed.
	Close() error

	// Get active GORM connection
	// Return an active GORM connection
	UseGorm() *gorm.DB
}

type service struct {
	sqlDB  *sql.DB
	gormDB *gorm.DB
}

var (
	database   = os.Getenv("DB_DATABASE")
	username   = os.Getenv("DB_USERNAME")
	password   = os.Getenv("DB_PASSWORD")
	port       = os.Getenv("DB_PORT")
	host       = os.Getenv("DB_HOST")
	dbInstance *service
)

func StartDB() Service {
	// Keep the connection alive
	if dbInstance != nil {
		return dbInstance
	}

	credential := fmt.Sprintf("%v:%v", username, password)
	databaseUrl := fmt.Sprintf("%v:%v/%v", host, port, database)
	connOptions := fmt.Sprintf("%v", "sslmode=disable")
	conn := fmt.Sprintf("postgres://%v@%v?%v", credential, databaseUrl, connOptions)

	sqlDB, err := sql.Open("pgx", conn)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize GORM
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "utils.",
			SingularTable: true,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// Ping to check if the connection is still alive
	if err = sqlDB.Ping(); err != nil {
		log.Println("Database connection lost")
		log.Fatal(err)
	}

	dbInstance = &service{
		sqlDB:  sqlDB,
		gormDB: gormDB,
	}

	fmt.Println("Database connection established.")

	return dbInstance
}

func (s *service) UseGorm() *gorm.DB {
	return s.gormDB
}

func (s *service) Close() error {
	log.Printf("Disconnected from database: %s", database)
	return s.sqlDB.Close()
}
