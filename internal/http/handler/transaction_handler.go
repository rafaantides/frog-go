package handler

import (
	"errors"
	"frog-go/internal/config"
	"frog-go/internal/core/dto"
	appError "frog-go/internal/core/errors"
	"frog-go/internal/core/ports/inbound"
	"frog-go/internal/utils"
	"frog-go/internal/utils/pagination"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	service inbound.TransactionService
}

func NewTransactionHandler(service inbound.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

// CreateTransactionHandler godoc
// @Summary Cria uma nova transação
// @Description Cria uma nova transação com os dados fornecidos no corpo da requisição
// @Tags Transações
// @Accept json
// @Produce json
// @Param request body dto.TransactionRequest true "Dados da transação"
// @Success 201 {object} dto.TransactionResponse
// @Router /api/v1/transactions [post]
func (h *TransactionHandler) CreateTransactionHandler(c *gin.Context) {
	ctx := c.Request.Context()
	var req dto.TransactionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(appError.NewAppError(http.StatusBadRequest, err))
		return
	}

	input, err := req.ToDomain()
	if err != nil {
		c.Error(appError.NewAppError(http.StatusBadRequest, err))
		return
	}

	data, err := h.service.CreateTransaction(ctx, *input)
	if err != nil {
		c.Error(appError.NewAppError(http.StatusInternalServerError, err))
		return
	}

	c.JSON(http.StatusCreated, data)
}

// GetTransactionByIDHandler godoc
// @Summary Busca uma transação por ID
// @Description Retorna os dados de uma transação com base no ID fornecido
// @Tags Transações
// @Accept json
// @Produce json
// @Param id path string true "ID da transação"
// @Success 200 {object} dto.TransactionResponse
// @Router /api/v1/transactions/{id} [get]
func (h *TransactionHandler) GetTransactionByIDHandler(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := utils.ToUUID(c.Param("id"))
	if err != nil {
		c.Error(appError.NewAppError(http.StatusBadRequest, err))
		return
	}

	data, err := h.service.GetTransactionByID(ctx, id)
	if err != nil {
		if errors.Is(err, appError.ErrNotFound) {
			c.Error(appError.NewAppError(http.StatusNotFound, err))
			return
		}
		c.Error(appError.NewAppError(http.StatusInternalServerError, err))
		return
	}

	c.JSON(http.StatusOK, data)
}

// ListTransactionsHandler godoc
// @Summary Lista transações com filtros e paginação
// @Description Lista todas as transações aplicando filtros e paginação
// @Tags Transações
// @Accept json
// @Produce json
// @Param status query []string false "Filtrar por status"
// @Param kinds query []string false "Filtrar por kind"
// @Param category_id query []string false "Filtrar por categorias"
// @Param min_amount query number false "Valor mínimo"
// @Param max_amount query number false "Valor máximo"
// @Param start_date query string false "Data inicial"
// @Param end_date query string false "Data final"
// @Param page query int false "Número da página"
// @Param limit query int false "Limite por página"
// @Param order_by query string false "Campo de ordenação"
// @Param order query string false "Ordem (asc, desc)"
// @Success 200 {array} dto.TransactionResponse
// @Router /api/v1/transactions [get]
func (h *TransactionHandler) ListTransactionsHandler(c *gin.Context) {
	ctx := c.Request.Context()
	var flt dto.TransactionFilters
	if err := c.ShouldBindQuery(&flt); err != nil {
		c.Error(appError.NewAppError(http.StatusBadRequest, err))
		return
	}

	pgn, err := pagination.NewPagination(c)

	if err != nil {
		c.Error(appError.NewAppError(http.StatusBadRequest, err))
		return
	}

	validColumns := map[string]bool{
		"id":            true,
		"title":         true,
		"category_id":   true,
		"category":      true,
		"amount":        true,
		"purchase_date": true,
		"due_date":      true,
		"status":        true,
		"created_at":    true,
		"updated_at":    true,
	}

	if err := pgn.ValidateOrderBy("purchase_date", config.OrderAsc, validColumns); err != nil {
		c.Error(appError.NewAppError(http.StatusBadRequest, err))
		return
	}

	response, total, err := h.service.ListTransactions(ctx, flt, pgn)

	if err != nil {
		c.Error(appError.NewAppError(http.StatusInternalServerError, err))
		return
	}

	pgn.SetPaginationHeaders(c, total)

	c.JSON(http.StatusOK, response)
}

