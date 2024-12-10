package models

import "time"

type ShiftType struct {
	ID          uint       `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at"`
	Name        string     `gorm:"not null;unique" json:"name"`
	Description string     `json:"description"`
	StartTime   string     `gorm:"not null" json:"start_time"`
	EndTime     string     `gorm:"not null" json:"end_time"`
	Color       string     `gorm:"not null" json:"color"`
	Shifts      []Shift    `json:"shifts"`
}
