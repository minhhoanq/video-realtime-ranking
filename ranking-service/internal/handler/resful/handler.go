package resful

import (
	"database/sql"
	"errors"
	"net/http"
	"video-realtime-ranking/ranking-service/internal/service"
)

type Handler struct {
	rankingService service.RankingService
}

func NewHandler(rankingService service.RankingService) *Handler {
	return &Handler{
		rankingService: rankingService,
	}
}

func (h *Handler) GetTopK(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	limit := readInt(r.URL.Query(), "limit", 10)
	offset := readInt(r.URL.Query(), "offset", 0)

	videoScore, err := h.rankingService.GetTopKVideoRanking(ctx, offset, limit)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			Error(w, err, http.StatusNotFound)
		}
		Error(w, err, http.StatusBadRequest)
	}

	Success(w, videoScore, http.StatusOK, nil)
}