// UpdateTransactionHandler godoc
// @Summary Atualiza uma transação existente
// @Description Atualiza os dados de uma transação com base no ID fornecido e nos dados enviados no corpo da requisição
// @Tags Transações
// @Accept json
// @Produce json
// @Param id path string true "ID da transação"
// @Param request body dto.TransactionRequest true "Dados atualizados da transação"
// @Success 200 {object} dto.TransactionResponse
// @Router /api/v1/transactions/{id} [put]
func (h *TransactionHandler) UpdateTransactionHandler(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := utils.ToUUID(c.Param("id"))
	if err != nil {
		c.Error(appError.NewAppError(http.StatusBadRequest, err))
		return
	}

	var req dto.TransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(appError.NewAppError(http.StatusBadRequest, err))
		return
	}

	input, err := req.ToDomain()
	if err != nil {
		c.Error(appError.NewAppError(http.StatusBadRequest, err))
		return
	}

	data, err := h.service.UpdateTransaction(ctx, id, *input)
	if err != nil {
		c.Error(appError.NewAppError(http.StatusInternalServerError, err))
		return
	}

	c.JSON(http.StatusOK, data)
}

// DeleteTransactionHandler godoc
// @Summary Remove uma transação
// @Description Exclui uma transação com base no ID fornecido
// @Tags Transações
// @Accept json
// @Produce json
// @Param id path string true "ID da transação"
// @Success 204 "Sem conteúdo"
// @Router /api/v1/transactions/{id} [delete]
func (h *TransactionHandler) DeleteTransactionHandler(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := utils.ToUUID(c.Param("id"))
	if err != nil {
		c.Error(appError.NewAppError(http.StatusBadRequest, err))
		return
	}

	err = h.service.DeleteTransactionByID(ctx, id)
	if err != nil {
		c.Error(appError.NewAppError(http.StatusInternalServerError, err))
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// TransactionsSummaryHandler godoc
// @Summary Retorna resumo de transações
// @Description Gera um resumo estatístico das transações baseado nos filtros fornecidos
// @Tags Transações
// @Accept json
// @Produce json
// @Param start_date query string false "Data inicial (YYYY-MM-DD)"
// @Param end_date query string false "Data final (YYYY-MM-DD)"
// @Param kinds query []string false "Tipos de transação (income, expense)"
// @Success 200 {object} dto.SummaryByDate
// @Router /api/v1/transactions/summary [get]
func (h *TransactionHandler) TransactionsSummaryHandler(c *gin.Context) {
	ctx := c.Request.Context()
	var flt dto.ChartFilters
	if err := c.ShouldBindQuery(&flt); err != nil {
		c.Error(appError.NewAppError(http.StatusBadRequest, err))
		return
	}

	response, err := h.service.TransactionsSummary(ctx, flt)

	if err != nil {
		c.Error(appError.NewAppError(http.StatusInternalServerError, err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// TransactionsGeneralStatsHandler godoc
// @Summary Retorna estatísticas gerais das transações
// @Description Fornece dados estatísticos agregados das transações com base nos filtros aplicados
// @Tags Transações
// @Accept json
// @Produce json
// @Param start_date query string false "Data inicial (YYYY-MM-DD)"
// @Param end_date query string false "Data final (YYYY-MM-DD)"
// @Param kinds query []string false "Tipos de transação (income, expense)"
// @Success 200 {object} dto.TransactionStatsSummary
// @Router /api/v1/transactions/stats [get]
func (h *TransactionHandler) TransactionsGeneralStatsHandler(c *gin.Context) {
	ctx := c.Request.Context()
	var flt dto.ChartFilters
	if err := c.ShouldBindQuery(&flt); err != nil {
		c.Error(appError.NewAppError(http.StatusBadRequest, err))
		return
	}

	response, err := h.service.TransactionsGeneralStats(ctx, flt)

	if err != nil {
		c.Error(appError.NewAppError(http.StatusInternalServerError, err))
		return
	}

	c.JSON(http.StatusOK, response)
}
