package handler

import (
	"errors"
	"fmt"
	"frog-go/internal/config"
	"frog-go/internal/core/dto"
	appError "frog-go/internal/core/errors"
	"frog-go/internal/core/ports/inbound"
	"frog-go/internal/utils"
	"frog-go/internal/utils/pagination"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	service inbound.CategoryService
}

func NewCategoryHandler(service inbound.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

// CreateCategoryHandler godoc
// @Summary Cria uma nova categoria
// @Description Cria uma nova categoria com os dados fornecidos no corpo da requisição
// @Tags Categorias
// @Accept json
// @Produce json
// @Param request body dto.CategoryRequest true "Dados da categoria"
// @Success 201 {object} dto.CategoryResponse
// @Router /api/v1/categories [post]
func (h *CategoryHandler) CreateCategoryHandler(c *gin.Context) {
	ctx := c.Request.Context()
	var req dto.CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(appError.NewAppError(http.StatusBadRequest, err))
		return
	}

	input, err := req.ToDomain()
	if err != nil {
		c.Error(appError.NewAppError(http.StatusBadRequest, err))
		return
	}

	data, err := h.service.CreateCategory(ctx, *input)
	if err != nil {
		c.Error(appError.NewAppError(http.StatusInternalServerError, err))
		return
	}

	c.JSON(http.StatusCreated, data)
}

// GetCategoryByIDHandler godoc
// @Summary Busca uma categoria por ID
// @Description Retorna os dados de uma categoria com base no ID fornecido
// @Tags Categorias
// @Accept json
// @Produce json
// @Param id path string true "ID da categoria"
// @Success 200 {object} dto.CategoryResponse
// @Router /api/v1/categories/{id} [get]
func (h *CategoryHandler) GetCategoryByIDHandler(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := utils.ToUUID(c.Param("id"))
	if err != nil {
		c.Error(appError.NewAppError(http.StatusBadRequest, err))
		return
	}

	data, err := h.service.GetCategoryByID(ctx, id)
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

// ListCategorysHandler godoc
// @Summary Lista categorias com filtros e paginação
// @Description Lista todas as categorias aplicando filtros e paginação
// @Tags Categorias
// @Accept json
// @Produce json
// @Param page query int false "Número da página"
// @Param limit query int false "Limite por página"
// @Param order_by query string false "Campo de ordenação (ex: name)"
// @Param order query string false "Ordem (asc, desc)"
// @Success 200 {array} dto.CategoryResponse
// @Router /api/v1/categories [get]
func (h *CategoryHandler) ListCategorysHandler(c *gin.Context) {
	ctx := c.Request.Context()
	pgn, err := pagination.NewPagination(c)

	if err != nil {
		c.Error(appError.NewAppError(http.StatusBadRequest, err))
		return
	}

	validColumns := map[string]bool{
		"id":          true,
		"name":        true,
		"description": true,
	}

	if err := pgn.ValidateOrderBy("name", config.OrderAsc, validColumns); err != nil {
		c.Error(appError.NewAppError(http.StatusBadRequest, err))
		return
	}

	fmt.Printf("%v", pgn)

	response, total, err := h.service.ListCategories(ctx, pgn)

	if err != nil {
		c.Error(appError.NewAppError(http.StatusInternalServerError, err))
		return
	}

	pgn.SetPaginationHeaders(c, total)

	c.JSON(http.StatusOK, response)
}

// UpdateCategoryHandler godoc
// @Summary Atualiza uma categoria existente
// @Description Atualiza os dados de uma categoria com base no ID fornecido e nos dados enviados no corpo da requisição
// @Tags Categorias
// @Accept json
// @Produce json
// @Param id path string true "ID da categoria"
// @Param request body dto.CategoryRequest true "Dados atualizados da categoria"
// @Success 200 {object} dto.CategoryResponse
// @Router /api/v1/categories/{id} [put]
func (h *CategoryHandler) UpdateCategoryHandler(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := utils.ToUUID(c.Param("id"))
	if err != nil {
		c.Error(appError.NewAppError(http.StatusBadRequest, err))
		return
	}

	var req dto.CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(appError.NewAppError(http.StatusBadRequest, err))
		return
	}

	input, err := req.ToDomain()
	if err != nil {
		c.Error(appError.NewAppError(http.StatusBadRequest, err))
		return
	}

	data, err := h.service.UpdateCategory(ctx, id, *input)
	if err != nil {
		c.Error(appError.NewAppError(http.StatusInternalServerError, err))
		return
	}

	c.JSON(http.StatusOK, data)
}

// DeleteCategoryHandler godoc
// @Summary Remove uma categoria
// @Description Exclui uma categoria com base no ID fornecido
// @Tags Categorias
// @Accept json
// @Produce json
// @Param id path string true "ID da categoria"
// @Success 204 "Sem conteúdo"
// @Router /api/v1/categories/{id} [delete]
func (h *CategoryHandler) DeleteCategoryHandler(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := utils.ToUUID(c.Param("id"))
	if err != nil {
		c.Error(appError.NewAppError(http.StatusBadRequest, err))
		return
	}

	err = h.service.DeleteCategoryByID(ctx, id)
	if err != nil {
		c.Error(appError.NewAppError(http.StatusInternalServerError, err))
		return
	}

	c.Status(http.StatusNoContent)
}
