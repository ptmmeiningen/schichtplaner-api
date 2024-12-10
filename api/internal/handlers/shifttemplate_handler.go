package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"shift-planner/api/internal/models"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type ShiftTemplateHandler struct {
	db *gorm.DB
}

func NewShiftTemplateHandler(db *gorm.DB) *ShiftTemplateHandler {
	return &ShiftTemplateHandler{db: db}
}

func (h *ShiftTemplateHandler) GetShiftTemplates(w http.ResponseWriter, r *http.Request) {
	var shiftTemplates []models.ShiftTemplate
	result := h.db.Preload("Employee").
		Order("name ASC").
		Find(&shiftTemplates)

	if result.Error != nil {
		SendErrorResponse(w, http.StatusInternalServerError, "Fehler beim Abrufen der Schichtvorlagen")
		log.Printf("GetShiftTemplates DB Error: %v", result.Error)
		return
	}

	SendSuccessResponse(w, "Schichtvorlagen erfolgreich abgerufen", shiftTemplates)
	log.Printf("GetShiftTemplates: %d Schichtvorlagen abgerufen", len(shiftTemplates))
}

func (h *ShiftTemplateHandler) CreateShiftTemplate(w http.ResponseWriter, r *http.Request) {
	var shiftTemplate models.ShiftTemplate
	if err := json.NewDecoder(r.Body).Decode(&shiftTemplate); err != nil {
		SendErrorResponse(w, http.StatusBadRequest, "Ungültige Eingabedaten")
		log.Printf("CreateShiftTemplate Decode Error: %v", err)
		return
	}

	if result := h.db.Create(&shiftTemplate); result.Error != nil {
		SendErrorResponse(w, http.StatusInternalServerError, "Fehler beim Erstellen der Schichtvorlage")
		log.Printf("CreateShiftTemplate DB Error: %v", result.Error)
		return
	}

	SendSuccessResponse(w, "Schichtvorlage erfolgreich erstellt", shiftTemplate)
	log.Printf("CreateShiftTemplate: Neue Schichtvorlage ID %d erstellt", shiftTemplate.ID)
}

func (h *ShiftTemplateHandler) GetShiftTemplate(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var shiftTemplate models.ShiftTemplate

	result := h.db.Preload("Employee").First(&shiftTemplate, params["id"])
	if result.Error != nil {
		SendErrorResponse(w, http.StatusNotFound, "Schichtvorlage nicht gefunden")
		log.Printf("GetShiftTemplate Error: ID %s nicht gefunden", params["id"])
		return
	}

	SendSuccessResponse(w, "Schichtvorlage erfolgreich abgerufen", shiftTemplate)
	log.Printf("GetShiftTemplate: Schichtvorlage ID %d abgerufen", shiftTemplate.ID)
}

func (h *ShiftTemplateHandler) UpdateShiftTemplate(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var existingTemplate models.ShiftTemplate
	var updatedTemplate models.ShiftTemplate

	// Lade existierende Vorlage
	if result := h.db.First(&existingTemplate, params["id"]); result.Error != nil {
		SendErrorResponse(w, http.StatusNotFound, "Schichtvorlage nicht gefunden")
		log.Printf("UpdateShiftTemplate Error: ID %s nicht gefunden", params["id"])
		return
	}

	// Dekodiere Update-Daten
	if err := json.NewDecoder(r.Body).Decode(&updatedTemplate); err != nil {
		SendErrorResponse(w, http.StatusBadRequest, "Ungültige Eingabedaten")
		log.Printf("UpdateShiftTemplate Decode Error: %v", err)
		return
	}

	// Aktualisiere die Felder
	result := h.db.Model(&existingTemplate).Updates(map[string]interface{}{
		"name":                    updatedTemplate.Name,
		"description":             updatedTemplate.Description,
		"color":                   updatedTemplate.Color,
		"employee_id":             updatedTemplate.EmployeeID,
		"monday_shift_type_id":    updatedTemplate.Monday.ShiftTypeID,
		"tuesday_shift_type_id":   updatedTemplate.Tuesday.ShiftTypeID,
		"wednesday_shift_type_id": updatedTemplate.Wednesday.ShiftTypeID,
		"thursday_shift_type_id":  updatedTemplate.Thursday.ShiftTypeID,
		"friday_shift_type_id":    updatedTemplate.Friday.ShiftTypeID,
		"saturday_shift_type_id":  updatedTemplate.Saturday.ShiftTypeID,
		"sunday_shift_type_id":    updatedTemplate.Sunday.ShiftTypeID,
	})

	if result.Error != nil {
		SendErrorResponse(w, http.StatusInternalServerError, "Fehler beim Aktualisieren der Schichtvorlage")
		log.Printf("UpdateShiftTemplate DB Error: %v", result.Error)
		return
	}

	// Lade aktualisierte Vorlage mit allen Beziehungen
	h.db.Preload("Employee").
		Preload("Monday.ShiftType").
		Preload("Tuesday.ShiftType").
		Preload("Wednesday.ShiftType").
		Preload("Thursday.ShiftType").
		Preload("Friday.ShiftType").
		Preload("Saturday.ShiftType").
		Preload("Sunday.ShiftType").
		First(&existingTemplate, existingTemplate.ID)

	SendSuccessResponse(w, "Schichtvorlage erfolgreich aktualisiert", existingTemplate)
	log.Printf("UpdateShiftTemplate: Schichtvorlage ID %d aktualisiert", existingTemplate.ID)
}

func (h *ShiftTemplateHandler) DeleteShiftTemplate(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var shiftTemplate models.ShiftTemplate

	if result := h.db.First(&shiftTemplate, params["id"]); result.Error != nil {
		SendErrorResponse(w, http.StatusNotFound, "Schichtvorlage nicht gefunden")
		log.Printf("DeleteShiftTemplate Error: ID %s nicht gefunden", params["id"])
		return
	}

	if result := h.db.Delete(&shiftTemplate); result.Error != nil {
		SendErrorResponse(w, http.StatusInternalServerError, "Fehler beim Löschen der Schichtvorlage")
		log.Printf("DeleteShiftTemplate DB Error: %v", result.Error)
		return
	}

	SendSuccessResponse(w, "Schichtvorlage erfolgreich gelöscht", nil)
	log.Printf("DeleteShiftTemplate: Schichtvorlage ID %d gelöscht", shiftTemplate.ID)
}
