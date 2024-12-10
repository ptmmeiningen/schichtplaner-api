package models

import "gorm.io/gorm"

func (s *Shift) BeforeDelete(tx *gorm.DB) error {
	// Keine speziellen Abhängigkeiten zu berücksichtigen
	return nil
}

func (s *Shift) BeforeSave(tx *gorm.DB) error {
	// Validiere ShiftType Existenz
	var shiftType ShiftType
	if result := tx.First(&shiftType, s.ShiftTypeID); result.Error != nil {
		return result.Error
	}

	// Validiere Employee Existenz
	var employee Employee
	if result := tx.First(&employee, s.EmployeeID); result.Error != nil {
		return result.Error
	}

	// Validiere Zeitraum
	if s.EndTime.Before(s.StartTime) {
		return gorm.ErrInvalidField
	}

	return nil
}
