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

type InvoiceHandler struct {
	service inbound.InvoiceService
}

func NewInvoiceHandler(service inbound.InvoiceService) *InvoiceHandler {
	return &InvoiceHandler{service: service}
}

// CreateInvoiceHandler godoc
// @Summary Cria uma nova fatura
// @Description Cria uma nova fatura com os dados fornecidos no corpo da requisição
// @Tags Faturas
// @Accept json
// @Produce json
// @Param request body dto.InvoiceRequest true "Dados da fatura"
// @Success 201 {object} dto.InvoiceResponse
// @Security BearerAuth
// @Router /api/v1/invoices [post]
func (h *InvoiceHandler) CreateInvoiceHandler(c *gin.Context) {
	ctx := c.Request.Context()
	var req dto.InvoiceRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(appError.NewAppError(http.StatusBadRequest, err))
		return
	}

	input, err := req.ToDomain()
	if err != nil {
		c.Error(appError.NewAppError(http.StatusBadRequest, err))
		return
	}

	data, err := h.service.CreateInvoice(ctx, *input)
	if err != nil {
		c.Error(appError.NewAppError(http.StatusInternalServerError, err))
		return
	}

	c.JSON(http.StatusCreated, data)
}

// GetInvoiceByIDHandler godoc
// @Summary Busca uma fatura por ID
// @Description Retorna os dados de uma fatura com base no ID fornecido
// @Tags Faturas
// @Accept json
// @Produce json
// @Param id path string true "ID da fatura"
// @Success 200 {object} dto.InvoiceResponse
// @Security BearerAuth
// @Router /api/v1/invoices/{id} [get]
func (h *InvoiceHandler) GetInvoiceByIDHandler(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := utils.ToUUID(c.Param("id"))
	if err != nil {
		c.Error(appError.NewAppError(http.StatusBadRequest, err))
		return
	}

	data, err := h.service.GetInvoiceByID(ctx, id)
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

// ListInvoicesHandler godoc
// @Summary Lista faturas com filtros e paginação
// @Description Lista todas as faturas aplicando filtros e paginação
// @Tags Faturas
// @Accept json
// @Produce json
// @Param page query int false "Número da página"
// @Param limit query int false "Limite por página"
// @Param order_by query string false "Campo de ordenação (ex: due_date)"
// @Param order query string false "Ordem (asc, desc)"
// @Param title query string false "Filtro por título"
// @Param status query string false "Filtro por status"
// @Param min_amount query number false "Valor mínimo"
// @Param max_amount query number false "Valor máximo"
// @Param due_date_start query string false "Data de vencimento inicial"
// @Param due_date_end query string false "Data de vencimento final"
// @Success 200 {array} dto.InvoiceResponse
// @Security BearerAuth
// @Router /api/v1/invoices [get]
func (h *InvoiceHandler) ListInvoicesHandler(c *gin.Context) {
	ctx := c.Request.Context()
	var flt dto.InvoiceFilters
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
		"id":         true,
		"title":      true,
		"amount":     true,
		"due_date":   true,
		"status_id":  true,
		"status":     true,
		"created_at": true,
		"updated_at": true,
	}

	if err := pgn.ValidateOrderBy("due_date", config.OrderAsc, validColumns); err != nil {
		c.Error(appError.NewAppError(http.StatusBadRequest, err))
		return
	}

	response, total, err := h.service.ListInvoices(ctx, flt, pgn)

	if err != nil {
		c.Error(appError.NewAppError(http.StatusInternalServerError, err))
		return
	}

	pgn.SetPaginationHeaders(c, total)

	c.JSON(http.StatusOK, response)
}

// UpdateInvoiceHandler godoc
// @Summary Atualiza uma fatura existente
// @Description Atualiza os dados de uma fatura com base no ID fornecido e nos dados enviados no corpo da requisição
// @Tags Faturas
// @Accept json
// @Produce json
// @Param id path string true "ID da fatura"
// @Param request body dto.InvoiceRequest true "Dados atualizados da fatura"
// @Success 200 {object} dto.InvoiceResponse
// @Security BearerAuth
// @Router /api/v1/invoices/{id} [put]
func (h *InvoiceHandler) UpdateInvoiceHandler(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := utils.ToUUID(c.Param("id"))
	if err != nil {
		c.Error(appError.NewAppError(http.StatusBadRequest, err))
		return
	}

	var req dto.InvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(appError.NewAppError(http.StatusBadRequest, err))
		return
	}
	input, err := req.ToDomain()
	if err != nil {
		c.Error(appError.NewAppError(http.StatusBadRequest, err))
		return
	}

	data, err := h.service.UpdateInvoice(ctx, id, *input)
	if err != nil {
		c.Error(appError.NewAppError(http.StatusInternalServerError, err))
		return
	}

	c.JSON(http.StatusOK, data)
}

// DeleteInvoiceHandler godoc
// @Summary Remove uma fatura
// @Description Exclui uma fatura com base no ID fornecido
// @Tags Faturas
// @Accept json
// @Produce json
// @Param id path string true "ID da fatura"
// @Success 204 "Sem conteúdo"
// @Security BearerAuth
// @Router /api/v1/invoices/{id} [delete]
func (h *InvoiceHandler) DeleteInvoiceHandler(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := utils.ToUUID(c.Param("id"))
	if err != nil {
		c.Error(appError.NewAppError(http.StatusBadRequest, err))
		return
	}

	err = h.service.DeleteInvoiceByID(ctx, id)
	if err != nil {
		c.Error(appError.NewAppError(http.StatusInternalServerError, err))
		return
	}

	c.Status(http.StatusNoContent)
}

// ListInvoiceDebtsHandler godoc
// @Summary Lista débitos associados a uma fatura
// @Description Lista todos os débitos (transações) vinculados a uma fatura específica, com filtros e paginação
// @Tags Faturas
// @Accept json
// @Produce json
// @Param id path string true "ID da fatura"
// @Param page query int false "Número da página"
// @Param limit query int false "Limite por página"
// @Param order_by query string false "Campo de ordenação (ex: record_date)"
// @Param order query string false "Ordem (asc, desc)"
// @Param title query string false "Filtro por título da transação"
// @Param status query string false "Filtro por status da transação"
// @Param record_type query string false "Tipo da transação (income, expense)"
// @Param category_id query string false "ID da categoria"
// @Param record_date_start query string false "Data inicial de registro"
// @Param record_date_end query string false "Data final de registro"
// @Success 200 {array} dto.TransactionResponse
// @Security BearerAuth
// @Router /api/v1/invoices/{id}/debts [get]
func (h *InvoiceHandler) ListInvoiceDebtsHandler(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := utils.ToUUID(c.Param("id"))
	if err != nil {
		c.Error(appError.NewAppError(http.StatusBadRequest, err))
		return
	}
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
		"id":          true,
		"title":       true,
		"category_id": true,
		"amount":      true,
		"record_date": true,
		"status":      true,
		"record_type": true,
		"created_at":  true,
		"updated_at":  true,
	}

	if err := pgn.ValidateOrderBy("record_date", config.OrderAsc, validColumns); err != nil {
		c.Error(appError.NewAppError(http.StatusBadRequest, err))
		return
	}

	response, total, err := h.service.ListInvoiceDebts(ctx, id, flt, pgn)

	if err != nil {
		c.Error(appError.NewAppError(http.StatusInternalServerError, err))
		return
	}

	pgn.SetPaginationHeaders(c, total)

	c.JSON(http.StatusOK, response)
}
