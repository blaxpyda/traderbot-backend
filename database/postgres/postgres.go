package postgres

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"thugcorp.io/final_bot/logger"
)

var (
	ErrMissingDBEnvVars = fmt.Errorf("missing required database environment variables")

)

type postgresDB struct {
	db *gorm.DB
	logger *logger.Logger
}

type PostgresDB interface {
	GetDB() *gorm.DB
	Connect() error
	IsConnected() bool
}

func NewPostgresDB(logger *logger.Logger) (PostgresDB, error) {
	db := &postgresDB{
		logger: logger,
	}
	if err := db.Connect(); err != nil {
		return nil, err
	}
	return db, nil
}



func (p *postgresDB) Connect() error {
	if p.db != nil {
		return nil
	}

	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		// Validate required environment variables
		host := os.Getenv("DB_HOST")
		port := os.Getenv("DB_PORT")
		user := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASSWORD")
		dbname := os.Getenv("DB_NAME")

		if host == "" || port == "" || user == "" || password == "" || dbname == "" {
			return ErrMissingDBEnvVars
		}

		// Construct DSN
		sslmode := os.Getenv("DB_SSLMODE")
		if sslmode == "" {
			sslmode = "disable"
		}
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, sslmode)
	
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		p.logger.Error("Failed to connect to Postgres database", err)
		return err
	}

	// set connection pool settings
	sqlDB, err := db.DB()
	if err != nil {
		p.logger.Error("Failed to get database instance", err)
		return nil
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(30)

	p.db = db
	return nil
}

func (p *postgresDB) GetDB() *gorm.DB { 
	return p.db
}

func (p *postgresDB) IsConnected() bool {
	return p.db != nil
}

func (p *postgresDB) Close() error {
	if p.db != nil {
		sqlDB, err := p.db.DB()
		if err != nil {
			p.logger.Error("Failed to get database instance for closing", err)
			return err
		}
		return sqlDB.Close()
	}
	return nil
}