package models

import "gorm.io/gorm"

func (st *ShiftType) BeforeDelete(tx *gorm.DB) error {
	// Lösche alle abhängigen Schichten
	if err := tx.Delete(&Shift{}, "shift_type_id = ?", st.ID).Error; err != nil {
		return err
	}

	// Setze ShiftType-Referenzen in ShiftTemplates auf null
	if err := tx.Model(&ShiftTemplate{}).Where("monday_shift_type_id = ?", st.ID).Update("monday_shift_type_id", nil).Error; err != nil {
		return err
	}
	if err := tx.Model(&ShiftTemplate{}).Where("tuesday_shift_type_id = ?", st.ID).Update("tuesday_shift_type_id", nil).Error; err != nil {
		return err
	}
	if err := tx.Model(&ShiftTemplate{}).Where("wednesday_shift_type_id = ?", st.ID).Update("wednesday_shift_type_id", nil).Error; err != nil {
		return err
	}
	if err := tx.Model(&ShiftTemplate{}).Where("thursday_shift_type_id = ?", st.ID).Update("thursday_shift_type_id", nil).Error; err != nil {
		return err
	}
	if err := tx.Model(&ShiftTemplate{}).Where("friday_shift_type_id = ?", st.ID).Update("friday_shift_type_id", nil).Error; err != nil {
		return err
	}
	if err := tx.Model(&ShiftTemplate{}).Where("saturday_shift_type_id = ?", st.ID).Update("saturday_shift_type_id", nil).Error; err != nil {
		return err
	}
	if err := tx.Model(&ShiftTemplate{}).Where("sunday_shift_type_id = ?", st.ID).Update("sunday_shift_type_id", nil).Error; err != nil {
		return err
	}

	return nil
}

func (st *ShiftType) BeforeSave(tx *gorm.DB) error {
	// Validiere Name
	if st.Name == "" {
		return gorm.ErrInvalidField
	}
	return nil
}
