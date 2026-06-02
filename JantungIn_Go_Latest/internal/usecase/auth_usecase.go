package usecase

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"time"

	"jantungin-api-server/internal/data/entity"
	"jantungin-api-server/internal/data/repository"
	"jantungin-api-server/internal/dto"
	"jantungin-api-server/internal/services"
	"jantungin-api-server/pkg/utils"

	otp "github.com/Myriadn/otp-lcm"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase interface {
	Register(ctx context.Context, req dto.AuthRegisterRequest) (*dto.AuthRegisterData, error)
	Login(ctx context.Context, req dto.AuthLoginRequest) (*dto.AuthLoginOTPData, error)
	LoginWithEmail(ctx context.Context, req dto.AuthLoginEmailRequest) (*dto.AuthLoginOTPData, error)
	VerifyOTP(ctx context.Context, req dto.AuthVerifyOTPRequest, userAgent, ipAddress, deviceFingerprint string) (*dto.AuthLoginData, error)
	GetProfile(ctx context.Context, userID string) (*dto.AuthUserResponse, error)
	UpdateProfile(ctx context.Context, userID string, req dto.UpdateProfileRequest) (*dto.UpdateProfileData, error)
}

type authUsecase struct {
	userRepo       repository.UserRepository
	userDeviceRepo repository.UserDeviceRepository
	cfg            *utils.Config
	otpStore       *otp.MemoryStore
	emailService   services.EmailService
}

func NewAuthUsecase(userRepo repository.UserRepository, userDeviceRepo repository.UserDeviceRepository, cfg *utils.Config, otpStore *otp.MemoryStore, emailService services.EmailService) AuthUsecase {
	return &authUsecase{
		userRepo:       userRepo,
		userDeviceRepo: userDeviceRepo,
		cfg:            cfg,
		otpStore:       otpStore,
		emailService:   emailService,
	}
}

const otpTTL = 30 * time.Second

var usernamePattern = regexp.MustCompile(`^[A-Za-z0-9]+([._][A-Za-z0-9]+)*$`)

func (u *authUsecase) Register(ctx context.Context, req dto.AuthRegisterRequest) (*dto.AuthRegisterData, error) {
	rawUsername := req.Username
	trimmedUsername := strings.TrimSpace(rawUsername)
	if trimmedUsername == "" {
		return nil, errors.New("username wajib diisi")
	}
	if trimmedUsername != rawUsername {
		return nil, errors.New("format username tidak valid")
	}
	if len(trimmedUsername) < 3 {
		return nil, errors.New("username minimal 3 karakter")
	}
	if len(trimmedUsername) > 16 {
		return nil, errors.New("username maksimal 16 karakter")
	}
	if !usernamePattern.MatchString(trimmedUsername) {
		return nil, errors.New("format username tidak valid")
	}

	username := strings.ToLower(trimmedUsername)

	if len(req.Password) < 6 {
		return nil, errors.New("password minimal 6 karakter")
	}

	existingUser, err := u.userRepo.FindByUsername(ctx, username)
	if err != nil {
		utils.Error("Failed to check username", zap.Error(err))
		return nil, errors.New("gagal memeriksa username")
	}
	if existingUser != nil {
		return nil, errors.New("username sudah terdaftar")
	}

	normalizedEmail := strings.ToLower(strings.TrimSpace(req.Email))

	// Cek email jika diisi
	if normalizedEmail != "" {
		existingUserByEmail, err := u.userRepo.FindByEmail(ctx, normalizedEmail)
		if err != nil {
			utils.Error("Failed to check email", zap.Error(err))
			return nil, errors.New("gagal memeriksa email")
		}
		if existingUserByEmail != nil {
			return nil, errors.New("email sudah terdaftar")
		}
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.Error("Failed to hash password", zap.Error(err))
		return nil, errors.New("gagal memproses password")
	}

	var dob *time.Time
	if req.DateOfBirth != "" {
		parsedDate, parseErr := time.Parse("2006-01-02", req.DateOfBirth)
		if parseErr != nil {
			return nil, errors.New("format tanggal lahir tidak valid, gunakan YYYY-MM-DD")
		}
		dob = &parsedDate
	}

	displayName := strings.TrimSpace(req.Name)
	if displayName == "" {
		displayName = username
	}

	newUser := entity.User{
		Name:        displayName,
		Username:    &username,
		Password:    string(hashedPassword),
		Role:        "user",
		DateOfBirth: dob,
	}
	if normalizedEmail != "" {
		newUser.Email = &normalizedEmail
	}

	if err := u.userRepo.Create(ctx, &newUser); err != nil {
		utils.Error("Failed to create user", zap.Error(err))
		return nil, errors.New("gagal membuat akun")
	}

	utils.Info("User registered successfully",
		zap.String("user_id", newUser.ID.String()),
		zap.String("role", newUser.Role),
	)

	return &dto.AuthRegisterData{
		ID:       newUser.ID.String(),
		Name:     newUser.Name,
		Username: newUser.Username,
		Email:    newUser.Email,
		Role:     newUser.Role,
	}, nil
}

