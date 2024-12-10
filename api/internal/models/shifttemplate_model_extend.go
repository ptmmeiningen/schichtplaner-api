package models

import "gorm.io/gorm"

func (st *ShiftTemplate) BeforeDelete(tx *gorm.DB) error {
	return nil
}

func (st *ShiftTemplate) BeforeSave(tx *gorm.DB) error {
	// Validiere Name
	if st.Name == "" {
		return gorm.ErrInvalidField
	}

	// Validiere Employee falls gesetzt
	if st.EmployeeID != nil {
		var employee Employee
		if result := tx.First(&employee, *st.EmployeeID); result.Error != nil {
			return result.Error
		}
	}

	// Validiere ShiftType IDs wenn gesetzt
	if st.Monday.ShiftTypeID != 0 {
		if err := validateShiftTypeID(tx, st.Monday.ShiftTypeID); err != nil {
			return err
		}
	}
	if st.Tuesday.ShiftTypeID != 0 {
		if err := validateShiftTypeID(tx, st.Tuesday.ShiftTypeID); err != nil {
			return err
		}
	}
	if st.Wednesday.ShiftTypeID != 0 {
		if err := validateShiftTypeID(tx, st.Wednesday.ShiftTypeID); err != nil {
			return err
		}
	}
	if st.Thursday.ShiftTypeID != 0 {
		if err := validateShiftTypeID(tx, st.Thursday.ShiftTypeID); err != nil {
			return err
		}
	}
	if st.Friday.ShiftTypeID != 0 {
		if err := validateShiftTypeID(tx, st.Friday.ShiftTypeID); err != nil {
			return err
		}
	}
	if st.Saturday.ShiftTypeID != 0 {
		if err := validateShiftTypeID(tx, st.Saturday.ShiftTypeID); err != nil {
			return err
		}
	}
	if st.Sunday.ShiftTypeID != 0 {
		if err := validateShiftTypeID(tx, st.Sunday.ShiftTypeID); err != nil {
			return err
		}
	}

	return nil
}

func validateShiftTypeID(tx *gorm.DB, id uint) error {
	var shiftType ShiftType
	if result := tx.First(&shiftType, id); result.Error != nil {
		return result.Error
	}
	return nil
}
