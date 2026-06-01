package wire

import (
	"jantungin-api-server/internal/adaptor"
	"jantungin-api-server/internal/data/repository"
	"jantungin-api-server/internal/usecase"
	"jantungin-api-server/pkg/middleware"
	"jantungin-api-server/pkg/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Wiring(cfg *utils.Config, db *gorm.DB) *gin.Engine {
	utils.Info("Starting route wiring process")

	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
		utils.Info("Gin mode set to release")
	} else {
		utils.Info("Gin mode set to debug")
	}

	router := gin.New()
	router.Use(middleware.Recovery())
	router.Use(middleware.Logger())
	router.Use(middleware.CORS(cfg))

	// Initialize repositories
	repo := repository.NewRepository(db)

	// Initialize usecases
	usecases := usecase.NewUseCase(repo.UserRepo, repo.DiagnosisRepo, repo.StatsRepo, repo.UserDeviceRepo, cfg, db)

	// Initialize adaptors
	adaptors := adaptor.NewAdaptor(usecases, cfg)

	// RequestTracker middleware — catat setiap request ke DB
	// Dipasang setelah router global middleware agar status code sudah tersedia
	router.Use(middleware.RequestTracker(repo.StatsRepo))

	// Register routes
	api := router.Group("/api/v1")
	registerAuthRoutes(api, adaptors, cfg)
	registerDiagnosisRoutes(api, adaptors, cfg)
	registerStatsRoutes(api, adaptors, cfg)
	registerPatientRoutes(api, adaptors, cfg)

	utils.Info("Route wiring completed")

	return router
}

func registerAuthRoutes(api *gin.RouterGroup, adaptors *adaptor.Adaptor, cfg *utils.Config) {
	// Auth routes (public)
	auth := api.Group("/auth")
	{
		auth.POST("/register", adaptors.AuthAdaptor.Register)
		auth.POST("/login", adaptors.AuthAdaptor.Login)
		auth.POST("/login-email", adaptors.AuthAdaptor.LoginWithEmail)
		auth.POST("/verify-otp", adaptors.AuthAdaptor.VerifyOTP)
	}

	// Auth routes (protected)
	authProtected := api.Group("/auth")
	authProtected.Use(middleware.AuthRequired(cfg))
	{
		authProtected.GET("/profile", adaptors.AuthAdaptor.GetProfile)
		authProtected.PUT("/profile", adaptors.AuthAdaptor.UpdateProfile)
		authProtected.POST("/logout", adaptors.AuthAdaptor.Logout)
	}
}

func registerDiagnosisRoutes(api *gin.RouterGroup, adaptors *adaptor.Adaptor, cfg *utils.Config) {
	// Semua user terauth: ambil history dan detail
	diagnosis := api.Group("/diagnosis")
	diagnosis.Use(middleware.AuthRequired(cfg))
	{
		diagnosis.GET("/history", adaptors.DiagnosisAdaptor.GetDiagnosisHistory)
		diagnosis.GET("/:id", adaptors.DiagnosisAdaptor.GetDiagnosisByID)
	}

	// Hanya admin/dokter: buat diagnosis
	diagnosisAdmin := api.Group("/diagnosis")
	diagnosisAdmin.Use(middleware.AuthRequired(cfg))
	diagnosisAdmin.Use(middleware.RoleRequired("admin", "dokter"))
	{
		diagnosisAdmin.POST("", adaptors.DiagnosisAdaptor.CreateDiagnosis)
	}

	// Hanya admin/dokter: endpoint admin
	admin := api.Group("/admin/diagnosis")
	admin.Use(middleware.AuthRequired(cfg))
	admin.Use(middleware.RoleRequired("admin", "dokter"))
	{
		admin.GET("/all", adaptors.DiagnosisAdaptor.GetAllDiagnoses)
		admin.GET("/patient/:patientId", adaptors.DiagnosisAdaptor.GetPatientDiagnoses)
	}
}

func registerPatientRoutes(api *gin.RouterGroup, adaptors *adaptor.Adaptor, cfg *utils.Config) {
	// Semua endpoint patient hanya untuk admin/dokter
	patients := api.Group("/admin/patients")
	patients.Use(middleware.AuthRequired(cfg))
	patients.Use(middleware.RoleRequired("admin", "dokter"))
	{
		// GET /api/v1/admin/patients/search?query=... — harus sebelum /:id
		patients.GET("/search", adaptors.PatientAdaptor.SearchPatients)
		// GET /api/v1/admin/patients
		patients.GET("", adaptors.PatientAdaptor.GetAllPatients)
		// GET /api/v1/admin/patients/:id
		patients.GET("/:id", adaptors.PatientAdaptor.GetPatientByID)
	}
}

func registerStatsRoutes(api *gin.RouterGroup, adaptors *adaptor.Adaptor, cfg *utils.Config) {
	// Public endpoint — untuk homepage (total users & diagnoses)
	api.GET("/stats", adaptors.StatsAdaptor.GetPublicStats)

	// Admin only endpoint — untuk dashboard admin (kunjungan, grafik harian, dll)
	adminStats := api.Group("/admin/stats")
	adminStats.Use(middleware.AuthRequired(cfg))
	adminStats.Use(middleware.RoleRequired("admin", "dokter"))
	{
		adminStats.GET("", adaptors.StatsAdaptor.GetAdminStats)
	}
}
