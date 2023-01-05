package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/marioheryanto/MIG/controller"
	"github.com/marioheryanto/MIG/database"
	"github.com/marioheryanto/MIG/library"
	"github.com/marioheryanto/MIG/repository"
	"github.com/marioheryanto/MIG/route"
)

func main() {
	db := database.Init()

	userRepository := repository.NewUserRepository(db)
	absensiRepository := repository.NewAbsensiRepository(db)
	activityRepository := repository.NewActivityRepository(db)

	userLibrary := library.NewUserLibrary(userRepository)
	absensiLibrary := library.NewAbsensiLibrary(absensiRepository)
	activityLibrary := library.NewActivityLibrary(activityRepository)

	userController := controller.NewUserController(userLibrary)
	absensiController := controller.NewAbsensiController(absensiLibrary)
	activityController := controller.NewActivityController(activityLibrary)

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))

	route.UserRoutes(app, userController)
	route.AbsensiRoutes(app, absensiController)
	route.ActivityRoutes(app, activityController)

	app.Listen("0.0.0.0:$PORT")
}
