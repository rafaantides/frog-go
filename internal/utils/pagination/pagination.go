package pagination

import (
	"fmt"
	"frog-go/internal/config"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Pagination struct {
	Page           int    `json:"page"`
	PageSize       int    `json:"page_size"`
	OrderBy        string `json:"order_by"`
	OrderDirection string `json:"order_direction"`
	Search         string `json:"search"`
}

func NewPagination(c *gin.Context) (*Pagination, error) {
	page, err := strconv.Atoi(c.DefaultQuery("page", config.DefaultPage))
	if err != nil || page < 1 {
		return nil, fmt.Errorf("%s: %s", "page", "invalid value")
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", config.DefaultPageSize))
	if err != nil || pageSize < 1 || pageSize > config.MaxPageSize {
		return nil, fmt.Errorf("%s: %s", "page_size", "invalid value")
	}

	orderBy := c.Query("order_by")
	orderDirection := c.Query("order_direction")
	search := c.Query("search")

	return &Pagination{
		Page:           page,
		PageSize:       pageSize,
		OrderBy:        orderBy,
		OrderDirection: orderDirection,
		Search:         search,
	}, nil
}

func (p *Pagination) Offset() int {
	return (p.Page - 1) * p.PageSize
}

func (p *Pagination) ValidateOrderBy(defaultOrder string, defaultDirection string, validColumns map[string]bool) error {
	if p.OrderBy == "" {
		p.OrderBy = defaultOrder
	}

	if p.OrderDirection == "" {
		p.OrderDirection = defaultDirection
	}

	if !validColumns[p.OrderBy] {
		return fmt.Errorf("%s: %s", "order_by", "invalid value")
	}

	if p.OrderDirection != config.OrderAsc && p.OrderDirection != config.OrderDesc {
		return fmt.Errorf("%s: %s", "order_direction", "invalid value")
	}
	return nil
}

func (p *Pagination) SetPaginationHeaders(c *gin.Context, total int) {
	totalPages := (total + p.PageSize - 1) / p.PageSize

	c.Header("X-Page", strconv.Itoa(p.Page))
	c.Header("X-Page-Size", strconv.Itoa(p.PageSize))
	c.Header("X-Total-Count", strconv.Itoa(total))
	c.Header("X-Total-Pages", strconv.Itoa(totalPages))
}
