package controllers

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/kevinhartarto/mytodolist/internal/database"
	"github.com/kevinhartarto/mytodolist/internal/models"
)

type DayController interface {

	// Get all days
	// return all days
	GetAllDays(c *fiber.Ctx) error

	// Get day by id or name
	// return day
	GetDay(c *fiber.Ctx) error
}

var (
	dayInstance *dayController
	days        []models.Day
	day         models.Day
)

type dayController struct {
	db database.Service
}

func NewDayController(db database.Service) *dayController {
	if dayInstance != nil {
		return dayInstance
	}

	dayInstance = &dayController{
		db: db,
	}

	return dayInstance
}

func (dc *dayController) GetAllDays(c *fiber.Ctx) error {
	dc.db.UseGorm().Find(&days)
	result, _ := json.Marshal(days)
	return c.SendString(string(result))
}

func (dc *dayController) GetDay(c *fiber.Ctx) error {
	if err := c.BodyParser(&day); err != nil {
		return err
	}

	if err := dc.db.UseGorm().First(&day).Error; err != nil {
		return err
	}

	result, _ := json.Marshal(&day)
	return c.SendString(string(result))
}
