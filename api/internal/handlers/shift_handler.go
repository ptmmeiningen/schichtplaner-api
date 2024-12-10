package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"shift-planner/api/internal/models"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type ShiftHandler struct {
	db *gorm.DB
}

func NewShiftHandler(db *gorm.DB) *ShiftHandler {
	return &ShiftHandler{db: db}
}

func (h *ShiftHandler) GetShifts(w http.ResponseWriter, r *http.Request) {
	var shifts []models.Shift
	result := h.db.Preload("ShiftType").Find(&shifts)
	if result.Error != nil {
		SendErrorResponse(w, http.StatusInternalServerError, "Fehler beim Abrufen der Schichten")
		log.Printf("GetShifts DB Error: %v", result.Error)
		return
	}

	SendSuccessResponse(w, "Schichten erfolgreich abgerufen", shifts)
	log.Printf("GetShifts: %d Schichten abgerufen", len(shifts))
}

func (h *ShiftHandler) GetShift(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var shift models.Shift

	result := h.db.Preload("ShiftType").First(&shift, params["id"])
	if result.Error != nil {
		SendErrorResponse(w, http.StatusNotFound, "Schicht nicht gefunden")
		log.Printf("GetShift Error: ID %s nicht gefunden", params["id"])
		return
	}

	SendSuccessResponse(w, "Schicht erfolgreich abgerufen", shift)
	log.Printf("GetShift: Schicht ID %d abgerufen", shift.ID)
}

func (h *ShiftHandler) CreateShift(w http.ResponseWriter, r *http.Request) {
	var shift models.Shift
	if err := json.NewDecoder(r.Body).Decode(&shift); err != nil {
		SendErrorResponse(w, http.StatusBadRequest, "Ungültige Eingabedaten")
		log.Printf("CreateShift Decode Error: %v", err)
		return
	}

	// Validiere ShiftType
	var shiftType models.ShiftType
	if result := h.db.First(&shiftType, shift.ShiftTypeID); result.Error != nil {
		SendErrorResponse(w, http.StatusBadRequest, "Ungültiger Schichttyp")
		log.Printf("CreateShift: ShiftType ID %d nicht gefunden", shift.ShiftTypeID)
		return
	}

	if result := h.db.Create(&shift); result.Error != nil {
		SendErrorResponse(w, http.StatusInternalServerError, "Fehler beim Erstellen der Schicht")
		log.Printf("CreateShift DB Error: %v", result.Error)
		return
	}

	h.db.Preload("ShiftType").First(&shift, shift.ID)

	SendSuccessResponse(w, "Schicht erfolgreich erstellt", shift)
	log.Printf("CreateShift: Neue Schicht ID %d erstellt", shift.ID)
}

func (h *ShiftHandler) UpdateShift(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var shift models.Shift

	if result := h.db.First(&shift, params["id"]); result.Error != nil {
		SendErrorResponse(w, http.StatusNotFound, "Schicht nicht gefunden")
		log.Printf("UpdateShift Error: ID %s nicht gefunden", params["id"])
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&shift); err != nil {
		SendErrorResponse(w, http.StatusBadRequest, "Ungültige Eingabedaten")
		log.Printf("UpdateShift Decode Error: %v", err)
		return
	}

	if result := h.db.Save(&shift); result.Error != nil {
		SendErrorResponse(w, http.StatusInternalServerError, "Fehler beim Aktualisieren der Schicht")
		log.Printf("UpdateShift DB Error: %v", result.Error)
		return
	}

	h.db.Preload("ShiftType").First(&shift, shift.ID)

	SendSuccessResponse(w, "Schicht erfolgreich aktualisiert", shift)
	log.Printf("UpdateShift: Schicht ID %d aktualisiert", shift.ID)
}

func (h *ShiftHandler) DeleteShift(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var shift models.Shift

	if result := h.db.First(&shift, params["id"]); result.Error != nil {
		SendErrorResponse(w, http.StatusNotFound, "Schicht nicht gefunden")
		log.Printf("DeleteShift Error: ID %s nicht gefunden", params["id"])
		return
	}

	if result := h.db.Delete(&shift); result.Error != nil {
		SendErrorResponse(w, http.StatusInternalServerError, "Fehler beim Löschen der Schicht")
		log.Printf("DeleteShift DB Error: %v", result.Error)
		return
	}

	SendSuccessResponse(w, "Schicht erfolgreich gelöscht", nil)
	log.Printf("DeleteShift: Schicht ID %d gelöscht", shift.ID)
}
