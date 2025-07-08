package consumers

import (
	"frog-go/internal/config"
	"frog-go/internal/config/bootstrap"
	"frog-go/internal/core/ports/inbound"
	"frog-go/internal/core/service"
)

type ConsumerFactory func(*bootstrap.WorkerDeps) inbound.Consumer

var Registry = map[string]ConsumerFactory{
	config.ResourceTransactions: func(b *bootstrap.WorkerDeps) inbound.Consumer {
		txnService := service.NewTransactionService(b.Repo)
		consumer := NewTransactionConsumer(txnService, b.Cfg)
		return consumer
	},
}
