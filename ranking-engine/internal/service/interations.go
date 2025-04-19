package service

import (
	"encoding/json"
	"net/http"
	"video-realtime-ranking/ranking-engine/internal/dataaccess/database"
)

type InteractionService interface {
	CreateInteraction(w http.ResponseWriter, r *http.Request)
}

type interactionService struct {
	interactionDataAccessor database.InteractionDataAccessor
}

func NewInteractionService(
	interactionDataAccessor database.InteractionDataAccessor,
) InteractionService {
	return &interactionService{
		interactionDataAccessor: interactionDataAccessor,
	}
}

func (i *interactionService) CreateInteraction(w http.ResponseWriter, r *http.Request) {
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

	interaction, err := i.interactionDataAccessor.CreateInteraction(ctx, &arg)
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
