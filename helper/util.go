package helper

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/marioheryanto/MIG/model"
)

func NewServiceError(code int, message string) model.ServiceError {
	return model.ServiceError{
		Code:    code,
		Message: message,
	}
}

func ReplaceServiceErrorForLogin(err error) error {
	serviceErr, ok := err.(model.ServiceError)
	if serviceErr.Code == http.StatusNotFound || !ok {
		serviceErr.Code = http.StatusBadRequest
		serviceErr.Message = "email atau password salah"
	}

	return serviceErr
}

func GenerateResponse(c *fiber.Ctx, r *model.Response, err error) model.Response {
	//safety assert
	serviceErr, ok := err.(model.ServiceError)
	if !ok {
		c.Status(http.StatusInternalServerError)
		r.Message = err.Error()
		return *r
	}

	code := http.StatusInternalServerError
	if serviceErr.Code != 0 {
		code = serviceErr.Code
	}

	c.Status(code)
	r.Message = serviceErr.Error()

	return *r
}
