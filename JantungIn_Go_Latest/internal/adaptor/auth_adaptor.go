package adaptor

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"jantungin-api-server/internal/dto"
	"jantungin-api-server/internal/usecase"
	"jantungin-api-server/pkg/middleware"
	"jantungin-api-server/pkg/utils"
)

type AuthAdaptor struct {
	authUsecase usecase.AuthUsecase
	cfg         *utils.Config
}

func NewAuthAdaptor(authUsecase usecase.AuthUsecase, cfg *utils.Config) *AuthAdaptor {
	return &AuthAdaptor{
		authUsecase: authUsecase,
		cfg:         cfg,
	}
}

func (h *AuthAdaptor) Register(c *gin.Context) {
	var req dto.AuthRegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Format data tidak valid", err.Error())
		return
	}

	data, err := h.authUsecase.Register(c.Request.Context(), req)
	if err != nil {
		switch err.Error() {
		case "username wajib diisi",
			"username minimal 3 karakter",
			"password minimal 6 karakter",
			"format tanggal lahir tidak valid, gunakan YYYY-MM-DD":
			utils.BadRequestResponse(c, err.Error(), nil)
		case "username sudah terdaftar",
			"email sudah terdaftar":
			utils.ErrorResponse(c, http.StatusConflict, err.Error(), nil)
		default:
			utils.InternalServerErrorResponse(c, err.Error(), nil)
		}
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "User registered successfully", data)
}

func (h *AuthAdaptor) Login(c *gin.Context) {
	var req dto.AuthLoginRequest

	if err := c.ShouldBind(&req); err != nil {
		utils.BadRequestResponse(c, "Format data tidak valid", err.Error())
		return
	}

	if strings.TrimSpace(req.Username) != "" {
		c.Set(middleware.RequestLogTargetUserKey, req.Username)
	}

	data, err := h.authUsecase.Login(c.Request.Context(), req)
	if err != nil {
		switch err.Error() {
		case "username wajib diisi", "email belum terdaftar":
			utils.BadRequestResponse(c, err.Error(), nil)
		case "username atau password tidak valid":
			utils.UnauthorizedResponse(c, err.Error())
		case "gagal mengirim otp":
			utils.InternalServerErrorResponse(c, err.Error(), nil)
		default:
			utils.InternalServerErrorResponse(c, err.Error(), nil)
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "OTP dikirim", data)
}

func (h *AuthAdaptor) LoginWithEmail(c *gin.Context) {
	var req dto.AuthLoginEmailRequest

	if err := c.ShouldBind(&req); err != nil {
		utils.BadRequestResponse(c, "Format data tidak valid", err.Error())
		return
	}

	if strings.TrimSpace(req.Email) != "" {
		c.Set(middleware.RequestLogTargetUserKey, req.Email)
	}

	data, err := h.authUsecase.LoginWithEmail(c.Request.Context(), req)
	if err != nil {
		switch err.Error() {
		case "email atau password tidak valid":
			utils.UnauthorizedResponse(c, err.Error())
		case "email belum terdaftar":
			utils.BadRequestResponse(c, err.Error(), nil)
		case "gagal mengirim otp":
			utils.InternalServerErrorResponse(c, err.Error(), nil)
		default:
			utils.InternalServerErrorResponse(c, err.Error(), nil)
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "OTP dikirim", data)
}

func (h *AuthAdaptor) VerifyOTP(c *gin.Context) {
	var req dto.AuthVerifyOTPRequest

	if err := c.ShouldBind(&req); err != nil {
		utils.BadRequestResponse(c, "Format data tidak valid", err.Error())
		return
	}

	if strings.TrimSpace(req.UserID) != "" {
		c.Set(middleware.RequestLogTargetUserKey, req.UserID)
	}

	// Extract device info from request
	userAgent := c.Request.Header.Get("User-Agent")
	ip := c.Request.Header.Get("X-Test-IP") // Testing header
	if ip == "" {
		ip = c.Request.Header.Get("X-Forwarded-For") // Production reverse proxy
	}
	if ip == "" {
		ip = c.ClientIP() // Default
	}
	deviceFingerprint := utils.GenerateDeviceFingerprint(userAgent, ip)

	data, err := h.authUsecase.VerifyOTP(c.Request.Context(), req, userAgent, ip, deviceFingerprint)
	if err != nil {
		switch err.Error() {
		case "user id wajib diisi", "otp wajib diisi", "invalid user ID":
			utils.BadRequestResponse(c, err.Error(), nil)
		case "otp tidak valid atau kadaluarsa":
			utils.UnauthorizedResponse(c, err.Error())
		case "user not found":
			utils.NotFoundResponse(c, err.Error())
		default:
			utils.InternalServerErrorResponse(c, err.Error(), nil)
		}
		return
	}

	h.setAuthCookie(c, data.Token)
	utils.SuccessResponse(c, http.StatusOK, "Login berhasil", data)
}

func (h *AuthAdaptor) Logout(c *gin.Context) {
	h.clearAuthCookie(c)
	utils.SuccessResponse(c, http.StatusOK, "Logout berhasil", nil)
}

func (h *AuthAdaptor) GetProfile(c *gin.Context) {
	userID, exists := c.Get(middleware.AuthUserIDKey)
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		utils.UnauthorizedResponse(c, "Invalid user identity")
		return
	}

	data, err := h.authUsecase.GetProfile(c.Request.Context(), userIDStr)
	if err != nil {
		switch err.Error() {
		case "user not found":
			utils.NotFoundResponse(c, err.Error())
		default:
			utils.InternalServerErrorResponse(c, err.Error(), nil)
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Profile retrieved successfully", data)
}

func (h *AuthAdaptor) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get(middleware.AuthUserIDKey)
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		utils.UnauthorizedResponse(c, "Invalid user identity")
		return
	}

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Format data tidak valid", err.Error())
		return
	}

	data, err := h.authUsecase.UpdateProfile(c.Request.Context(), userIDStr, req)
	if err != nil {
		switch err.Error() {
		case "user not found":
			utils.NotFoundResponse(c, err.Error())
		case "format tanggal lahir tidak valid, gunakan YYYY-MM-DD":
			utils.BadRequestResponse(c, err.Error(), nil)
		default:
			utils.InternalServerErrorResponse(c, err.Error(), nil)
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Profile updated successfully", data)
}

func (h *AuthAdaptor) setAuthCookie(c *gin.Context, token string) {
	if strings.TrimSpace(token) == "" {
		return
	}

	secure := h.cfg.App.Env == "production"
	maxAge := int(h.cfg.JWT.AccessTokenExpire.Seconds())
	if maxAge < 0 {
		maxAge = 0
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(middleware.AuthTokenCookieName, token, maxAge, "/", "", secure, true)
}

func (h *AuthAdaptor) clearAuthCookie(c *gin.Context) {
	secure := h.cfg.App.Env == "production"
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(middleware.AuthTokenCookieName, "", -1, "/", "", secure, true)
}
