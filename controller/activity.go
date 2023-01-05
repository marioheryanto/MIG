package controller

import (
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/marioheryanto/MIG/helper"
	"github.com/marioheryanto/MIG/library"
	"github.com/marioheryanto/MIG/model"
)

type ActivityController struct {
	Library library.ActivityLibraryInterface
}

type ActivityControllerInterface interface {
	Create(c *fiber.Ctx) error
	Edit(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	GetAll(c *fiber.Ctx) error
}

func NewActivityController(library library.ActivityLibraryInterface) ActivityControllerInterface {
	return ActivityController{
		Library: library,
	}
}

func (controller ActivityController) Create(c *fiber.Ctx) error {
	request := model.Activity{}
	response := model.Response{}
	tokenString := c.Cookies("jwt")

	if tokenString == "" {
		c.Status(http.StatusUnauthorized)
		response.Message = "please login first"
		return c.JSON(response)
	}

	err := c.BodyParser(&request)
	if err != nil {
		return c.JSON(helper.GenerateResponse(c, &response, err))
	}

	err = helper.ValidateActivity(request, false)
	if err != nil {
		return c.JSON(helper.GenerateResponse(c, &response, err))
	}

	err = controller.Library.Create(tokenString, request)
	if err != nil {
		return c.JSON(helper.GenerateResponse(c, &response, err))
	}

	c.Status(http.StatusCreated)
	response.Data = "success create activity"
	return c.JSON(response)
}

func (controller ActivityController) Edit(c *fiber.Ctx) error {
	request := model.Activity{}
	response := model.Response{}
	tokenString := c.Cookies("jwt")

	if tokenString == "" {
		c.Status(http.StatusUnauthorized)
		response.Message = "please login first"
		return c.JSON(response)
	}

	err := c.BodyParser(&request)
	if err != nil {
		return c.JSON(helper.GenerateResponse(c, &response, err))
	}

	err = helper.ValidateActivity(request, true)
	if err != nil {
		return c.JSON(helper.GenerateResponse(c, &response, err))
	}

	err = controller.Library.Edit(tokenString, request)
	if err != nil {
		return c.JSON(helper.GenerateResponse(c, &response, err))
	}

	c.Status(http.StatusOK)
	response.Data = "success edit activity"
	return c.JSON(response)
}

func (controller ActivityController) Delete(c *fiber.Ctx) error {
	response := model.Response{}
	tokenString := c.Cookies("jwt")

	if tokenString == "" {
		c.Status(http.StatusUnauthorized)
		response.Message = "please login first"
		return c.JSON(response)
	}

	activityId := strings.TrimSpace(c.Query("id"))
	if activityId == "" {
		c.Status(http.StatusBadRequest)
		response.Message = "activity id is empty"
		return c.JSON(response)
	}

	err := controller.Library.Delete(tokenString, activityId)
	if err != nil {
		return c.JSON(helper.GenerateResponse(c, &response, err))
	}

	c.Status(http.StatusOK)
	response.Data = "success delete activity"
	return c.JSON(response)
}

func (controller ActivityController) GetAll(c *fiber.Ctx) error {
	response := model.Response{}
	tokenString := c.Cookies("jwt")

	if tokenString == "" {
		c.Status(http.StatusUnauthorized)
		response.Message = "please login first"
		return c.JSON(response)
	}

	from := strings.TrimSpace(c.Query("from"))
	to := strings.TrimSpace(c.Query("to"))

	if from == "" || to == "" {
		c.Status(http.StatusBadRequest)
		response.Message = "range date is empty"
		return c.JSON(response)
	}

	data, err := controller.Library.GetAll(tokenString, from, to)
	if err != nil {
		return c.JSON(helper.GenerateResponse(c, &response, err))
	}

	c.Status(http.StatusOK)
	response.Data = data
	return c.JSON(response)
}
