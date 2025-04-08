package restHandlers

import (
	"encoding/json"
	"net/http"

	"github.com/GPTA/src/models"
)

// GetAllAIPromptsHandler returns all AI prompts from the database
func GetAllAIPrompts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var prompts []models.AIPrompt
	if err := models.DB.Find(&prompts).Error; err != nil {
		http.Error(w, "Error fetching prompts: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert to JSON and send response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(prompts); err != nil {
		http.Error(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
