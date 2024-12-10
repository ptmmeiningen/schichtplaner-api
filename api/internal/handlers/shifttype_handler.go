package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"shift-planner/api/internal/models"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type ShiftTypeHandler struct {
	db *gorm.DB
}

func NewShiftTypeHandler(db *gorm.DB) *ShiftTypeHandler {
	return &ShiftTypeHandler{db: db}
}

func (h *ShiftTypeHandler) GetShiftTypes(w http.ResponseWriter, r *http.Request) {
	var shiftTypes []models.ShiftType
	result := h.db.Order("name ASC").Find(&shiftTypes)

	if result.Error != nil {
		SendErrorResponse(w, http.StatusInternalServerError, "Fehler beim Abrufen der Schichttypen")
		log.Printf("GetShiftTypes Error: %v", result.Error)
		return
	}

	SendSuccessResponse(w, "Schichttypen erfolgreich abgerufen", shiftTypes)
	log.Printf("GetShiftTypes: %d Schichttypen abgerufen", len(shiftTypes))
}

func (h *ShiftTypeHandler) CreateShiftType(w http.ResponseWriter, r *http.Request) {
	var shiftType models.ShiftType
	if err := json.NewDecoder(r.Body).Decode(&shiftType); err != nil {
		SendErrorResponse(w, http.StatusBadRequest, "Ungültige Eingabedaten")
		log.Printf("CreateShiftType Decode Error: %v", err)
		return
	}

	result := h.db.Create(&shiftType)
	if result.Error != nil {
		SendErrorResponse(w, http.StatusInternalServerError, "Fehler beim Erstellen des Schichttyps")
		log.Printf("CreateShiftType DB Error: %v", result.Error)
		return
	}

	SendSuccessResponse(w, "Schichttyp erfolgreich erstellt", shiftType)
	log.Printf("CreateShiftType: Schichttyp %s erstellt", shiftType.Name)
}

func (h *ShiftTypeHandler) GetShiftType(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var shiftType models.ShiftType

	result := h.db.First(&shiftType, params["id"])
	if result.Error != nil {
		SendErrorResponse(w, http.StatusNotFound, "Schichttyp nicht gefunden")
		log.Printf("GetShiftType Error: ID %s nicht gefunden", params["id"])
		return
	}

	SendSuccessResponse(w, "Schichttyp erfolgreich abgerufen", shiftType)
	log.Printf("GetShiftType: Schichttyp %s abgerufen", shiftType.Name)
}

func (h *ShiftTypeHandler) UpdateShiftType(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var shiftType models.ShiftType

	if result := h.db.First(&shiftType, params["id"]); result.Error != nil {
		SendErrorResponse(w, http.StatusNotFound, "Schichttyp nicht gefunden")
		log.Printf("UpdateShiftType Error: ID %s nicht gefunden", params["id"])
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&shiftType); err != nil {
		SendErrorResponse(w, http.StatusBadRequest, "Ungültige Eingabedaten")
		log.Printf("UpdateShiftType Decode Error: %v", err)
		return
	}

	if result := h.db.Save(&shiftType); result.Error != nil {
		SendErrorResponse(w, http.StatusInternalServerError, "Fehler beim Aktualisieren des Schichttyps")
		log.Printf("UpdateShiftType DB Error: %v", result.Error)
		return
	}

	SendSuccessResponse(w, "Schichttyp erfolgreich aktualisiert", shiftType)
	log.Printf("UpdateShiftType: Schichttyp %s aktualisiert", shiftType.Name)
}

func (h *ShiftTypeHandler) DeleteShiftType(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var shiftType models.ShiftType

	if result := h.db.First(&shiftType, params["id"]); result.Error != nil {
		SendErrorResponse(w, http.StatusNotFound, "Schichttyp nicht gefunden")
		log.Printf("DeleteShiftType Error: ID %s nicht gefunden", params["id"])
		return
	}

	if result := h.db.Delete(&shiftType); result.Error != nil {
		SendErrorResponse(w, http.StatusInternalServerError, "Fehler beim Löschen des Schichttyps")
		log.Printf("DeleteShiftType DB Error: %v", result.Error)
		return
	}

	SendSuccessResponse(w, "Schichttyp erfolgreich gelöscht", nil)
	log.Printf("DeleteShiftType: Schichttyp %s gelöscht", shiftType.Name)
}

// Hilfsfunktionen für die Antwortverarbeitung
func SendErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ApiResponse{
		Success: false,
		Message: message,
		Data:    nil,
	})
}

func SendSuccessResponse(w http.ResponseWriter, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ApiResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}
