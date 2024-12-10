package models

import "gorm.io/gorm"

func (e *Employee) BeforeDelete(tx *gorm.DB) error {
	// Lösche alle abhängigen ShiftTemplates
	if err := tx.Delete(&ShiftTemplate{}, "employee_id = ?", e.ID).Error; err != nil {
		return err
	}

	// Lösche alle abhängigen Shifts
	if err := tx.Delete(&Shift{}, "employee_id = ?", e.ID).Error; err != nil {
		return err
	}

	// Lösche alle Verknüpfungen zu Departments
	if err := tx.Table("employee_departments").Where("employee_id = ?", e.ID).Delete(&struct{}{}).Error; err != nil {
		return err
	}

	return nil
}

func (e *Employee) BeforeSave(tx *gorm.DB) error {
	// Validiere Pflichtfelder
	if e.FirstName == "" || e.LastName == "" || e.Email == "" {
		return gorm.ErrInvalidField
	}

	// Validiere Email-Eindeutigkeit
	var existingEmployee Employee
	if result := tx.Where("email = ? AND id != ?", e.Email, e.ID).First(&existingEmployee); result.Error == nil {
		return gorm.ErrDuplicatedKey
	}

	return nil
}
