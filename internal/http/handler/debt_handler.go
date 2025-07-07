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

type DebtHandler struct {
	service inbound.DebtService
}

func NewDebtHandler(service inbound.DebtService) *DebtHandler {
	return &DebtHandler{service: service}
}

// CreateDebtHandler godoc
// @Summary Cria uma nova dívida
// @Description Cria uma nova dívida com os dados fornecidos no corpo da requisição
// @Tags Dívidas
// @Accept json
// @Produce json
// @Param request body dto.DebtRequest true "Dados da dívida"
// @Success 201 {object} dto.DebtResponse
// @Router /v1/debts [post]
func (h *DebtHandler) CreateDebtHandler(c *gin.Context) {
	ctx := c.Request.Context()
	var req dto.DebtRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(appError.NewAppError(http.StatusBadRequest, err))
		return
	}

	input, err := req.ToDomain()
	if err != nil {
		c.Error(appError.NewAppError(http.StatusBadRequest, err))
		return
	}

	data, err := h.service.CreateDebt(ctx, *input)
	if err != nil {
		c.Error(appError.NewAppError(http.StatusInternalServerError, err))
		return
	}

	c.JSON(http.StatusCreated, data)
}

// GetDebtByIDHandler godoc
// @Summary Busca uma dívida por ID
// @Description Retorna os dados de uma dívida com base no ID fornecido
// @Tags Dívidas
// @Accept json
// @Produce json
// @Param id path string true "ID da dívida"
// @Success 200 {object} dto.DebtResponse
// @Router /v1/debts/{id} [get]
func (h *DebtHandler) GetDebtByIDHandler(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := utils.ToUUID(c.Param("id"))
	if err != nil {
		c.Error(appError.NewAppError(http.StatusBadRequest, err))
		return
	}

	data, err := h.service.GetDebtByID(ctx, id)
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

// ListDebtsHandler godoc
// @Summary Lista dívidas com filtros e paginação
// @Description Lista todas as dívidas aplicando filtros e paginação
// @Tags Dívidas
// @Accept json
// @Produce json
// @Param status query []string false "Filtrar por status"
// @Param category_id query []string false "Filtrar por categorias"
// @Param min_amount query number false "Valor mínimo"
// @Param max_amount query number false "Valor máximo"
// @Param start_date query string false "Data inicial"
// @Param end_date query string false "Data final"
// @Param page query int false "Número da página"
// @Param limit query int false "Limite por página"
// @Param order_by query string false "Campo de ordenação"
// @Param order query string false "Ordem (asc, desc)"
// @Success 200 {array} dto.DebtResponse
// @Router /v1/debts [get]
func (h *DebtHandler) ListDebtsHandler(c *gin.Context) {
	ctx := c.Request.Context()
	var flt dto.DebtFilters
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

	response, total, err := h.service.ListDebts(ctx, flt, pgn)

	if err != nil {
		c.Error(appError.NewAppError(http.StatusInternalServerError, err))
		return
	}

	pgn.SetPaginationHeaders(c, total)

	c.JSON(http.StatusOK, response)
}

func (h *DebtHandler) UpdateDebtHandler(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := utils.ToUUID(c.Param("id"))
	if err != nil {
		c.Error(appError.NewAppError(http.StatusBadRequest, err))
		return
	}

	var req dto.DebtRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(appError.NewAppError(http.StatusBadRequest, err))
		return
	}

	input, err := req.ToDomain()
	if err != nil {
		c.Error(appError.NewAppError(http.StatusBadRequest, err))
		return
	}

	data, err := h.service.UpdateDebt(ctx, id, *input)
	if err != nil {
		c.Error(appError.NewAppError(http.StatusInternalServerError, err))
		return
	}

	c.JSON(http.StatusOK, data)
}

func (h *DebtHandler) DeleteDebtHandler(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := utils.ToUUID(c.Param("id"))
	if err != nil {
		c.Error(appError.NewAppError(http.StatusBadRequest, err))
		return
	}

	err = h.service.DeleteDebtByID(ctx, id)
	if err != nil {
		c.Error(appError.NewAppError(http.StatusInternalServerError, err))
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *DebtHandler) DebtsSummaryHandler(c *gin.Context) {
	ctx := c.Request.Context()
	var flt dto.ChartFilters
	if err := c.ShouldBindQuery(&flt); err != nil {
		c.Error(appError.NewAppError(http.StatusBadRequest, err))
		return
	}

	response, err := h.service.DebtsSummary(ctx, flt)

	if err != nil {
		c.Error(appError.NewAppError(http.StatusInternalServerError, err))
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *DebtHandler) DebtsGeneralStatsHandler(c *gin.Context) {
	ctx := c.Request.Context()
	var flt dto.ChartFilters
	if err := c.ShouldBindQuery(&flt); err != nil {
		c.Error(appError.NewAppError(http.StatusBadRequest, err))
		return
	}

	response, err := h.service.DebtsGeneralStats(ctx, flt)

	if err != nil {
		c.Error(appError.NewAppError(http.StatusInternalServerError, err))
		return
	}

	c.JSON(http.StatusOK, response)
}
