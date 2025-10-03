package api

import (
	"math"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/agnosto/fansly-scraper/internal/service"
)

// Creator represents a content creator
type Creator = service.Creator

// handleListCreators handles GET /api/v1/creators
// Query parameters:
//   - limit: number of creators to return (default: 20, max: 100)
//   - offset: number of creators to skip (default: 0)
//   - sort: field to sort by (name, last_updated, default: name)
//   - order: sort order (asc, desc, default: asc)
func (s *Server) handleListCreators(w http.ResponseWriter, r *http.Request) {
	// Parse and validate query parameters
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		limit = 20 // Default limit
	}
	if limit > 100 {
		limit = 100 // Max limit
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil || offset < 0 {
		offset = 0
	}

	sortBy := r.URL.Query().Get("sort")
	switch sortBy {
	case "", "name", "last_updated":
		// Valid sort fields
	default:
		respondWithError(w, http.StatusBadRequest, "Invalid sort field. Must be one of: name, last_updated")
		return
	}

	order := r.URL.Query().Get("order")
	if order != "" && order != "asc" && order != "desc" {
		respondWithError(w, http.StatusBadRequest, "Invalid order. Must be one of: asc, desc")
	}

	s.log.Infof("Listing creators with limit=%d, offset=%d, sort=%s, order=%s", 
		limit, offset, sortBy, order)

	// Get creators from the scraper service
	creators, err := s.scraperSvc.GetCreators(r.Context(), limit, offset)
	if err != nil {
		s.log.Errorf("Failed to get creators: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch creators")
		return
	}

	// Apply sorting
	if sortBy == "last_updated" {
		sort.Slice(creators, func(i, j int) bool {
			if order == "desc" {
				return creators[i].LastUpdated.After(creators[j].LastUpdated)
			}
			return creators[i].LastUpdated.Before(creators[j].LastUpdated)
		})
	} else { // sort by name (default)
		sort.Slice(creators, func(i, j int) bool {
			if order == "desc" {
				return creators[i].Name > creators[j].Name
			}
			return creators[i].Name < creators[j].Name
		})
	}

	total := len(creators)
	// Prepare response
	response := map[string]interface{}{
		"data": creators,
		"meta": map[string]interface{}{
			"total":       total,
			"count":       len(creators),
			"per_page":    limit,
			"current_page": (offset / limit) + 1,
			"total_pages": int(math.Ceil(float64(total) / float64(limit))),
		},
	}
	respondWithJSON(w, http.StatusOK, response)
}

// handleGetCreatorContent gets content from a specific creator
func (s *Server) handleGetCreatorContent(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement getting creator content
	creatorID := r.URL.Query().Get("creatorId")
	s.log.Infof("Get creator content for: %s", creatorID)
	respondWithJSON(w, http.StatusNotImplemented, map[string]string{
		"message":  "Get creator content not yet implemented",
		"creatorId": creatorID,
	})
}

// handleGetMedia gets specific media content
func (s *Server) handleGetMedia(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement getting media
	creatorID := r.URL.Query().Get("creatorId")
	mediaID := r.URL.Query().Get("mediaId")

	s.log.Infof("Get media %s for creator %s", mediaID, creatorID)
	respondWithJSON(w, http.StatusNotImplemented, map[string]string{
		"message":  "Get media not yet implemented",
		"creatorId": creatorID,
		"mediaId":  mediaID,
	})
}

// handleStartMonitoring starts monitoring a creator for new content
func (s *Server) handleStartMonitoring(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement start monitoring
	s.log.Info("Start monitoring endpoint hit")
	respondWithJSON(w, http.StatusNotImplemented, map[string]string{
		"message": "Start monitoring not yet implemented",
	})
}

// handleStopMonitoring stops monitoring
func (s *Server) handleStopMonitoring(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement stop monitoring
	s.log.Info("Stop monitoring endpoint hit")
	respondWithJSON(w, http.StatusNotImplemented, map[string]string{
		"message": "Stop monitoring not yet implemented",
	})
}
