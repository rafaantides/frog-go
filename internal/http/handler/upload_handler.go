package handler

import (
	"frog-go/internal/core/ports/inbound"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	service inbound.UploadService
}

// ProcessFileHandler trata o upload de um arquivo para importação de dados.
//
// @Summary      Processar arquivo
// @Description  Recebe um arquivo e os parâmetros necessários para processamento assíncrono.
// @Tags         upload
// @Accept       multipart/form-data
// @Produce      json
// @Param        file     formData  file   true   "Arquivo para importação (ex: .csv, .xlsx)"
// @Param        resource formData  string true   "Nome do recurso (ex: transactions)"
// @Param        model    formData  string true   "Nome do modelo alvo (ex: Nubank)"
// @Param        action   formData  string true   "Ação desejada (ex: create)"
// @Success      202      {object} map[string]string "Arquivo recebido, processamento em andamento"
// @Failure      400      {object} map[string]string "Erro nos parâmetros ou no upload"
// @Router       /api/v1/upload [post]
func NewUploadHandler(service inbound.UploadService) *UploadHandler {
	return &UploadHandler{service: service}
}

func (h *UploadHandler) ProcessFileHandler(c *gin.Context) {
	resource := c.PostForm("resource")
	action := c.PostForm("action")
	model := c.PostForm("model")

	if resource == "" || action == "" || model == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parâmetros 'resource', 'model' e 'action' são obrigatórios"})
		return
	}

	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Falha ao obter o arquivo"})
		return
	}
	defer file.Close()

	if err := h.service.ImportFile(resource, model, action, file, fileHeader); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar o arquivo"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "Arquivo recebido, processamento em andamento"})
}
