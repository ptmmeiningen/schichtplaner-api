package models

import "time"

type Shift struct {
	ID          uint       `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at"`
	EmployeeID  uint       `gorm:"not null" json:"employee_id"`
	Employee    Employee   `gorm:"foreignKey:EmployeeID" json:"employee"`
	ShiftTypeID uint       `gorm:"not null" json:"shift_type_id"`
	ShiftType   ShiftType  `gorm:"foreignKey:ShiftTypeID" json:"shift_type"`
	StartTime   time.Time  `gorm:"not null" json:"start_time"`
	EndTime     time.Time  `gorm:"not null" json:"end_time"`
	Description string     `json:"description"`
}
