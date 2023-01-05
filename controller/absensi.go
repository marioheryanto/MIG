package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/marioheryanto/MIG/helper"
	"github.com/marioheryanto/MIG/library"
	"github.com/marioheryanto/MIG/model"
)

type AbsensiController struct {
	Library library.AbsensiLibraryInterface
}

type AbsensiControllerInterface interface {
	CheckInOut(c *fiber.Ctx) error
	GetAll(c *fiber.Ctx) error
}

func NewAbsensiController(library library.AbsensiLibraryInterface) AbsensiControllerInterface {
	return AbsensiController{
		Library: library,
	}
}

func (controller AbsensiController) CheckInOut(c *fiber.Ctx) error {
	response := model.Response{}
	tokenString := c.Cookies("jwt")

	if tokenString == "" {
		c.Status(http.StatusUnauthorized)
		response.Message = "please login first"
		return c.JSON(response)
	}

	check := strings.ToLower(strings.TrimSpace(c.Query("check")))
	if !strings.Contains("in|out", check) || check == "" {
		c.Status(http.StatusBadRequest)
		response.Message = "value must 'in' or 'out'"
		return c.JSON(response)
	}

	err := controller.Library.CheckInOut(tokenString, check)
	if err != nil {
		return c.JSON(helper.GenerateResponse(c, &response, err))
	}

	c.Status(http.StatusOK)
	response.Data = fmt.Sprintf("success to check %v", check)
	return c.JSON(response)
}

func (controller AbsensiController) GetAll(c *fiber.Ctx) error {
	response := model.Response{}
	tokenString := c.Cookies("jwt")

	if tokenString == "" {
		c.Status(http.StatusUnauthorized)
		response.Message = "please login first"
		return c.JSON(response)
	}

	month, _ := strconv.Atoi(c.Query("month"))
	year, _ := strconv.Atoi(c.Query("year"))

	data, err := controller.Library.GetAll(tokenString, month, year)
	if err != nil {
		return c.JSON(helper.GenerateResponse(c, &response, err))
	}

	c.Status(http.StatusOK)
	response.Data = data
	return c.JSON(response)
}
