package database

import (
	"fitness-club/models"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() error {
	dsn := "host=localhost user=admin password=password dbname=fitness_club port=5432 sslmode=disable TimeZone=UTC"

	var err error
	var db *gorm.DB

	// Попытки подключения с задержкой
	for i := 0; i < 5; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Printf("Attempt %d: Failed to connect to database: %v", i+1, err)
			time.Sleep(3 * time.Second)
			continue
		}
		break
	}

	if err != nil {
		return fmt.Errorf("failed to connect to database after 5 attempts: %w", err)
	}

	DB = db
	log.Println("Successfully connected to database")
	return nil
}

func Migrate() error {
	err := DB.AutoMigrate(
		&models.Member{},
		&models.MembershipType{},
		&models.Membership{},
		&models.Trainer{},
		&models.TrainingSession{},
	)
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database migrated successfully")
	return nil
}
