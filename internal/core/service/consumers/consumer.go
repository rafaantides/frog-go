package consumers

import (
	"context"
	"encoding/json"
	"fmt"
	"frog-go/internal/config"
	"frog-go/internal/core/dto"
	"frog-go/internal/core/ports/outbound/messagebus"
	"frog-go/internal/utils/logger"
	"net/http"
	"time"
)

type Message struct {
	ID   string
	Body []byte
}

type Consumer struct {
	mbus       messagebus.MessageBus
	log        *logger.Logger
	httpClient *http.Client
}

func NewSpiderConsumer(
	mbus messagebus.MessageBus,
) *Consumer {

	httpClient := &http.Client{}

	return &Consumer{
		mbus:       mbus,
		log:        logger.NewLogger("Consumer"),
		httpClient: httpClient,
	}
}

func (c *Consumer) ProcessMessage(
	queue string,
	timeoutSeconds int,
	maxAttempts int,
	messageBody []byte,
) error {
	// ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSeconds)*time.Second)
	_, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSeconds)*time.Second)
	defer cancel()

	var msg dto.MessageData
	if err := json.Unmarshal(messageBody, &msg); err != nil {
		return fmt.Errorf("failed to unmarshal MessageData: %w", err)
	}
	c.log.Info("Processing message: %+v", msg)
	// start := time.Now()

	// processed := c.processRequest(msg)
	c.processRequest(msg)

	// processingTime := time.Since(start)
	// processingTimeformatted := utils.FormatDuration(processingTime)

	return nil
}

func (c *Consumer) processRequest(msg dto.MessageData) dto.ProcessResult {
	return dto.ProcessResult{
		Retry:     true,
		Processed: false,
		Status:    config.StatusSuccess,
	}

}
