package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/marioheryanto/MIG/controller"
)

func UserRoutes(app *fiber.App, c controller.UserControllerInterface) {
	user := app.Group("/user")
	user.Post("/register", c.Register)
	user.Post("/login", c.Login)
	user.Get("/home", c.Home)
	user.Post("/logout", c.Logout)
}