func (u *authUsecase) Login(ctx context.Context, req dto.AuthLoginRequest) (*dto.AuthLoginOTPData, error) {
	username := strings.ToLower(strings.TrimSpace(req.Username))
	if username == "" {
		return nil, errors.New("username wajib diisi")
	}

	foundUser, err := u.userRepo.FindByUsername(ctx, username)
	if err != nil {
		utils.Error("Failed to find user by username", zap.Error(err))
		return nil, errors.New("username atau password tidak valid")
	}
	if foundUser == nil {
		utils.Warn("No user found with matching username")
		return nil, errors.New("username atau password tidak valid")
	}

	// Verifikasi password
	if err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(req.Password)); err != nil {
		utils.Warn("Password verification failed",
			zap.String("user_id", foundUser.ID.String()),
		)
		return nil, errors.New("username atau password tidak valid")
	}

	if foundUser.Email == nil || strings.TrimSpace(*foundUser.Email) == "" {
		return nil, errors.New("email belum terdaftar")
	}
	if foundUser.Username == nil || strings.TrimSpace(*foundUser.Username) == "" {
		return nil, errors.New("username tidak ditemukan")
	}

	normalizedUsername := strings.ToLower(strings.TrimSpace(*foundUser.Username))
	code := otp.GenerateOTP(normalizedUsername, time.Now())

	if u.cfg.App.Env == "development" {
		utils.Debug("OTP generated (dev)",
			zap.String("user_id", foundUser.ID.String()),
			zap.String("otp_code", code),
			zap.Int("expires_in", int(otpTTL.Seconds())),
		)
	}

	if err := u.emailService.SendOTP(ctx, *foundUser.Email, foundUser.Name, code, otpTTL); err != nil {
		utils.Error("Failed to send OTP", zap.Error(err))
		return nil, errors.New("gagal mengirim otp")
	}

	u.otpStore.Save(foundUser.ID.String(), code, otpTTL)

	utils.Info("OTP sent",
		zap.String("user_id", foundUser.ID.String()),
		zap.String("role", foundUser.Role),
	)

	return &dto.AuthLoginOTPData{
		ID:           foundUser.ID.String(),
		Name:         foundUser.Name,
		Username:     foundUser.Username,
		Email:        foundUser.Email,
		Role:         foundUser.Role,
		OTPRequired:  true,
		OTPExpiresIn: int(otpTTL.Seconds()),
	}, nil
}

func (u *authUsecase) LoginWithEmail(ctx context.Context, req dto.AuthLoginEmailRequest) (*dto.AuthLoginOTPData, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))
	user, err := u.userRepo.FindByEmail(ctx, email)
	if err != nil {
		utils.Error("Failed to find user by email", zap.Error(err))
		return nil, errors.New("email atau password tidak valid")
	}
	if user == nil {
		utils.Warn("No user found with email", zap.String("email", email))
		return nil, errors.New("email atau password tidak valid")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		utils.Warn("Password verification failed for email login",
			zap.String("user_id", user.ID.String()),
		)
		return nil, errors.New("email atau password tidak valid")
	}

	if user.Email == nil || strings.TrimSpace(*user.Email) == "" {
		return nil, errors.New("email belum terdaftar")
	}
	if user.Username == nil || strings.TrimSpace(*user.Username) == "" {
		return nil, errors.New("username tidak ditemukan")
	}

	normalizedUsername := strings.ToLower(strings.TrimSpace(*user.Username))
	code := otp.GenerateOTP(normalizedUsername, time.Now())

	if u.cfg.App.Env == "development" {
		utils.Debug("OTP generated (dev)",
			zap.String("user_id", user.ID.String()),
			zap.String("otp_code", code),
			zap.Int("expires_in", int(otpTTL.Seconds())),
		)
	}

	if err := u.emailService.SendOTP(ctx, *user.Email, user.Name, code, otpTTL); err != nil {
		utils.Error("Failed to send OTP", zap.Error(err))
		return nil, errors.New("gagal mengirim otp")
	}

	u.otpStore.Save(user.ID.String(), code, otpTTL)

	utils.Info("OTP sent for email login",
		zap.String("user_id", user.ID.String()),
		zap.String("role", user.Role),
	)

	return &dto.AuthLoginOTPData{
		ID:           user.ID.String(),
		Name:         user.Name,
		Username:     user.Username,
		Email:        user.Email,
		Role:         user.Role,
		OTPRequired:  true,
		OTPExpiresIn: int(otpTTL.Seconds()),
	}, nil
}

