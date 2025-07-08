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

	kinds := c.QueryArray("kinds")

	response, total, err := h.service.ListCategories(ctx, kinds, pgn)

	if err != nil {
		c.Error(appError.NewAppError(http.StatusInternalServerError, err))
		return
	}

	pgn.SetPaginationHeaders(c, total)

	c.JSON(http.StatusOK, response)
}

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

	c.JSON(http.StatusNoContent, nil)
}
