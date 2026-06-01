package usecase

import (
	"jantungin-api-server/internal/data/repository"
	"jantungin-api-server/internal/services"
	"jantungin-api-server/pkg/utils"

	otp "github.com/Myriadn/otp-lcm"

	"gorm.io/gorm"
)

type UseCase struct {
	AuthUseCase      AuthUsecase
	DiagnosisUseCase DiagnosisUsecase
	StatsUseCase     StatsUsecase
	PatientUseCase   PatientUsecase
}

func NewUseCase(userRepo repository.UserRepository, diagnosisRepo repository.DiagnosisRepository, statsRepo repository.StatsRepository, userDeviceRepo repository.UserDeviceRepository, cfg *utils.Config, db *gorm.DB) *UseCase {
	mlClient := services.NewMLClient(cfg.App.MLServiceURL)
	emailService := services.NewEmailService(cfg)
	otpStore := otp.NewMemoryStore()

	return &UseCase{
		AuthUseCase:      NewAuthUsecase(userRepo, userDeviceRepo, cfg, otpStore, emailService),
		DiagnosisUseCase: NewDiagnosisUsecase(diagnosisRepo, userRepo, mlClient),
		StatsUseCase:     NewStatsUsecase(statsRepo),
		PatientUseCase:   NewPatientUsecase(userRepo),
	}
}
