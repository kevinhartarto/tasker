package models

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	TaskId      uuid.UUID `json:"task_id"`
	Task        string    `json:"task"`
	Description string    `json:"description"`
	Finished    bool      `json:"finished"`
	CreatedAt   time.Time `json:"created"`
	UpdatedAt   time.Time `json:"updated"`
}

type Reminder struct {
	ReminderId        uuid.UUID `json:"reminder_id"`
	TaskId            uuid.UUID `json:"task_id"`
	Reminder          string    `json:"reminder"`
	Description       string    `json:"description"`
	StartTime         time.Time `json:"start_time"`
	Frequency         string    `json:"frequency"`
	RepeatDays        []string  `json:"repeat_days" gorm:"type:text"`
	RepeatSameday     bool      `json:"repeat_sameday"`
	RepeatUntil       time.Time `json:"repeat_until"`
	Interval          int       `json:"inteval"`
	IntervalInMinutes int       `json:"interval_in_minutes"`
	CreatedAt         time.Time `json:"created"`
	UpdatedAt         time.Time `json:"updated"`
}
