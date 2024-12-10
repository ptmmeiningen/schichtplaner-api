package models

import "time"

type Employee struct {
	ID             uint            `gorm:"primarykey" json:"id"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
	DeletedAt      *time.Time      `gorm:"index" json:"deleted_at"`
	FirstName      string          `gorm:"not null" json:"first_name"`
	LastName       string          `gorm:"not null" json:"last_name"`
	Email          string          `gorm:"unique;not null" json:"email"`
	Password       string          `gorm:"not null" json:"password"`
	Color          string          `gorm:"not null" json:"color"`
	Departments    []Department    `gorm:"many2many:employee_departments;" json:"departments"`
	Shifts         []Shift         `json:"shifts"`
	ShiftTemplates []ShiftTemplate `json:"shift_templates"`
}
