package handlers

import (
	"net/http"
	"github.com/gorilla/mux"
)

type DownloadHandler struct {}

func NewDownloadHandler() *DownloadHandler {
	return &DownloadHandler{}
}

// GetZip godoc
// @Summary      Obtém a URL de download de um arquivo ZIP
// @Description  Retorna uma URL pré-assinada para o download do arquivo ZIP processado.
// @Tags         downloads
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "ID do Job"
// @Success      200  {object}  map[string]string
// @Router       /downloads/{id} [get]
func (h *DownloadHandler) GetZip(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	zipID := vars["id"]

	// Buscar jobID com status Done
	// Gerar Presigned URL S3
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"url": "https://s3.aws.com/fiapx/` + zipID + `.zip"}`))
}
