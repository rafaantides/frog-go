package bootstrap

import (
	"fmt"
	"frog-go/internal/adapters/messagebus/rabbitmq"
	"frog-go/internal/adapters/repository/postgresql"
	"frog-go/internal/config"
	"frog-go/internal/core/ports/outbound/messagebus"
	"frog-go/internal/core/ports/outbound/repository"
)

type WorkerDeps struct {
	Repo  repository.Repository
	Mbus  messagebus.MessageBus
	Cfg   *config.ConfigConsumer
}

func InitWorker(envPath string) (*WorkerDeps, error) {

	cfg, err := config.LoadConfig(envPath)
	if err != nil {
		return nil, err
	}

	repo, err := postgresql.NewPostgreSQL(
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
		cfg.SeedPath,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %v", err)
	}

	mbus, err := rabbitmq.NewRabbitMQ(
		cfg.MessageBusUser,
		cfg.MessageBusPass,
		cfg.MessageBusHost,
		cfg.MessageBusPort,
	)
	if err != nil {
		repo.Close()
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	return &WorkerDeps{
		Repo:  repo,
		Mbus:  mbus,
		Cfg:   config.LoadConsumerConfig(envPath),
	}, nil
}
