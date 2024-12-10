package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"shift-planner/api/internal/models"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type DepartmentHandler struct {
	db *gorm.DB
}

func NewDepartmentHandler(db *gorm.DB) *DepartmentHandler {
	return &DepartmentHandler{db: db}
}

func (h *DepartmentHandler) GetDepartments(w http.ResponseWriter, r *http.Request) {
	var departments []models.Department
	result := h.db.Preload("Employees.Departments").
		Order("name ASC").
		Find(&departments)

	if result.Error != nil {
		SendErrorResponse(w, http.StatusInternalServerError, "Fehler beim Abrufen der Abteilungen")
		log.Printf("GetDepartments DB Error: %v", result.Error)
		return
	}

	SendSuccessResponse(w, "Abteilungen erfolgreich abgerufen", departments)
	log.Printf("GetDepartments: %d Abteilungen abgerufen", len(departments))
}

func (h *DepartmentHandler) CreateDepartment(w http.ResponseWriter, r *http.Request) {
	var department models.Department
	if err := json.NewDecoder(r.Body).Decode(&department); err != nil {
		SendErrorResponse(w, http.StatusBadRequest, "Ungültige Eingabedaten")
		log.Printf("CreateDepartment Decode Error: %v", err)
		return
	}

	result := h.db.Create(&department)
	if result.Error != nil {
		SendErrorResponse(w, http.StatusInternalServerError, "Fehler beim Erstellen der Abteilung")
		log.Printf("CreateDepartment DB Error: %v", result.Error)
		return
	}

	h.db.Preload("Employees.Departments").First(&department, department.ID)

	SendSuccessResponse(w, "Abteilung erfolgreich erstellt", department)
	log.Printf("CreateDepartment: Neue Abteilung ID %d erstellt", department.ID)
}

func (h *DepartmentHandler) GetDepartment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var department models.Department

	result := h.db.Preload("Employees.Departments").First(&department, params["id"])
	if result.Error != nil {
		SendErrorResponse(w, http.StatusNotFound, "Abteilung nicht gefunden")
		log.Printf("GetDepartment Error: ID %s nicht gefunden", params["id"])
		return
	}

	SendSuccessResponse(w, "Abteilung erfolgreich abgerufen", department)
	log.Printf("GetDepartment: Abteilung ID %d abgerufen", department.ID)
}

func (h *DepartmentHandler) UpdateDepartment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var department models.Department

	if result := h.db.First(&department, params["id"]); result.Error != nil {
		SendErrorResponse(w, http.StatusNotFound, "Abteilung nicht gefunden")
		log.Printf("UpdateDepartment Error: ID %s nicht gefunden", params["id"])
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&department); err != nil {
		SendErrorResponse(w, http.StatusBadRequest, "Ungültige Eingabedaten")
		log.Printf("UpdateDepartment Decode Error: %v", err)
		return
	}

	if result := h.db.Save(&department); result.Error != nil {
		SendErrorResponse(w, http.StatusInternalServerError, "Fehler beim Aktualisieren der Abteilung")
		log.Printf("UpdateDepartment DB Error: %v", result.Error)
		return
	}

	h.db.Preload("Employees.Departments").First(&department, department.ID)

	SendSuccessResponse(w, "Abteilung erfolgreich aktualisiert", department)
	log.Printf("UpdateDepartment: Abteilung ID %d aktualisiert", department.ID)
}

func (h *DepartmentHandler) DeleteDepartment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var department models.Department

	if result := h.db.First(&department, params["id"]); result.Error != nil {
		SendErrorResponse(w, http.StatusNotFound, "Abteilung nicht gefunden")
		log.Printf("DeleteDepartment Error: ID %s nicht gefunden", params["id"])
		return
	}

	tx := h.db.Begin()

	if err := tx.Model(&department).Association("Employees").Clear(); err != nil {
		tx.Rollback()
		SendErrorResponse(w, http.StatusInternalServerError, "Fehler beim Löschen der Mitarbeiterzuweisungen")
		return
	}

	if result := tx.Delete(&department); result.Error != nil {
		tx.Rollback()
		SendErrorResponse(w, http.StatusInternalServerError, "Fehler beim Löschen der Abteilung")
		log.Printf("DeleteDepartment DB Error: %v", result.Error)
		return
	}

	tx.Commit()

	SendSuccessResponse(w, "Abteilung erfolgreich gelöscht", nil)
	log.Printf("DeleteDepartment: Abteilung ID %d gelöscht", department.ID)
}
