package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jerryluo/cmd/internal/logging"
)

// LogListResponse is the response for the list logs endpoint
type LogListResponse struct {
	Logs   []logging.LogSummary `json:"logs"`
	Total  int                  `json:"total"`
	Limit  int                  `json:"limit"`
	Offset int                  `json:"offset"`
}

// ErrorResponse is a standard error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// handleListLogs handles GET /api/v1/logs
func handleListLogs(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	// Get all logs
	logs, err := logging.ListLogs()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to list logs: "+err.Error())
		return
	}

	// Apply search filter
	if search := query.Get("search"); search != "" {
		logs, err = logging.SearchLogs(search)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Failed to search logs: "+err.Error())
			return
		}
	}

	// Apply status filter
	if status := query.Get("status"); status != "" {
		filtered := make([]logging.LogSummary, 0)
		for _, log := range logs {
			if string(log.FinalStatus) == status {
				filtered = append(filtered, log)
			}
		}
		logs = filtered
	}

	// Apply model filter
	if model := query.Get("model"); model != "" {
		filtered := make([]logging.LogSummary, 0)
		for _, log := range logs {
			if log.Model == model {
				filtered = append(filtered, log)
			}
		}
		logs = filtered
	}

	// Apply time range filters
	if from := query.Get("from"); from != "" {
		fromTime, err := time.Parse(time.RFC3339, from)
		if err == nil {
			filtered := make([]logging.LogSummary, 0)
			for _, log := range logs {
				if !log.Timestamp.Before(fromTime) {
					filtered = append(filtered, log)
				}
			}
			logs = filtered
		}
	}

	if to := query.Get("to"); to != "" {
		toTime, err := time.Parse(time.RFC3339, to)
		if err == nil {
			filtered := make([]logging.LogSummary, 0)
			for _, log := range logs {
				if !log.Timestamp.After(toTime) {
					filtered = append(filtered, log)
				}
			}
			logs = filtered
		}
	}

	// Apply sorting
	sortField := query.Get("sort")
	sortOrder := query.Get("order")
	if sortOrder == "" {
		sortOrder = "desc"
	}

	if sortField != "" {
		sortLogs(logs, sortField, sortOrder == "asc")
	}

	// Store total before pagination
	total := len(logs)

	// Apply pagination
	limit := 100
	if l := query.Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 1000 {
			limit = parsed
		}
	}

	offset := 0
	if o := query.Get("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	// Apply offset and limit
	if offset >= len(logs) {
		logs = []logging.LogSummary{}
	} else {
		end := offset + limit
		if end > len(logs) {
			end = len(logs)
		}
		logs = logs[offset:end]
	}

	response := LogListResponse{
		Logs:   logs,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}

	writeJSON(w, http.StatusOK, response)
}

// handleGetLog handles GET /api/v1/logs/{id}
func handleGetLog(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "Log ID is required")
		return
	}

	log, err := logging.ReadLogWithID(id)
	if err != nil {
		writeError(w, http.StatusNotFound, "Log not found: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, log)
}

// sortLogs sorts the logs by the given field
func sortLogs(logs []logging.LogSummary, field string, ascending bool) {
	for i := 0; i < len(logs)-1; i++ {
		for j := i + 1; j < len(logs); j++ {
			swap := false

			switch field {
			case "timestamp":
				if ascending {
					swap = logs[j].Timestamp.Before(logs[i].Timestamp)
				} else {
					swap = logs[j].Timestamp.After(logs[i].Timestamp)
				}
			case "status":
				if ascending {
					swap = string(logs[j].FinalStatus) < string(logs[i].FinalStatus)
				} else {
					swap = string(logs[j].FinalStatus) > string(logs[i].FinalStatus)
				}
			case "model":
				if ascending {
					swap = logs[j].Model < logs[i].Model
				} else {
					swap = logs[j].Model > logs[i].Model
				}
			case "query":
				if ascending {
					swap = strings.ToLower(logs[j].UserQuery) < strings.ToLower(logs[i].UserQuery)
				} else {
					swap = strings.ToLower(logs[j].UserQuery) > strings.ToLower(logs[i].UserQuery)
				}
			}

			if swap {
				logs[i], logs[j] = logs[j], logs[i]
			}
		}
	}
}

// writeJSON writes a JSON response
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// writeError writes a JSON error response
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, ErrorResponse{Error: message})
}
