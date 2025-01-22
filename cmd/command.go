package cmd

import (
	"log"

	"github.com/polarisjrex0406/federico-app/config"
	"github.com/polarisjrex0406/federico-app/migrations"
	"gorm.io/gorm"
)

func Commands(db *gorm.DB) bool {
	cfg, err := config.GetConfig()
	if err != nil {
		return false
	}

	if cfg.Command.Migrate {
		if err := migrations.Migrate(db); err != nil {
			log.Fatalf("error migration: %v", err)
		}
		log.Println("migration completed successfully")
	}

	if cfg.Command.Seed {
		if err := migrations.Seeder(db); err != nil {
			log.Fatalf("error migration seeder: %v", err)
		}
		log.Println("seeder completed successfully")
	}

	return true
}
