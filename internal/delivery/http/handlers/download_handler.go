package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/fiap-x/video-processor/internal/usecase/video"
	"github.com/gorilla/mux"
)

type DownloadHandler struct {
	downloadUseCase video.DownloadUseCase
}

func NewDownloadHandler(uc video.DownloadUseCase) *DownloadHandler {
	return &DownloadHandler{downloadUseCase: uc}
}

// GetZip godoc
// @Summary      Obtém a URL de download de um arquivo ZIP
// @Description  Retorna uma URL pré-assinada para o download do arquivo ZIP processado.
// @Tags         downloads
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "ID do Job"
// @Success      200  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      422  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /downloads/{id} [get]
func (h *DownloadHandler) GetZip(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobID := vars["id"]

	// Chamar o Use Case
	url, err := h.downloadUseCase.Execute(jobID)
	if err != nil {
		log.Printf("Erro ao executar o use case de download: %v", err)
		// TODO: Tratar os diferentes tipos de erro para retornar o status code correto (404, 422, 500)
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	// Retornar a URL
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"url": url,
	})
}
