package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/fiap-x/video-processor/internal/delivery/http/middleware"
	"github.com/fiap-x/video-processor/internal/usecase/video"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type VideoHandler struct {
	uploadUseCase video.UploadUseCase
	listUseCase   video.ListUseCase
}

func NewVideoHandler(uploadUC video.UploadUseCase, listUC video.ListUseCase) *VideoHandler {
	return &VideoHandler{
		uploadUseCase: uploadUC,
		listUseCase:   listUC,
	}
}

// UploadVideo godoc
// @Summary      Realiza o upload de um vídeo
// @Description  Recebe um arquivo de vídeo para processamento e retorna um ID de job. Requer autenticação JWT.
// @Tags         videos
// @Accept       multipart/form-data
// @Produce      json
// @Security     BearerAuth
// @Param        video formData file true "Arquivo de vídeo para upload"
// @Success      202   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Router       /videos/upload [post]
func (h *VideoHandler) UploadVideo(w http.ResponseWriter, r *http.Request) {
	// 1. Extrair IDs do contexto (colocados pelo AuthMiddleware)
	userID, err := uuid.Parse(r.Context().Value(middleware.UserIDKey).(string))
	if err != nil {
		http.Error(w, "ID de usuário inválido no token", http.StatusUnauthorized)
		return
	}
	tenantID, err := uuid.Parse(r.Context().Value(middleware.TenantIDKey).(string))
	if err != nil {
		http.Error(w, "ID de tenant inválido no token", http.StatusUnauthorized)
		return
	}

	// 2. Parse do formulário multipart
	r.ParseMultipartForm(500 << 20) // Limite de 500MB

	file, handler, err := r.FormFile("video")
	if err != nil {
		http.Error(w, "Erro ao ler o arquivo: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// 3. Salvar o arquivo temporariamente
	uploadDir := "./uploads"
	os.MkdirAll(uploadDir, os.ModePerm)
	filePath := filepath.Join(uploadDir, handler.Filename)

	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Não foi possível criar o arquivo no servidor: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Não foi possível salvar o arquivo: "+err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Arquivo '%s' salvo temporariamente em '%s' para o usuário %s", handler.Filename, filePath, userID)

	// 4. Chamar o Use Case com o caminho do arquivo e os IDs
	job, err := h.uploadUseCase.Execute(filePath, userID, tenantID)
	if err != nil {
		log.Printf("Erro ao executar o use case de upload: %v", err)
		http.Error(w, "Erro interno ao processar o upload.", http.StatusInternalServerError)
		return
	}

	// 5. Retornar a resposta
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{
		"jobId":   job.ID.String(), // Convertendo o UUID para string
		"message": fmt.Sprintf("Upload recebido. O processamento do arquivo '%s' foi iniciado.", handler.Filename),
	})
}

// ListVideos godoc
// @Summary      Lista os vídeos de um usuário
// @Description  Retorna uma lista de vídeos com seus status para um determinado usuário.
// @Tags         users
// @Produce      json
// @Param        id   path      string  true  "ID do Usuário"
// @Success      200  {array}   video.VideoDTO
// @Router       /users/{id}/videos [get]
func (h *VideoHandler) ListVideos(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	videos, err := h.listUseCase.Execute(userID)
	if err != nil {
		log.Printf("Erro ao listar vídeos do usuário %s: %v", userID, err)
		http.Error(w, "Erro ao buscar os vídeos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(videos)
}
