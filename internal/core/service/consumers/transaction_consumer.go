package consumers

import (
	"context"
	"encoding/json"
	"fmt"
	"frog-go/internal/config"
	"frog-go/internal/core/dto"
	"frog-go/internal/core/ports/inbound"
	"frog-go/internal/utils"
	"frog-go/internal/utils/logger"
	"strings"
	"time"
)

type Message struct {
	ID   string
	Body []byte
}

type TransactionConsumer struct {
	service inbound.TransactionService
	cfg     *config.ConfigConsumer
	log     *logger.Logger
}

func NewTransactionConsumer(
	service inbound.TransactionService,
	cfg *config.ConfigConsumer,
) *TransactionConsumer {
	return &TransactionConsumer{
		service: service,
		cfg:     cfg,
		log:     logger.NewLogger("TransactionConsumer"),
	}
}

func (c *TransactionConsumer) ProcessMessage(
	timeoutSeconds int,
	messageBody []byte,
) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSeconds)*time.Second)
	defer cancel()

	var msg dto.ImportTxnMessage
	if err := json.Unmarshal(messageBody, &msg); err != nil {
		return fmt.Errorf("failed to unmarshal ImportTxnMessage: %w", err)
	}

	userID, err := utils.ToUUID(msg.UserID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	c.log.Info("Processing message: %+v", msg)
	switch msg.Action {
	case config.ActionCreate:
		input, err := msg.Data.Transaction.ToDomain()
		if err != nil {
			return fmt.Errorf("failed to parse debt: %w", err)
		}

		if c.shouldSkipTitle(input.Title) {
			c.log.Info("Skipping title: %s", input.Title)
			return nil
		}

		if _, err := c.service.CreateTransaction(ctx, userID, *input); err != nil {
			return fmt.Errorf("failed to create transaction: %w", err)
		}

	default:
		return fmt.Errorf("invalid action: %s", msg.Action)
	}

	return nil
}

func (c *TransactionConsumer) shouldSkipTitle(title string) bool {
	for _, skip := range c.cfg.SkipTitles {
		if strings.EqualFold(skip, title) {
			return true
		}
	}
	return false
}
