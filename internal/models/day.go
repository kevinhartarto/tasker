package models

import "github.com/google/uuid"

type Day struct {
	DayId   int    `json:"day_id" gorm:"primaryKey"`
	DayName string `json:"day_name"`
}

type DayGroup struct {
	DayId       int       `json:"day_id"`
	TaskGroupId uuid.UUID `json:"task_group_id"`
}
