package handlers

import (
	"encoding/json"
	"net/http"
	"github.com/google/uuid"
	"log"
)

type VideoHandler struct {
	// Aqui entraria o seu usecase de upload
	// uploadUseCase usecase.UploadUseCase
}

func NewVideoHandler() *VideoHandler {
	return &VideoHandler{}
}

// UploadVideo godoc
// @Summary      Realiza o upload de um vídeo
// @Description  Recebe um arquivo de vídeo para processamento e retorna um ID de job.
// @Tags         videos
// @Accept       multipart/form-data
// @Produce      json
// @Param        video formData file true "Arquivo de vídeo para upload"
// @Success      202   {object}  map[string]string
// @Router       /videos/upload [post]
func (h *VideoHandler) UploadVideo(w http.ResponseWriter, r *http.Request) {
	// Define um limite para o tamanho do arquivo (ex: 500MB)
	r.ParseMultipartForm(500 << 20)

	file, handler, err := r.FormFile("video")
	if err != nil {
		log.Println("Erro ao obter o arquivo:", err)
		http.Error(w, "Erro ao ler o arquivo enviado.", http.StatusBadRequest)
		return
	}
	defer file.Close()

	log.Printf("Arquivo recebido: %+v\n", handler.Filename)
	log.Printf("Tamanho do arquivo: %+v bytes\n", handler.Size)

	// Aqui você chamaria o seu use case para salvar o arquivo no Minio
	// e publicar uma mensagem no RabbitMQ.

	jobID := uuid.New().String()
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{
		"jobId":   jobID,
		"message": "Upload recebido. O processamento foi iniciado.",
	})
}
