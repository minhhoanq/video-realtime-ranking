package resful

import (
	"encoding/json"
	"net/http"
	"video-realtime-ranking/interaction-processing-service/internal/dataaccess/database"
	"video-realtime-ranking/interaction-processing-service/internal/service"
)

type Handler struct {
	interactionService service.InteractionService
}

func NewHandler(interactionService service.InteractionService) *Handler {
	return &Handler{
		interactionService: interactionService,
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
	arg := database.SendInteractionRequest{
		UserID:          req.UserID,
		VideoID:         req.VideoID,
		InteractionType: req.InteractionType,
	}

	interaction, err := h.interactionService.CreateInteraction(ctx, &arg)
	if err != nil {
		http.Error(w, "Failed to create interaction", http.StatusInternalServerError)
		return
	}

	response := &database.SendInteractionResponse{
		UserID:  interaction.UserID,
		VideoID: interaction.UserID,
		ID:      interaction.ID,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
