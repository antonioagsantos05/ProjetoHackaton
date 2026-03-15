package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/minio/minio-go/v7"
	"github.com/streadway/amqp"
)

type HealthHandler struct {
	db          *sql.DB
	rabbitConn  *amqp.Connection
	minioClient *minio.Client
}

func NewHealthHandler(db *sql.DB, rabbitConn *amqp.Connection, minioClient *minio.Client) *HealthHandler {
	return &HealthHandler{
		db:          db,
		rabbitConn:  rabbitConn,
		minioClient: minioClient,
	}
}

// CheckHealth godoc
// @Summary      Verifica a saúde da API e suas dependências
// @Description  Retorna o status da API e de serviços como Banco de Dados, Mensageria e Storage.
// @Tags         health
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      503  {object}  map[string]interface{}
// @Router       /health [get]
func (h *HealthHandler) CheckHealth(w http.ResponseWriter, r *http.Request) {
	dependencies := make(map[string]string)
	overallStatus := "UP"
	httpStatus := http.StatusOK

	// 1. Verificar Banco de Dados
	if err := h.db.Ping(); err != nil {
		dependencies["database"] = "DOWN"
		overallStatus = "DOWN"
	} else {
		dependencies["database"] = "UP"
	}

	// 2. Verificar RabbitMQ
	if h.rabbitConn.IsClosed() {
		dependencies["messaging"] = "DOWN"
		overallStatus = "DOWN"
	} else {
		dependencies["messaging"] = "UP"
	}

	// 3. Verificar Minio
	_, err := h.minioClient.BucketExists(context.Background(), "videos")
	if err != nil {
		dependencies["storage"] = "DOWN"
		overallStatus = "DOWN"
	} else {
		dependencies["storage"] = "UP"
	}

	// Se qualquer dependência estiver fora, o status geral é Service Unavailable
	if overallStatus == "DOWN" {
		httpStatus = http.StatusServiceUnavailable
	}

	response := map[string]interface{}{
		"status":       overallStatus,
		"dependencies": dependencies,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(response)
}
