package models

import (
	"math/rand"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedDatabase(db *gorm.DB) {
	// Admin
	var adminCount int64
	db.Model(&Admin{}).Count(&adminCount)

	if adminCount == 0 {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		admin := Admin{
			Username:    "admin",
			Password:    string(hashedPassword),
			Email:       "admin@example.com",
			FirstName:   "Admin",
			LastName:    "User",
			IsActive:    true,
			IsSuperUser: true,
		}
		db.Create(&admin)
	}

	// Departments
	itDepartment := Department{
		Name:        "IT",
		Description: "Informationstechnologie Abteilung",
		Color:       "#3b82f6",
	}
	db.Create(&itDepartment)

	hrDepartment := Department{
		Name:        "HR",
		Description: "Human Resources Abteilung",
		Color:       "#22c55e",
	}
	db.Create(&hrDepartment)

	marketingDepartment := Department{
		Name:        "Marketing",
		Description: "Marketing Abteilung",
		Color:       "#f43f5e",
	}
	db.Create(&marketingDepartment)

	salesDepartment := Department{
		Name:        "Vertrieb",
		Description: "Vertriebsabteilung",
		Color:       "#a855f7",
	}
	db.Create(&salesDepartment)

	// ShiftTypes
	früh := ShiftType{
		Name:        "Früh",
		Description: "Frühschicht",
		StartTime:   "06:00",
		EndTime:     "14:00",
		Color:       "#0ea5e9",
	}
	db.Create(&früh)

	spät := ShiftType{
		Name:        "Spät",
		Description: "Spätschicht",
		StartTime:   "14:00",
		EndTime:     "22:00",
		Color:       "#6366f1",
	}
	db.Create(&spät)

	nacht := ShiftType{
		Name:        "Nacht",
		Description: "Nachtschicht",
		StartTime:   "22:00",
		EndTime:     "06:00",
		Color:       "#8b5cf6",
	}
	db.Create(&nacht)

	bereitschaft := ShiftType{
		Name:        "Bereitschaft",
		Description: "Bereitschaftsdienst",
		StartTime:   "00:00",
		EndTime:     "24:00",
		Color:       "#84cc16",
	}
	db.Create(&bereitschaft)

	teilzeit := ShiftType{
		Name:        "Teilzeit",
		Description: "Teilzeitschicht",
		StartTime:   "09:00",
		EndTime:     "15:00",
		Color:       "#f59e0b",
	}
	db.Create(&teilzeit)

	// Mitarbeiter
	employeeNames := []struct {
		FirstName string
		LastName  string
		Color     string
	}{
		{"Anna", "Schmidt", "#ef4444"},
		{"Thomas", "Weber", "#f97316"},
		{"Sarah", "Meyer", "#84cc16"},
		{"Michael", "Wagner", "#06b6d4"},
		{"Laura", "Fischer", "#8b5cf6"},
		{"Felix", "Koch", "#ec4899"},
		{"Julia", "Becker", "#f43f5e"},
		{"David", "Hoffmann", "#10b981"},
		{"Lisa", "Schulz", "#6366f1"},
		{"Jonas", "Richter", "#14b8a6"},
		{"Nina", "Wolf", "#f59e0b"},
		{"Tim", "Schäfer", "#3b82f6"},
		{"Lena", "Bauer", "#a855f7"},
		{"Paul", "Klein", "#d946ef"},
		{"Marie", "Krause", "#0ea5e9"},
		{"Lukas", "Schwarz", "#22c55e"},
		{"Sophie", "Schneider", "#be123c"},
		{"Jan", "Zimmermann", "#7c3aed"},
		{"Emma", "König", "#0d9488"},
		{"Finn", "Lang", "#b91c1c"},
		{"Hannah", "Schmitt", "#c2410c"},
		{"Leon", "Werner", "#15803d"},
		{"Mia", "Peters", "#1d4ed8"},
		{"Ben", "Neumann", "#7e22ce"},
		{"Clara", "Schmitz", "#be185d"},
		{"Noah", "Krüger", "#115e59"},
		{"Lea", "Friedrich", "#854d0e"},
		{"Luis", "Scholz", "#1e40af"},
		{"Sophia", "Möller", "#86198f"},
		{"Max", "Hartmann", "#991b1b"},
	}

	departments := []Department{itDepartment, hrDepartment, marketingDepartment, salesDepartment}

	var employees []Employee
	for _, name := range employeeNames {
		employee := Employee{
			FirstName:   name.FirstName,
			LastName:    name.LastName,
			Email:       strings.ToLower(name.FirstName + "." + name.LastName + "@example.com"),
			Password:    "password123",
			Color:       name.Color,
			Departments: []Department{},
		}

		numDepts := rand.Intn(3) + 1
		selectedDepts := make(map[int]bool)
		for i := 0; i < numDepts; i++ {
			for {
				deptIndex := rand.Intn(len(departments))
				if !selectedDepts[deptIndex] {
					selectedDepts[deptIndex] = true
					employee.Departments = append(employee.Departments, departments[deptIndex])
					break
				}
			}
		}

		db.Create(&employee)
		employees = append(employees, employee)
	}

	// Schichtvorlagen
	shiftTemplates := []struct {
		Name        string
		Description string
		Color       string
		Monday      uint
		Tuesday     uint
		Wednesday   uint
		Thursday    uint
		Friday      uint
		Saturday    uint
		Sunday      uint
	}{
		{
			Name:        "Frühschicht",
			Description: "Nur Frühschichten",
			Color:       "#0ea5e9",
			Monday:      früh.ID,
			Tuesday:     früh.ID,
			Wednesday:   früh.ID,
			Thursday:    früh.ID,
			Friday:      früh.ID,
		},
		{
			Name:        "Spätschicht",
			Description: "Nur Spätschichten",
			Color:       "#6366f1",
			Monday:      spät.ID,
			Tuesday:     spät.ID,
			Wednesday:   spät.ID,
			Thursday:    spät.ID,
			Friday:      spät.ID,
		},
		{
			Name:        "Nachtschicht",
			Description: "Nur Nachtschichten",
			Color:       "#8b5cf6",
			Monday:      nacht.ID,
			Tuesday:     nacht.ID,
			Wednesday:   nacht.ID,
			Thursday:    nacht.ID,
			Friday:      nacht.ID,
		},
		{
			Name:        "Früh/Spät Mix",
			Description: "Gemischte Schichten",
			Color:       "#22c55e",
			Monday:      früh.ID,
			Tuesday:     spät.ID,
			Wednesday:   früh.ID,
			Thursday:    spät.ID,
			Friday:      früh.ID,
		},
		{
			Name:        "Teilzeit",
			Description: "Teilzeitschichten",
			Color:       "#f59e0b",
			Monday:      teilzeit.ID,
			Tuesday:     teilzeit.ID,
			Wednesday:   teilzeit.ID,
		},
	}

	for _, template := range shiftTemplates {
		employeeID := employees[rand.Intn(len(employees))].ID
		db.Create(&ShiftTemplate{
			Name:        template.Name,
			Description: template.Description,
			Color:       template.Color,
			EmployeeID:  &employeeID,
			Monday:      ShiftDay{ShiftTypeID: template.Monday},
			Tuesday:     ShiftDay{ShiftTypeID: template.Tuesday},
			Wednesday:   ShiftDay{ShiftTypeID: template.Wednesday},
			Thursday:    ShiftDay{ShiftTypeID: template.Thursday},
			Friday:      ShiftDay{ShiftTypeID: template.Friday},
			Saturday:    ShiftDay{ShiftTypeID: template.Saturday},
			Sunday:      ShiftDay{ShiftTypeID: template.Sunday},
		})
	}
}
