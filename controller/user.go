package controller

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/marioheryanto/MIG/helper"
	"github.com/marioheryanto/MIG/library"
	"github.com/marioheryanto/MIG/model"
)

type UserController struct {
	Library library.UserLibraryInterface
}

type UserControllerInterface interface {
	Register(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	Home(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
}

func NewUserController(library library.UserLibraryInterface) UserControllerInterface {
	return UserController{
		Library: library,
	}
}

func (controller UserController) Register(c *fiber.Ctx) error {
	request := model.UserRequest{}
	response := model.Response{}

	err := c.BodyParser(&request)
	if err != nil {
		return c.JSON(helper.GenerateResponse(c, &response, err))
	}

	err = helper.ValidateRegister(request)
	if err != nil {
		return c.JSON(helper.GenerateResponse(c, &response, err))
	}

	err = controller.Library.Register(request)
	if err != nil {
		return c.JSON(helper.GenerateResponse(c, &response, err))
	}

	c.Status(http.StatusCreated)
	response.Data = "Register success"
	return c.JSON(response)
}

func (controller UserController) Login(c *fiber.Ctx) error {
	request := model.UserRequest{}
	response := model.Response{}

	err := c.BodyParser(&request)
	if err != nil {
		return c.JSON(helper.GenerateResponse(c, &response, err))
	}

	err = helper.ValidateLogin(request)
	if err != nil {
		return c.JSON(helper.GenerateResponse(c, &response, err))
	}

	cookie, err := controller.Library.Login(request)
	if err != nil {
		return c.JSON(helper.GenerateResponse(c, &response, err))
	}

	c.Cookie(&cookie)

	c.Status(http.StatusOK)
	response.Data = "logged in"
	return c.JSON(response)
}

func (controller UserController) Home(c *fiber.Ctx) error {
	response := model.Response{}
	tokenString := c.Cookies("jwt")

	if tokenString == "" {
		c.Status(http.StatusUnauthorized)
		response.Message = "please login first"
		return c.JSON(response)
	}

	user, err := controller.Library.Home(tokenString)
	if err != nil {
		return c.JSON(helper.GenerateResponse(c, &response, err))
	}

	c.Status(http.StatusOK)
	response.Data = user
	return c.JSON(response)
}

func (controller UserController) Logout(c *fiber.Ctx) error {
	response := model.Response{}

	tokenString := c.Cookies("jwt")

	if tokenString == "" {
		c.Status(http.StatusUnauthorized)
		response.Message = "please login first"
		return c.JSON(response)
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().In(helper.LoadLocationJakarta()).Add(-5 * time.Minute),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	c.Status(http.StatusOK)
	response.Data = "logged out"
	return c.JSON(response)
}
