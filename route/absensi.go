package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/marioheryanto/MIG/controller"
)

func AbsensiRoutes(app *fiber.App, c controller.AbsensiControllerInterface) {
	absensi := app.Group("/absensi")
	absensi.Post("/check", c.CheckInOut)
	absensi.Get("/all", c.GetAll)
}
