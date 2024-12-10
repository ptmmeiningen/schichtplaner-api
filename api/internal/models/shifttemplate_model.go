package models

import (
	"time"
)

type ShiftTemplate struct {
	ID          uint       `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at"`
	Name        string     `gorm:"not null" json:"name"`
	EmployeeID  *uint      `json:"employee_id"`
	Employee    *Employee  `gorm:"foreignKey:EmployeeID" json:"employee"`
	Description string     `json:"description"`
	Color       string     `json:"color"`
	Monday      ShiftDay   `gorm:"embedded;embeddedPrefix:monday_" json:"monday"`
	Tuesday     ShiftDay   `gorm:"embedded;embeddedPrefix:tuesday_" json:"tuesday"`
	Wednesday   ShiftDay   `gorm:"embedded;embeddedPrefix:wednesday_" json:"wednesday"`
	Thursday    ShiftDay   `gorm:"embedded;embeddedPrefix:thursday_" json:"thursday"`
	Friday      ShiftDay   `gorm:"embedded;embeddedPrefix:friday_" json:"friday"`
	Saturday    ShiftDay   `gorm:"embedded;embeddedPrefix:saturday_" json:"saturday"`
	Sunday      ShiftDay   `gorm:"embedded;embeddedPrefix:sunday_" json:"sunday"`
}

type ShiftDay struct {
	ShiftTypeID uint      `gorm:"column:shift_type_id" json:"shift_type_id"`
	ShiftType   ShiftType `gorm:"-" json:"-"`
}
