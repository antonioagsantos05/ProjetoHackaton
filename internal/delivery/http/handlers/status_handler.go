package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/fiap-x/video-processor/internal/usecase/job"
	"github.com/gorilla/mux"
)

type StatusHandler struct {
	statusUseCase job.StatusUseCase
}

func NewStatusHandler(uc job.StatusUseCase) *StatusHandler {
	return &StatusHandler{statusUseCase: uc}
}

// CheckJobStatus godoc
// @Summary      Verifica o status de um job
// @Description  Retorna os detalhes de um job de processamento e do vídeo associado.
// @Tags         jobs
// @Produce      json
// @Param        id   path      string  true  "ID do Job"
// @Success      200  {object}  job.StatusResponse
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /jobs/{id}/status [get]
func (h *StatusHandler) CheckJobStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobID := vars["id"]

	// Chamar o Use Case
	response, err := h.statusUseCase.Execute(jobID)
	if err != nil {
		log.Printf("Erro ao executar o use case de status: %v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Retornar a resposta completa
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
