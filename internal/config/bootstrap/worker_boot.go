package bootstrap

import (
	"fmt"
	"frog-go/internal/adapters/messagebus/rabbitmq"
	"frog-go/internal/adapters/notifier/discord"
	"frog-go/internal/adapters/cachestorage/redis"
	"frog-go/internal/adapters/repository/postgresql"
	"frog-go/internal/config"
	"frog-go/internal/core/ports/outbound/cachestorage"
	"frog-go/internal/core/ports/outbound/messagebus"
	"frog-go/internal/core/ports/outbound/notifier"
	"frog-go/internal/core/ports/outbound/repository"
)

type WorkerDeps struct {
	Repo  repository.Repository
	Mbus  messagebus.MessageBus
	Cache cachestorage.CacheStorage
	Noti  notifier.Notifier
	Cfg   *config.ConfigConsumer
}

func InitWorker(envPath string) (*WorkerDeps, error) {

	cfg, err := config.LoadConfig(envPath)
	if err != nil {
		return nil, err
	}

	wcfg, err := config.LoadWorkerConfig(envPath)
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

	cch, err := redis.NewRedis(
		wcfg.CachePass,
		wcfg.CacheHost,
		wcfg.CachePort,
	)
	if err != nil {
		repo.Close()
		mbus.Close()
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	noti := discord.NewDiscord(cch, wcfg.NotifierWebhookURL)

	return &WorkerDeps{
		Repo:  repo,
		Mbus:  mbus,
		Cache: cch,
		Noti:  noti,
		Cfg:   config.LoadConsumerConfig(envPath),
	}, nil
}
