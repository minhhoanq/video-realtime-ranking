package resful

import (
	"encoding/json"
	"errors"
	"net/http"
	"video-realtime-ranking/interaction-processing-service/internal/dataaccess/database"
	"video-realtime-ranking/interaction-processing-service/internal/dataaccess/kafka/producer"
	"video-realtime-ranking/interaction-processing-service/internal/service"
)

type Handler struct {
	interactionService             service.InteractionService
	interactionCreateKafkaProducer producer.InteractionCreateProducer
}

func NewHandler(interactionService service.InteractionService,
	interactionCreateKafkaProducer producer.InteractionCreateProducer,
) *Handler {
	return &Handler{
		interactionService:             interactionService,
		interactionCreateKafkaProducer: interactionCreateKafkaProducer,
	}
}

func (h *Handler) CreateInteraction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// Parse JSON body
	var req database.SendInteractionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, err, http.StatusBadRequest)
		return
	}

	// Validate (optional)
	if req.UserID == "" || req.VideoID == "" || req.InteractionType == "" {
		Error(w, errors.New("Missing required fields"), http.StatusBadRequest)
		return
	}

	// Tạo argument cho data layer
	message := producer.InteractionCreate{
		UserID:          req.UserID,
		VideoID:         req.VideoID,
		InteractionType: req.InteractionType,
	}

	err := h.interactionCreateKafkaProducer.Produce(ctx, message)
	if err != nil {
		Error(w, err, http.StatusBadRequest)
	}

	Success(w, err, http.StatusOK, nil)
}
