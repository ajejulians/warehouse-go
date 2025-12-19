package controller

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"warehouse-go/api-gateaway/middleware"

	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	userServiceURL string
	jwtConfig      middleware.JWTConfig
}

type LoginRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	UserID uint `json:"user_id"`
	Email string `json:"email"`
	Roles string `json:"roles"`
}

type UserServiceResponse struct {
	Data struct {
		UserID uint `json:"user_id"`
		Email string `json:"email"`
		Role []string `json:"roles"`
	} `json:"data"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User struct {
		ID uint `json:"id"`
		Email string `json:"email"`
		Roles string `json:"roles"`
	} `json:"user"`
}

func NewAuthController(userServiceURL string, jwtConfig middleware.JWTConfig) *AuthController {
	return &AuthController{
		userServiceURL: userServiceURL,
		jwtConfig: jwtConfig,
	}
}

func (a *AuthController) Login(c *fiber.Ctx) error {
	var loginRequest LoginRequest
	if err := c.BodyParser(&loginRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message" : err.Error(),
		})
	}

	if loginRequest.Email == "" || loginRequest.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error" : "Bad Request",
			"message" : "Email and password are required",
 		})
	}

	loginResp, err := a.forwardLoginRequest(loginRequest)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message" : "Internal Server Error",
		})
	}

	token, err := middleware.GenerateJWT(loginResp.UserID, loginResp.Email, loginResp.Roles, a.jwtConfig)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message" : err.Error(),
		})
	}

	response := AuthResponse{
		Token: token,
		User: struct{
			ID uint `json:"id"`
			Email string `json:"email"`
			Roles string `json:"roles"`
		}{
			ID: loginResp.UserID,
			Email: loginResp.Email,
			Roles: loginResp.Roles,
		},
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message" : "Login successful",
		"data" : response,
	})
}

func (ac *AuthController) forwardLoginRequest(loginReq LoginRequest) (*LoginResponse, error) {
	reqBody, err := json.Marshal(loginReq)
	if err != nil {
		return nil, err 
	}

	req, err :=http.NewRequest("POST", ac.userServiceURL+"api/v1/auth/login", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Gateaway", "warehouse-api-gateaway")
	req.Header.Set("X-Internal-Request", "true")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, err
	}

	var userServiceResp UserServiceResponse
	if err := json.Unmarshal(respBody, &userServiceResp); err != nil {
		return nil, err
	}

	roleStr := ""
	if len(userServiceResp.Data.Role) > 0 {
		roleStr = userServiceResp.Data.Role[0]
		for i := 1; i < len(userServiceResp.Data.Role); i++ {
			roleStr += "," + userServiceResp.Data.Role[i]
		}
	}

	loginResp := LoginResponse{
		UserID: userServiceResp.Data.UserID,
		Email: userServiceResp.Data.Email,
		Roles: roleStr,
	}

	return &loginResp, nil
}