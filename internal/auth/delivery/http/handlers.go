package http

import (
	"github.com/gofiber/fiber/v3"
	"url-shortener/internal/auth"
	"url-shortener/internal/models"
	"url-shortener/pkg/http"
)

type request struct {
	Email    string `json:"email" validate:"required,email" example:"example@example.com"`
	Password string `json:"password,omitempty" validate:"required,min=6" example:"very_strong_password"`
}

type authHandlers struct {
	authService auth.Service
}

func NewAuthHandlers(authService auth.Service) auth.Handlers {
	return &authHandlers{authService: authService}
}

// Login godoc
//
//	@Summary		User login
//	@Description	Authenticate user and generate access token
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			user	body		request						true	"User credentials"
//	@Success		200		{object}	http.Response[auth.Token]	"Successfully authenticated"
//	@Failure		400		{object}	http.Response[any]			"Invalid request payload"
//	@Failure		401		{object}	http.Response[any]			"Invalid credentials"
//	@Failure		500		{object}	http.Response[any]			"Internal server error"
//	@Router			/auth/login [post]
func (h *authHandlers) Login(ctx fiber.Ctx) error {
	var req request
	if err := ctx.Bind().JSON(&req); err != nil {
		return err
	}

	user := models.User{
		Email:    req.Email,
		Password: req.Password,
	}
	token, err := h.authService.Login(ctx.Context(), &user)
	if err != nil {
		return err
	}

	data := http.OK[*auth.Token](token)
	return ctx.JSON(data)
}

// Register godoc
//
//	@Summary		Register new user
//	@Description	Create a new user account and return authentication tokens
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			user	body		request						true	"User registration data"
//	@Success		201		{object}	http.Response[auth.Token]	"User successfully registered"
//	@Failure		400		{object}	http.Response[any]			"Invalid request payload"
//	@Failure		409		{object}	http.Response[any]			"User already exists"
//	@Failure		500		{object}	http.Response[any]			"Internal server error"
//	@Router			/auth/register [post]
func (h *authHandlers) Register(ctx fiber.Ctx) error {
	var req request
	if err := ctx.Bind().JSON(&req); err != nil {
		return err
	}

	user := models.User{
		Email:    req.Email,
		Password: req.Password,
	}
	token, err := h.authService.Register(ctx.Context(), &user)
	if err != nil {
		return err
	}

	data := http.OK[*auth.Token](token)
	return ctx.Status(fiber.StatusCreated).JSON(data)
}

// GetMe godoc
//
//	@Summary		Get current user info
//	@Description	Returns authenticated user's information
//	@Tags			auth
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	http.Response[models.User]	"Successfully retrieved user data"
//	@Failure		401	{object}	http.Response[any]			"Unauthorized - Missing or invalid token"
//	@Failure		500	{object}	http.Response[any]			"Internal server error"
//	@Router			/auth/me [get]
func (h *authHandlers) GetMe(ctx fiber.Ctx) error {
	authData := fiber.Locals[*auth.Data](ctx, "authData")

	data := http.OK[*models.User](authData.User)
	return ctx.JSON(data)
}

// RefreshToken godoc
//
//	@Summary	    Refresh token
//	@Description	Refresh token. Requires valid refresh token in Authorization header.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	http.Response[auth.Token]	"Returns new access and refresh tokens"
//	@Failure		401	{object}	http.Response[any]			"Unauthorized - Missing or invalid token"
//	@Failure		500	{object}	http.Response[any]			"Internal server error"
//	@Router			/auth/refresh-token [post]
func (h *authHandlers) RefreshToken(ctx fiber.Ctx) error {
	authData := fiber.Locals[*auth.Data](ctx, "authData")
	res, err := h.authService.RefreshToken(ctx.Context(), authData)
	if err != nil {
		return err
	}

	data := http.OK[*auth.Token](res)
	return ctx.JSON(data)
}
