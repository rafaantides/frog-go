package main

import (
	"flag"
	"fmt"
	"os"

	"frog-go/internal/config/bootstrap"
	"frog-go/internal/core/service/consumers"
	"frog-go/internal/utils/logger"
	"frog-go/internal/worker"
)

var (
	limit   int
	timeout int
	debug   bool
	queue   string
	port    string
	envPath string
)

func init() {
	flag.IntVar(&limit, "limit", 5, "Número máximo de mensagens processadas simultaneamente (concorrência)")
	flag.IntVar(&timeout, "timeout", 30, "Timeout em segundos para processamento de cada mensagem")
	flag.StringVar(&queue, "queue", "development", "Nome da fila a ser processada")

	flag.StringVar(&port, "port", "8080", "Porta para executar o servidor API")
	flag.BoolVar(&debug, "debug", false, "Habilita o modo debug")
	flag.StringVar(&envPath, "env", ".env", "Caminho para o arquivo .env")
}

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		fmt.Println("Error: você deve fornecer exatamente um argumento indicando o tipo de consumidor")
		os.Exit(1)
	}

	startConsumer(args[0])
}

func startConsumer(resource string) {
	log := logger.NewLogger("Worker")

	boot, err := bootstrap.InitWorker(envPath)
	if err != nil {
		log.Fatal("%v", err)
	}
	defer boot.Repo.Close()
	defer boot.Mbus.Close()
	defer boot.Cache.Close()

	factory, ok := consumers.Registry[resource]
	if !ok {
		log.Fatal("Consumer inválido: %s", resource)
	}

	consumer := factory(boot)
	stopChan := make(chan struct{})

	log.Info("Iniciando worker com fila: %s | limite: %d | timeout: %ds", queue, limit, timeout)

	w := worker.NewWorker(consumer, log, boot.Mbus, boot.Noti, stopChan)
	w.Start(queue, limit, timeout)
}
