package handlers

import (
	"encoding/json"
	"net/http"
	"github.com/gorilla/mux"
)

type StatusHandler struct {}

func NewStatusHandler() *StatusHandler {
	return &StatusHandler{}
}

// CheckJobStatus godoc
// @Summary      Verifica o status de um job
// @Description  Retorna o status atual de um job de processamento de vídeo.
// @Tags         jobs
// @Produce      json
// @Param        id   path      string  true  "ID do Job"
// @Success      200  {object}  map[string]string
// @Router       /jobs/{id}/status [get]
func (h *StatusHandler) CheckJobStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobID := vars["id"]

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"id": jobID,
		"status": "PROCESSING",
		"progress": "45%",
	})
}