func (u *authUsecase) VerifyOTP(ctx context.Context, req dto.AuthVerifyOTPRequest, userAgent, ipAddress, deviceFingerprint string) (*dto.AuthLoginData, error) {
	if strings.TrimSpace(req.UserID) == "" {
		return nil, errors.New("user id wajib diisi")
	}
	if strings.TrimSpace(req.Code) == "" {
		return nil, errors.New("otp wajib diisi")
	}

	uid, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	user, err := u.userRepo.FindByID(ctx, uid)
	if err != nil {
		utils.Error("Failed to find user by ID", zap.Error(err))
		return nil, errors.New("gagal mengambil user")
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	if !u.otpStore.Validate(user.ID.String(), strings.TrimSpace(req.Code)) {
		return nil, errors.New("otp tidak valid atau kadaluarsa")
	}

	token, err := u.generateToken(user)
	if err != nil {
		utils.Error("Failed to generate token", zap.Error(err))
		return nil, errors.New("gagal membuat token")
	}

	if err := u.userDeviceRepo.CreateOrUpdate(ctx, user.ID, userAgent, ipAddress, deviceFingerprint); err != nil {
		utils.Warn("Failed to track device login",
			zap.String("user_id", user.ID.String()),
			zap.Error(err),
		)
	}

	utils.Info("User logged in after OTP verification",
		zap.String("user_id", user.ID.String()),
		zap.String("role", user.Role),
	)

	return &dto.AuthLoginData{
		ID:       user.ID.String(),
		Name:     user.Name,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
		Token:    token,
	}, nil
}

func (u *authUsecase) GetProfile(ctx context.Context, userID string) (*dto.AuthUserResponse, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	user, err := u.userRepo.FindByID(ctx, uid)
	if err != nil {
		utils.Error("Failed to find user by ID", zap.Error(err))
		return nil, errors.New("gagal mengambil profil")
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	return &dto.AuthUserResponse{
		ID:       user.ID.String(),
		Name:     user.Name,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	}, nil
}

func (u *authUsecase) UpdateProfile(ctx context.Context, userID string, req dto.UpdateProfileRequest) (*dto.UpdateProfileData, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	user, err := u.userRepo.FindByID(ctx, uid)
	if err != nil {
		utils.Error("Failed to find user by ID", zap.Error(err))
		return nil, errors.New("gagal mengambil profil")
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.DateOfBirth != "" {
		parsedDate, parseErr := time.Parse("2006-01-02", req.DateOfBirth)
		if parseErr != nil {
			return nil, errors.New("format tanggal lahir tidak valid, gunakan YYYY-MM-DD")
		}
		user.DateOfBirth = &parsedDate
	}

	if err := u.userRepo.Update(ctx, user); err != nil {
		utils.Error("Failed to update user", zap.Error(err))
		return nil, errors.New("gagal memperbarui profil")
	}

	utils.Info("User profile updated",
		zap.String("user_id", user.ID.String()),
	)

	result := &dto.UpdateProfileData{
		ID:       user.ID.String(),
		Name:     user.Name,
		Username: user.Username,
		Email:    user.Email,
	}
	if user.DateOfBirth != nil {
		formatted := user.DateOfBirth.Format("2006-01-02")
		result.DateOfBirth = &formatted
	}

	return result, nil
}

// generateToken membuat JWT access token dari data user
func (u *authUsecase) generateToken(user *entity.User) (string, error) {
	email := ""
	if user.Email != nil {
		email = *user.Email
	}

	return utils.GenerateAccessToken(
		user.ID.String(),
		email,
		user.Role, // RoleID diisi dengan role string
		user.Role, // RoleCode diisi dengan role string
		u.cfg,
	)
}
