package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// --- Models ---

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	City      string  `json:"city,omitempty"`
	Country   string  `json:"country,omitempty"`
}

type Sighting struct {
	ID         string   `json:"id"`
	KaijuID    string   `json:"kaiju_id"`
	Location   Location `json:"location"`
	ReportedAt string   `json:"reported_at"`
	Confirmed  bool     `json:"confirmed"`
	Notes      string   `json:"notes,omitempty"`
}

type Kaiju struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Species      string  `json:"species"`
	HeightMeters float64 `json:"height_meters"`
	WeightTonnes float64 `json:"weight_tonnes"`
	ThreatLevel  string  `json:"threat_level"`
	FirstSeen    string  `json:"first_seen"`
	Status       string  `json:"status"`
	Description  string  `json:"description,omitempty"`
}

type Pagination struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalCount int `json:"total_count"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// --- Seed data ---

var kaijuDB = []Kaiju{
	{
		ID:           "k-001",
		Name:         "Goraxus",
		Species:      "Mega Primate",
		HeightMeters: 102,
		WeightTonnes: 58000,
		ThreatLevel:  "omega",
		FirstSeen:    "1954-11-03",
		Status:       "active",
		Description:  "A colossal primate first observed near Skull Island. Known for territorial behavior and extraordinary climbing ability.",
	},
	{
		ID:           "k-002",
		Name:         "Tidestrider",
		Species:      "Leviathan",
		HeightMeters: 78,
		WeightTonnes: 42000,
		ThreatLevel:  "gamma",
		FirstSeen:    "1998-05-20",
		Status:       "dormant",
		Description:  "An enormous aquatic creature that patrols deep ocean trenches. Rarely surfaces but causes massive tidal disruptions when it does.",
	},
	{
		ID:           "k-003",
		Name:         "Volcanor",
		Species:      "Igneous Titan",
		HeightMeters: 135,
		WeightTonnes: 75000,
		ThreatLevel:  "omega",
		FirstSeen:    "1971-03-12",
		Status:       "active",
		Description:  "Born from volcanic activity in the Pacific Ring of Fire. Generates extreme heat and has been linked to several eruptions.",
	},
	{
		ID:           "k-004",
		Name:         "Skytalon",
		Species:      "Avian Rex",
		HeightMeters: 55,
		WeightTonnes: 18000,
		ThreatLevel:  "beta",
		FirstSeen:    "2010-07-08",
		Status:       "active",
		Description:  "A winged kaiju capable of supersonic flight. Nests in mountain ranges and hunts over wide territories.",
	},
	{
		ID:           "k-005",
		Name:         "Crustara",
		Species:      "Arthropod Colossus",
		HeightMeters: 40,
		WeightTonnes: 12000,
		ThreatLevel:  "alpha",
		FirstSeen:    "2020-01-15",
		Status:       "neutralized",
		Description:  "A heavily armored crustacean kaiju that emerged from the Mariana Trench. Successfully neutralized in 2023.",
	},
}

var sightingsDB = []Sighting{
	{
		ID:      "s-1001",
		KaijuID: "k-001",
		Location: Location{
			Latitude: 35.6762, Longitude: 139.6503,
			City: "Tokyo", Country: "Japan",
		},
		ReportedAt: "2024-08-15T14:30:00Z",
		Confirmed:  true,
		Notes:      "Emerged from Tokyo Bay at dawn. Headed northwest.",
	},
	{
		ID:      "s-1002",
		KaijuID: "k-001",
		Location: Location{
			Latitude: 37.7749, Longitude: -122.4194,
			City: "San Francisco", Country: "United States",
		},
		ReportedAt: "2025-01-10T08:15:00Z",
		Confirmed:  false,
		Notes:      "Unconfirmed sonar contact beneath the Golden Gate.",
	},
	{
		ID:      "s-1003",
		KaijuID: "k-002",
		Location: Location{
			Latitude: -33.8688, Longitude: 151.2093,
			City: "Sydney", Country: "Australia",
		},
		ReportedAt: "2024-12-01T22:00:00Z",
		Confirmed:  true,
		Notes:      "Massive displacement wave detected off the coast. Visual confirmation by naval patrol.",
	},
	{
		ID:      "s-1004",
		KaijuID: "k-003",
		Location: Location{
			Latitude: 19.4326, Longitude: -99.1332,
			City: "Mexico City", Country: "Mexico",
		},
		ReportedAt: "2025-02-20T06:45:00Z",
		Confirmed:  true,
		Notes:      "Seismic readings followed by visual sighting near Popocatépetl volcano.",
	},
	{
		ID:      "s-1005",
		KaijuID: "k-003",
		Location: Location{
			Latitude: 35.3606, Longitude: 138.7274,
			City: "Mount Fuji", Country: "Japan",
		},
		ReportedAt: "2024-06-10T03:20:00Z",
		Confirmed:  false,
		Notes:      "Thermal anomaly detected. Unconfirmed aerial sighting.",
	},
	{
		ID:      "s-1006",
		KaijuID: "k-004",
		Location: Location{
			Latitude: 27.9881, Longitude: 86.9250,
			City: "Everest Region", Country: "Nepal",
		},
		ReportedAt: "2025-03-01T11:00:00Z",
		Confirmed:  true,
		Notes:      "Spotted circling above the summit. Sonic boom reported by climbing expedition.",
	},
}

// --- Helpers ---

var validThreatLevels = map[string]bool{
	"alpha": true, "beta": true, "gamma": true, "omega": true,
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, code int, message string) {
	writeJSON(w, code, ErrorResponse{Code: code, Message: message})
}

func parsePagination(r *http.Request) (page, pageSize int, err error) {
	page = 1
	pageSize = 10

	if v := r.URL.Query().Get("page"); v != "" {
		page, err = strconv.Atoi(v)
		if err != nil || page < 1 {
			return 0, 0, fmt.Errorf("invalid query parameter: page must be a positive integer")
		}
	}
	if v := r.URL.Query().Get("page_size"); v != "" {
		pageSize, err = strconv.Atoi(v)
		if err != nil || pageSize < 1 || pageSize > 100 {
			return 0, 0, fmt.Errorf("invalid query parameter: page_size must be between 1 and 100")
		}
	}
	return page, pageSize, nil
}

func paginate[T any](items []T, page, pageSize int) ([]T, Pagination) {
	total := len(items)
	start := (page - 1) * pageSize
	if start > total {
		start = total
	}
	end := start + pageSize
	if end > total {
		end = total
	}
	return items[start:end], Pagination{
		Page:       page,
		PageSize:   pageSize,
		TotalCount: total,
	}
}

// --- Handlers ---

func handleListKaiju(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed.")
		return
	}

	page, pageSize, err := parsePagination(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	threatLevel := r.URL.Query().Get("threat_level")
	if threatLevel != "" && !validThreatLevels[threatLevel] {
		writeError(w, http.StatusBadRequest,
			"Invalid query parameter: threat_level must be one of alpha, beta, gamma, omega.")
		return
	}

	filtered := kaijuDB
	if threatLevel != "" {
		filtered = nil
		for _, k := range kaijuDB {
			if k.ThreatLevel == threatLevel {
				filtered = append(filtered, k)
			}
		}
	}

	data, pagination := paginate(filtered, page, pageSize)
	writeJSON(w, http.StatusOK, map[string]any{
		"data":       data,
		"pagination": pagination,
	})
}

func handleGetKaiju(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed.")
		return
	}

	for _, k := range kaijuDB {
		if k.ID == id {
			writeJSON(w, http.StatusOK, k)
			return
		}
	}
	writeError(w, http.StatusNotFound, "Kaiju not found.")
}

func handleListSightings(w http.ResponseWriter, r *http.Request, kaijuID string) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed.")
		return
	}

	found := false
	for _, k := range kaijuDB {
		if k.ID == kaijuID {
			found = true
			break
		}
	}
	if !found {
		writeError(w, http.StatusNotFound, "Kaiju not found.")
		return
	}

	page, pageSize, err := parsePagination(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	var sightings []Sighting
	for _, s := range sightingsDB {
		if s.KaijuID == kaijuID {
			sightings = append(sightings, s)
		}
	}

	data, pagination := paginate(sightings, page, pageSize)
	writeJSON(w, http.StatusOK, map[string]any{
		"data":       data,
		"pagination": pagination,
	})
}

// --- Router ---

func router(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimSuffix(r.URL.Path, "/")

	// GET /v1/kaiju
	if path == "/v1/kaiju" {
		handleListKaiju(w, r)
		return
	}

	// GET /v1/kaiju/{kaijuId}/sightings
	if strings.HasPrefix(path, "/v1/kaiju/") && strings.HasSuffix(path, "/sightings") {
		id := strings.TrimPrefix(path, "/v1/kaiju/")
		id = strings.TrimSuffix(id, "/sightings")
		if id != "" {
			handleListSightings(w, r, id)
			return
		}
	}

	// GET /v1/kaiju/{kaijuId}
	if strings.HasPrefix(path, "/v1/kaiju/") {
		id := strings.TrimPrefix(path, "/v1/kaiju/")
		if id != "" && !strings.Contains(id, "/") {
			handleGetKaiju(w, r, id)
			return
		}
	}

	writeError(w, http.StatusNotFound, "Not found.")
}

func main() {
	http.HandleFunc("/", router)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port
	log.Printf("Kaiju Registry API listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
