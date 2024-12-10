package models

import "time"

type Department struct {
	ID          uint       `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at"`
	Name        string     `gorm:"not null;unique" json:"name"`
	Description string     `json:"description"`
	Color       string     `gorm:"not null" json:"color"`
	Employees   []Employee `gorm:"many2many:employee_departments;" json:"employees"`
}
