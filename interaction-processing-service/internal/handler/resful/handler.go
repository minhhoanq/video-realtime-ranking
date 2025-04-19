package resful

import (
	"encoding/json"
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
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate (optional)
	if req.UserID == "" || req.VideoID == "" || req.InteractionType == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Tạo argument cho data layer
	message := producer.InteractionCreate{
		UserID:          req.UserID,
		VideoID:         req.VideoID,
		InteractionType: req.InteractionType,
	}

	err := h.interactionCreateKafkaProducer.Produce(ctx, message)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(err)
}
