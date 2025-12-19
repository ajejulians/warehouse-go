package app

import "github.com/gofiber/fiber/v2"

func SetupRoutes(app *fiber.App, container *Container) {
	api := app.Group("/api/v1")

	roles := api.Group("/roles")
	roles.Post("/", container.RoleController.CreateRole)
	roles.Get("/", container.RoleController.GetAllRoles)
	roles.Get("/:id", container.RoleController.GetAllRoleByID)
	roles.Put("/:id", container.RoleController.UpdateRole)
	roles.Delete("/:id", container.RoleController.DeleteRole)

	users := api.Group("/users")
	users.Post("/", container.UserController.CreateUser)
	users.Get("/", container.UserController.GetAllUsers)
	users.Get("/:id", container.UserController.GetUserByID)
	users.Get("/email/:email", container.UserController.GetUserByID)
	users.Put("/:id", container.UserController.UpdateUser)
	users.Delete("/:id", container.UserController.DeleteUser)

	assignRole := api.Group("/assign-role")
	assignRole.Post("/", container.UserController.AssignUserToRole)
	assignRole.Get("/", container.UserController.GetAllUserRoles)
	assignRole.Get("/:userRoleID", container.UserController.GetUserRoleByID)
	assignRole.Put("/:userRoleID", container.UserController.EditAssignUserToRole)

	users.Get("/role/:roleName", container.UserController.GetUserByRoleName)

	auth := api.Group("/auth")
	auth.Post("/login", container.AuthController.Login)

	upload := api.Group("/upload")
	upload.Post("/photo", container.UploadController.UploadPhoto)

}