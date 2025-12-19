package controller

import (
	"warehouse-go/user-service/controller/request"
	"warehouse-go/user-service/controller/response"
	"warehouse-go/user-service/pkg/conv"
	"warehouse-go/user-service/pkg/validator"
	"warehouse-go/user-service/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type AuthControllerInterface interface {
	Login(c *fiber.Ctx) error
}

type AuthController struct {
	AuthService usecase.UserUsecaseInterface
}

// Login implements AuthControllerInterface.
func (a *AuthController) Login(c *fiber.Ctx) error {
	ctx := c.Context()
	var loginRequest request.LoginRequest
	if err := c.BodyParser(&loginRequest); err != nil {
		log.Errorf("[AuthController] Login - 1: %v", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	if err := validator.Validate(loginRequest); err != nil {
		log.Errorf("[AuthController] Login - 2: %v", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	user, err := a.AuthService.GetUserByEmail(ctx, loginRequest.Email)
	if err != nil {
		log.Errorf("[AuthController] Login - 3: %v", err.Error())
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid email or password",
		})
	}

	if user == nil {
		log.Errorf("[AuthController] Login - 4: user not found")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User not found",
		})
	}

	isSame := conv.CheckPasswordHash(loginRequest.Password, user.Password)
	if !isSame {
		log.Errorf("[AuthController] Login - 5: %v", err.Error())
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid email or password",
		})
	}	 

	var roles []string 
	for _, r := range user.Roles {
		roles = append(roles, r.Name)
	}

	loginResp := response.LoginResponse{
		UserID:   uint(user.ID),
		Email:    user.Email,
		Role: roles,
	}


	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
		"data":    loginResp,
	})


}

func NewAuthController(authService usecase.UserUsecaseInterface) AuthControllerInterface {
	return &AuthController{
		AuthService: authService,
	}
}
