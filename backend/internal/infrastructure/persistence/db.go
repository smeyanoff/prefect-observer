package persistence

import (
	"crm-uplift-ii24-backend/config"
	"crm-uplift-ii24-backend/internal/domain/entity"
	"crm-uplift-ii24-backend/pkg/logging"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ConnectDB устанавливает соединение с базой данных PostgreSQL
func ConnectDB(cfg *config.Config) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	maxRetries := 10
	retryDelay := time.Duration(1 * time.Second)

	dsn := fmt.Sprintf("host=%s user=%s password=%s database=%s port=%s sslmode=%s",
		cfg.DB.Host, cfg.DB.User, cfg.DB.Pwd, cfg.DB.Database, cfg.DB.Port, cfg.DB.SSLMode)

	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
		})
		if err == nil {
			logging.Info("Connected to the database successfully")
			return db, nil
		}

		logging.Warn("Connect to the database failed.", zap.Int("retry", i))
		time.Sleep(retryDelay) // Задержка перед повторной попыткой
	}
	return nil, fmt.Errorf("could not connect to the database after %d attempts: %w", maxRetries, err)
}

// AutoMigrate выполняет автоматическую миграцию всех моделей
func AutoMigrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&entity.Sendpost{},
		&entity.Stage{},
		&entity.SendpostSchedule{},
	); err != nil {
		return err
	}

	if err := db.Migrator().CreateConstraint(&entity.Sendpost{}, "FirstStage"); err != nil {
		return err
	}
	if err := db.Migrator().CreateConstraint(&entity.Stage{}, "Sendpost"); err != nil {
		return err
	}
	if err := db.Migrator().CreateConstraint(&entity.Stage{}, "NextStage"); err != nil {
		return err
	}
	if err := db.Migrator().CreateConstraint(&entity.Stage{}, "SubStages"); err != nil {
		return err
	}

	logging.Logger.Info("Database migration completed successfully")
	return nil
}
