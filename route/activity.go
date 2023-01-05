package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/marioheryanto/MIG/controller"
)

func ActivityRoutes(app *fiber.App, c controller.ActivityControllerInterface) {
	activity := app.Group("/activity")
	activity.Post("/create", c.Create)
	activity.Put("/edit", c.Edit)
	activity.Delete("/delete", c.Delete)
	activity.Get("/all", c.GetAll)
}
