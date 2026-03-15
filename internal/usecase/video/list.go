package video

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// VideoDTO representa os dados do vídeo que serão retornados na listagem.
type VideoDTO struct {
	ID        uuid.UUID  `json:"id"`
	Title     string     `json:"title"`
	Status    int16      `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
}

// ListUseCase define a interface para o caso de uso de listagem de vídeos de um usuário.
type ListUseCase interface {
	Execute(userID string) ([]VideoDTO, error)
}

// listUseCase é a implementação do ListUseCase.
type listUseCase struct {
	db *sql.DB
}

// NewListUseCase cria uma nova instância do listUseCase.
func NewListUseCase(db *sql.DB) ListUseCase {
	return &listUseCase{db: db}
}

// Execute orquestra o processo de listagem de vídeos.
func (uc *listUseCase) Execute(userID string) ([]VideoDTO, error) {
	query := `
		SELECT id, titulo, status, created_at
		FROM video
		WHERE usuario_id = $1
		ORDER BY created_at DESC
	`
	rows, err := uc.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar vídeos no banco de dados: %w", err)
	}
	defer rows.Close()

	var videos []VideoDTO
	for rows.Next() {
		var v VideoDTO
		if err := rows.Scan(&v.ID, &v.Title, &v.Status, &v.CreatedAt); err != nil {
			return nil, fmt.Errorf("erro ao escanear linha de vídeo: %w", err)
		}
		videos = append(videos, v)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("erro após iterar sobre os vídeos: %w", err)
	}

	// Se não houver vídeos, retornamos um slice vazio em vez de nil para que o JSON fique "[]"
	if videos == nil {
		videos = []VideoDTO{}
	}

	return videos, nil
}
