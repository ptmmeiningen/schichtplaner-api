package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"shift-planner/api/internal/models"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type EmployeeHandler struct {
	db *gorm.DB
}

func NewEmployeeHandler(db *gorm.DB) *EmployeeHandler {
	return &EmployeeHandler{db: db}
}

func (h *EmployeeHandler) GetEmployees(w http.ResponseWriter, r *http.Request) {
	var employees []models.Employee
	result := h.db.Preload("Departments.Employees").
		Order("first_name ASC, last_name ASC").
		Find(&employees)

	if result.Error != nil {
		SendErrorResponse(w, http.StatusInternalServerError, "Fehler beim Abrufen der Mitarbeiter")
		log.Printf("GetEmployees DB Error: %v", result.Error)
		return
	}

	SendSuccessResponse(w, "Mitarbeiter erfolgreich abgerufen", employees)
	log.Printf("GetEmployees: %d Mitarbeiter abgerufen", len(employees))
}

func (h *EmployeeHandler) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	var employeeInput struct {
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		Email       string `json:"email"`
		Password    string `json:"password"`
		Color       string `json:"color"`
		Departments []uint `json:"departments"`
	}

	if err := json.NewDecoder(r.Body).Decode(&employeeInput); err != nil {
		SendErrorResponse(w, http.StatusBadRequest, "Ungültige Eingabedaten")
		log.Printf("CreateEmployee Decode Error: %v", err)
		return
	}

	employee := models.Employee{
		FirstName: employeeInput.FirstName,
		LastName:  employeeInput.LastName,
		Email:     employeeInput.Email,
		Password:  employeeInput.Password,
		Color:     employeeInput.Color,
	}

	tx := h.db.Begin()

	if err := tx.Create(&employee).Error; err != nil {
		tx.Rollback()
		SendErrorResponse(w, http.StatusInternalServerError, "Fehler beim Erstellen des Mitarbeiters")
		log.Printf("CreateEmployee DB Error: %v", err)
		return
	}

	if len(employeeInput.Departments) > 0 {
		var departments []models.Department
		if err := tx.Find(&departments, employeeInput.Departments).Error; err != nil {
			tx.Rollback()
			SendErrorResponse(w, http.StatusBadRequest, "Abteilungen nicht gefunden")
			return
		}
		if err := tx.Model(&employee).Association("Departments").Replace(departments); err != nil {
			tx.Rollback()
			SendErrorResponse(w, http.StatusInternalServerError, "Fehler beim Zuweisen der Abteilungen")
			return
		}
	}

	tx.Commit()

	h.db.Preload("Departments.Employees").First(&employee, employee.ID)

	SendSuccessResponse(w, "Mitarbeiter erfolgreich erstellt", employee)
	log.Printf("CreateEmployee: Neuer Mitarbeiter ID %d erstellt", employee.ID)
}

func (h *EmployeeHandler) GetEmployee(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var employee models.Employee

	result := h.db.Preload("Departments.Employees").First(&employee, params["id"])
	if result.Error != nil {
		SendErrorResponse(w, http.StatusNotFound, "Mitarbeiter nicht gefunden")
		log.Printf("GetEmployee Error: ID %s nicht gefunden", params["id"])
		return
	}

	SendSuccessResponse(w, "Mitarbeiter erfolgreich abgerufen", employee)
	log.Printf("GetEmployee: Mitarbeiter ID %d abgerufen", employee.ID)
}

func (h *EmployeeHandler) UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var employeeInput struct {
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		Email       string `json:"email"`
		Password    string `json:"password"`
		Color       string `json:"color"`
		Departments []uint `json:"departments"`
	}

	if err := json.NewDecoder(r.Body).Decode(&employeeInput); err != nil {
		SendErrorResponse(w, http.StatusBadRequest, "Ungültige Eingabedaten")
		log.Printf("UpdateEmployee Decode Error: %v", err)
		return
	}

	var employee models.Employee
	if result := h.db.First(&employee, params["id"]); result.Error != nil {
		SendErrorResponse(w, http.StatusNotFound, "Mitarbeiter nicht gefunden")
		log.Printf("UpdateEmployee Error: ID %s nicht gefunden", params["id"])
		return
	}

	tx := h.db.Begin()

	employee.FirstName = employeeInput.FirstName
	employee.LastName = employeeInput.LastName
	employee.Email = employeeInput.Email
	if employeeInput.Password != "" {
		employee.Password = employeeInput.Password
	}
	employee.Color = employeeInput.Color

	if err := tx.Save(&employee).Error; err != nil {
		tx.Rollback()
		SendErrorResponse(w, http.StatusInternalServerError, "Fehler beim Aktualisieren des Mitarbeiters")
		log.Printf("UpdateEmployee DB Error: %v", err)
		return
	}

	if len(employeeInput.Departments) > 0 {
		var departments []models.Department
		if err := tx.Find(&departments, employeeInput.Departments).Error; err != nil {
			tx.Rollback()
			SendErrorResponse(w, http.StatusBadRequest, "Abteilungen nicht gefunden")
			return
		}
		if err := tx.Model(&employee).Association("Departments").Replace(departments); err != nil {
			tx.Rollback()
			SendErrorResponse(w, http.StatusInternalServerError, "Fehler beim Aktualisieren der Abteilungen")
			return
		}
	}

	tx.Commit()

	h.db.Preload("Departments.Employees").First(&employee, employee.ID)

	SendSuccessResponse(w, "Mitarbeiter erfolgreich aktualisiert", employee)
	log.Printf("UpdateEmployee: Mitarbeiter ID %d aktualisiert", employee.ID)
}

func (h *EmployeeHandler) DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var employee models.Employee

	if result := h.db.First(&employee, params["id"]); result.Error != nil {
		SendErrorResponse(w, http.StatusNotFound, "Mitarbeiter nicht gefunden")
		log.Printf("DeleteEmployee Error: ID %s nicht gefunden", params["id"])
		return
	}

	tx := h.db.Begin()

	if err := tx.Model(&employee).Association("Departments").Clear(); err != nil {
		tx.Rollback()
		SendErrorResponse(w, http.StatusInternalServerError, "Fehler beim Löschen der Abteilungszuweisungen")
		return
	}

	if err := tx.Delete(&employee).Error; err != nil {
		tx.Rollback()
		SendErrorResponse(w, http.StatusInternalServerError, "Fehler beim Löschen des Mitarbeiters")
		log.Printf("DeleteEmployee DB Error: %v", err)
		return
	}

	tx.Commit()

	SendSuccessResponse(w, "Mitarbeiter erfolgreich gelöscht", nil)
	log.Printf("DeleteEmployee: Mitarbeiter ID %d gelöscht", employee.ID)
}
