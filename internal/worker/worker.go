package worker

import (
	"context"
	"frog-go/internal/core/ports/inbound"
	"frog-go/internal/core/ports/outbound/messagebus"
	"frog-go/internal/utils/logger"
	"sync"
	"time"
)

type Worker struct {
	ctx      context.Context
	consumer inbound.Consumer

	log  *logger.Logger
	mbus messagebus.MessageBus

	stopChan chan struct{} // Canal para sinalizar parada segura
	mu       sync.Mutex
}

func NewWorker(
	consumer inbound.Consumer,
	log *logger.Logger,
	mbus messagebus.MessageBus,
	stopChan chan struct{},
) *Worker {
	ctx := context.Background()
	return &Worker{
		ctx:      ctx,
		consumer: consumer,
		log:      log,
		mbus:     mbus,
		stopChan: stopChan, // Inicializa o canal de parada
		mu:       sync.Mutex{},
	}
}

func (w *Worker) Start(queue string, limit, timeoutSeconds int) {
	w.log.Start(
		"Processo iniciado... Fila: %s | Concorrência: %d mensagens | Timeout: %ds",
		queue, limit, timeoutSeconds,
	)

	var (
		wg        sync.WaitGroup
		semaphore = make(chan struct{}, limit)
	)

	for {
		select {
		case <-w.stopChan:
			w.log.Warn("Sinal de parada recebido. Encerrando worker...")
			wg.Wait()
			close(semaphore)
			w.log.Success("Worker finalizado com sucesso.")
			return
		default:
			messageHandler, err := w.mbus.Consume(queue)
			if err != nil {
				w.log.Error("Erro ao iniciar o consumo da fila %s: %v", queue, err)
				time.Sleep(5 * time.Second) // aguarda um tempo antes de tentar de novo
				continue
			}

			msgs := messageHandler.Messages()

			for {
				select {
				case <-w.stopChan:
					w.log.Warn("Sinal de parada recebido. Encerrando consumo atual...")
					messageHandler.Close()
					wg.Wait()
					close(semaphore)
					w.log.Success("Worker finalizado com sucesso.")
					return

				case msg, ok := <-msgs:
					if !ok {
						w.log.Warn("Canal de mensagens fechado. Tentando reconectar...")
						messageHandler.Close()
						time.Sleep(2 * time.Second)
						break
					}

					select {
					case semaphore <- struct{}{}:
						wg.Add(1)
						go func(msg messagebus.Message) {
							defer wg.Done()
							defer func() { <-semaphore }()

							if err := w.processMessage(timeoutSeconds, msg.Body()); err != nil {
								w.log.Error("Erro ao processar mensagem: %v", err)
							}

							if err := msg.Ack(); err != nil {
								w.log.Error("Falha ao confirmar a mensagem: %v", err)
							}
						}(msg)
					case <-w.stopChan:
						messageHandler.Close()
						wg.Wait()
						close(semaphore)
						w.log.Success("Worker finalizado com sucesso.")
						return
					}
				}
			}
		}
	}
}

func (w *Worker) processMessage(timeoutSeconds int, messageBody []byte) error {
	if err := w.consumer.ProcessMessage(timeoutSeconds, messageBody); err != nil {
		w.log.Error("ctx: %s | %v", w.ctx, err)
		return err
	}
	return nil
}

func (w *Worker) Stop() {
	// Método para parar o worker com segurança
	w.mu.Lock()
	defer w.mu.Unlock()

	select {
	case <-w.stopChan:
		// Se já foi fechado, não faz nada
	default:
		close(w.stopChan)
	}
}
