package handler

import (
	appError "frog-go/internal/core/errors"
	"frog-go/internal/core/ports/inbound"
	"frog-go/internal/utils"
	"frog-go/internal/utils/utilsctx"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	service inbound.UploadService
}

// ProcessFileHandler trata o upload de um arquivo para importação de dados.
//
// @Summary Processar arquivo
// @Description Recebe um arquivo e os parâmetros necessários para processamento assíncrono.
// @Tags Upload
// @Accept multipart/form-data
// @Produce json
// @Param file formData  file   true   "Arquivo para importação (ex: .csv, .xlsx)"
// @Param invoice_id formData  string true   "ID da fatura (opcional)"
// @Param model formData  string true   "Nome do modelo alvo (ex: Nubank)"
// @Param action formData  string true   "Ação desejada (ex: create)"
// @Success 202 {object} map[string]string "Arquivo recebido, processamento em andamento"
// @Failure 400 {object} map[string]string "Erro nos parâmetros ou no upload"
// @Security BearerAuth
// @Router /api/v1/upload [post]
func NewUploadHandler(service inbound.UploadService) *UploadHandler {
	return &UploadHandler{service: service}
}

func (h *UploadHandler) ProcessFileHandler(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := utilsctx.GetUserID(ctx)
	if err != nil {
		c.Error(appError.NewAppError(http.StatusBadRequest, err))
		return
	}

	action := c.PostForm("action")
	model := c.PostForm("model")

	invoiceID, err := utils.ToNillableUUID(c.PostForm("invoice_id"))
	if err != nil {
		c.Error(
			appError.NewAppError(
				http.StatusBadRequest,
				appError.InvalidParam("invoice_id", err),
			),
		)
		return
	}

	if action == "" || model == "" {
		c.Error(
			appError.NewAppError(
				http.StatusBadRequest,
				appError.EmptyField("model and action"),
			),
		)
		return
	}

	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		c.Error(
			appError.NewAppError(
				http.StatusBadRequest,
				appError.FailedToFind("flile", err),
			),
		)
		return
	}
	defer file.Close()

	if err := h.service.ImportFile(userID, model, action, invoiceID, file, fileHeader); err != nil {
		c.Error(appError.NewAppError(http.StatusInternalServerError, err))
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "Arquivo recebido, processamento em andamento"})
}
