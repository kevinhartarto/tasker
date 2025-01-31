package models

import (
	"time"

	"github.com/google/uuid"
)

type TaskGroup struct {
	TaskGroupId uuid.UUID `json:"task_group_id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	TaskGroup   string    `json:"task_group"`
	Reminder    bool      `json:"reminder" gorm:"default:false"`
	Repeats     bool      `json:"repeats" gorm:"default:false"`
	Deadline    time.Time `json:"deadline"`
	NotifyTime  time.Time `json:"notify_time"`
	CreatedAt   time.Time `json:"created"`
	UpdatedAt   time.Time `json:"updated"`
	Deprecated  bool      `json:"deprecated" gorm:"default:false"`
}

type Task struct {
	TaskId     uuid.UUID `json:"task_id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Task       string    `json:"task"`
	CreatedAt  time.Time `json:"created"`
	UpdatedAt  time.Time `json:"updated"`
	Deprecated bool      `json:"deprecated" gorm:"default:false"`
}

type TaskGroupTask struct {
	TaskGroupId uuid.UUID `json:"task_group_id"`
	TaskId      uuid.UUID `json:"task_id"`
}
