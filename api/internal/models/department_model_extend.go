package models

import "gorm.io/gorm"

func (d *Department) BeforeDelete(tx *gorm.DB) error {
	// Lösche alle Verknüpfungen zu Mitarbeitern in der Verbindungstabelle
	if err := tx.Table("employee_departments").Where("department_id = ?", d.ID).Delete(&struct{}{}).Error; err != nil {
		return err
	}
	return nil
}

func (d *Department) BeforeSave(tx *gorm.DB) error {
	// Validiere Pflichtfelder
	if d.Name == "" {
		return gorm.ErrInvalidField
	}

	// Validiere Name-Eindeutigkeit
	var existingDepartment Department
	if result := tx.Where("name = ? AND id != ?", d.Name, d.ID).First(&existingDepartment); result.Error == nil {
		return gorm.ErrDuplicatedKey
	}

	return nil
}
