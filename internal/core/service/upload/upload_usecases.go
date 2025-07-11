package upload

import (
	"encoding/json"
	"fmt"
	"frog-go/internal/config"
	"frog-go/internal/core/dto"
	"frog-go/internal/core/errors"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

func (s *uploadService) processTransactions(model, action, filename string, rows [][]string, idx map[string]int) error {
	jobID := uuid.New().String()

	baseName := strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))
	parts := strings.Split(baseName, "_")
	if len(parts) < 2 {
		return fmt.Errorf("invalid filename format: %s", filename)
	}

	dueDate := parts[1]

	for _, row := range rows[1:] {

		transaction, err := buildTransactionRequest(model, dueDate, row, idx)
		if err != nil {
			return fmt.Errorf("failed to build request: %w", err)
		}

		msg := dto.ImportTxnMessage{
			JobID:    jobID,
			Filename: filename,
			Action:   action,
			Data: struct {
				Transaction dto.TransactionRequest `json:"transaction"`
			}{
				Transaction: *transaction,
			},
		}

		messageBytes, err := json.Marshal(msg)
		if err != nil {
			return fmt.Errorf("failed to serialize message: %w", err)
		}

		if err := s.mb.SendMessage(config.ResourceTransactions, messageBytes); err != nil {
			return err
		}
	}

	return nil
}

func buildTransactionRequest(model, dueDate string, row []string, idx map[string]int) (*dto.TransactionRequest, error) {
	switch model {
	case config.ModelNubank:
		return nubankToRequest(dueDate, row, idx)
	default:
		return nil, fmt.Errorf("unknown model: %s", model)
	}
}

func nubankToRequest(dueDate string, row []string, idx map[string]int) (*dto.TransactionRequest, error) {

	amount, err := strconv.ParseFloat(getValue(row, idx, "amount"), 64)
	if err != nil {
		return nil, errors.InvalidParam("amount", err)
	}

	return &dto.TransactionRequest{
		DueDate:      &dueDate,
		PurchaseDate: getValue(row, idx, "date"),
		Title:        getValue(row, idx, "title"),
		Amount:       amount,
	}, nil
}

func getValue(row []string, idx map[string]int, key string) string {
	if i, ok := idx[key]; ok && i < len(row) {
		return row[i]
	}
	return ""
}
