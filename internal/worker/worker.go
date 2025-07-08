package worker

import (
	"context"
	"fmt"
	"frog-go/internal/core/ports/inbound"
	"frog-go/internal/core/ports/outbound/messagebus"
	"frog-go/internal/core/ports/outbound/notifier"
	"frog-go/internal/utils/logger"
	"sync"
	"time"
)

type Worker struct {
	ctx      context.Context
	consumer inbound.Consumer

	log  *logger.Logger
	mbus messagebus.MessageBus
	noti notifier.Notifier

	stopChan chan struct{} // Canal para sinalizar parada segura
	mu       sync.Mutex
}

func NewWorker(
	consumer inbound.Consumer,
	log *logger.Logger,
	mbus messagebus.MessageBus,
	noti notifier.Notifier,
	stopChan chan struct{},
) *Worker {
	ctx := context.Background()
	return &Worker{
		ctx:      ctx,
		consumer: consumer,
		log:      log,
		mbus:     mbus,
		noti:     noti,
		stopChan: stopChan, // Inicializa o canal de parada
		mu:       sync.Mutex{},
	}
}

func (w *Worker) Start(queue string, limit, timeoutSeconds int) {
	w.log.Start(
		"Processo iniciado... Fila: %s | Concorrência: %d mensagens | Timeout: %ds",
		queue, limit, timeoutSeconds,
	)

	processInfo := fmt.Sprintf(
		"\n\nFila: %s\nConcorrência: %d mensagens\nTimeout: %ds",
		queue,
		limit,
		timeoutSeconds,
	)

	w.noti.SendMessage(w.ctx, fmt.Sprintf("**Processo iniciado 🚀**%s", processInfo))

	var (
		wg          sync.WaitGroup
		semaphore   = make(chan struct{}, limit)
		idleTimer   = time.NewTimer(5 * time.Minute)
		hadErrors   bool
		mu          sync.Mutex
		idleTimeout = 5 * time.Minute
	)

	messageHandler, err := w.mbus.Consume(queue)
	if err != nil {
		w.log.Fatal("Erro ao iniciar o consumo da fila %s: %v", queue, err)
	}
	defer messageHandler.Close()

	msgs := messageHandler.Messages()

	go func() {
		for {
			select {
			case <-idleTimer.C:
				if len(semaphore) == 0 {
					w.log.Warn("Nenhuma mensagem recebida e nenhum processamento em andamento nos últimos 5 minutos. Encerrando worker.")
					if err := w.mbus.DeleteQueue(queue); err != nil {
						w.log.Error("Falha ao deletar a fila %s: %v", queue, err)
					} else {
						w.log.Info("Fila %s deletada com sucesso.", queue)
					}
					w.Stop()
					return
				}
				idleTimer.Reset(idleTimeout)
			case <-w.stopChan:
				return
			}
		}
	}()

consumeLoop:
	for {
		select {
		case <-w.stopChan:
			w.log.Warn("Sinal de parada recebido. Encerrando worker...")
			break consumeLoop

		case msg, ok := <-msgs:
			if !ok {
				w.log.Warn("Canal de mensagens fechado. Nenhuma mensagem restante para processar.")
				break consumeLoop
			}

			select {
			case semaphore <- struct{}{}:
				if !idleTimer.Stop() {
					select {
					case <-idleTimer.C:
					default:
					}
				}
				idleTimer.Reset(idleTimeout)

				wg.Add(1)

				go func(msg messagebus.Message) {
					defer wg.Done()
					defer func() { <-semaphore }()

					if err := w.processMessage(timeoutSeconds, msg.Body()); err != nil {
						mu.Lock()
						hadErrors = true
						mu.Unlock()
						w.log.Error("Erro ao processar mensagem: %v", err)
					}

					if err := msg.Ack(); err != nil {
						w.log.Error("Falha ao confirmar a mensagem: %v", err)
					}
				}(msg)
			case <-w.stopChan:
				break consumeLoop
			}
		}
	}

	idleTimer.Stop()
	wg.Wait()
	close(semaphore)

	if hadErrors {
		w.noti.SendMessage(w.ctx, fmt.Sprintf("**Processo finalizado com erros de processamento ❌**%s", processInfo))
	} else {
		w.noti.SendMessage(w.ctx, fmt.Sprintf("**Processo finalizado com sucesso ✅**%s", processInfo))
	}

	w.log.Success("Worker finalizado com sucesso.")
}

func (w *Worker) processMessage(timeoutSeconds int, messageBody []byte) error {
	if err := w.consumer.ProcessMessage(timeoutSeconds, messageBody); err != nil {
		w.log.Error("ctx: %s | %v", w.ctx, err)
		return err
	}
	return nil
}

// Método para parar o worker com segurança
func (w *Worker) Stop() {
	w.mu.Lock()
	defer w.mu.Unlock()

	select {
	case <-w.stopChan:
		// Se já foi fechado, não faz nada
	default:
		close(w.stopChan)
	}
}
