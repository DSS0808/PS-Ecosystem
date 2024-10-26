package app

import (
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"strconv"

	"github.com/VikaPaz/pantheon/internal/models"
	"github.com/VikaPaz/pantheon/internal/repository"
	"github.com/VikaPaz/pantheon/internal/server/grpc"
	"github.com/VikaPaz/pantheon/internal/server/rest"
	service "github.com/VikaPaz/pantheon/internal/srvice"
	"github.com/joho/godotenv"
	"github.com/pressly/goose"
	"github.com/sirupsen/logrus"
)

type S struct {
}

func Run() {
	logger := NewLogger(logrus.DebugLevel, &logrus.TextFormatter{
		FullTimestamp: true,
	})

	if err := godotenv.Overload(); err != nil {
		logger.Errorf("Error loading .env file: %e", models.ErrLoadEnvFailed)
		return
	}

	confPostgres := repository.Config{
		Host:     os.Getenv("HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		User:     os.Getenv("USER"),
		Password: os.Getenv("PASSWORD"),
		Dbname:   os.Getenv("DB_NAME"),
	}

	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		logger.Error("Error load grps .env config. have not GRPC_PORT")
	}

	dbConn, err := repository.Connection(confPostgres)
	if err != nil {
		logger.Errorf("Error connecting to database: %v, config: %v", err, confPostgres)
		return
	}
	logger.Infof("Connected to PostgreSQL")

	err = runMigrations(logger, dbConn)
	if err != nil {
		logger.Errorf("can't run migrations: %v", err)
		return
	}

	userRepo := repository.NewUserRepository(dbConn, logger)

	userSvc := service.NewService(userRepo, logger)

	s := S{}
	rest := rest.NewSrver(s)
	go func() {
		rest.Run("8900")
	}()

	grpcTaskServer := grpc.NewUserHandler(userSvc, logger)
	logger.Info("created grpc server")
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Errorf("can't run grpc server: %v", err)
			}
		}()

		grpc.Run(grpcTaskServer, grpcPort)
	}()
	logger.Infof("grpc server is running on port: %s", grpcPort)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
}

func runMigrations(logger *logrus.Logger, dbConn *sql.DB) error {
	upMigration, err := strconv.ParseBool(os.Getenv("RUN_MIGRATION"))
	if err != nil {
		return err
	}

	if !upMigration {
		return nil
	}

	migrationDir := os.Getenv("MIGRATION_DIR")
	if migrationDir == "" {
		logger.Infof("no migration dir provided; skipping migrations")
		return nil
	}
	err = goose.Up(dbConn, migrationDir)
	if err != nil {
		return fmt.Errorf("migration dir: %v, %w", migrationDir, err)
	}
	logger.Infof("migrations are applied successfully")

	return nil
}

func NewLogger(level logrus.Level, formatter logrus.Formatter) *logrus.Logger {
	logger := logrus.New()
	logger.SetLevel(level)
	logger.SetFormatter(formatter)
	return logger
}
