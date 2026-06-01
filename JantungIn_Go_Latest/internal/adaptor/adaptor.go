package adaptor

import (
	"jantungin-api-server/internal/usecase"
	"jantungin-api-server/pkg/utils"
)

type Adaptor struct {
	AuthAdaptor      *AuthAdaptor
	DiagnosisAdaptor *DiagnosisAdaptor
	StatsAdaptor     *StatsAdaptor
	PatientAdaptor   *PatientAdaptor
}

func NewAdaptor(usecases *usecase.UseCase, cfg *utils.Config) *Adaptor {
	return &Adaptor{
		AuthAdaptor:      NewAuthAdaptor(usecases.AuthUseCase, cfg),
		DiagnosisAdaptor: NewDiagnosisAdaptor(usecases.DiagnosisUseCase),
		StatsAdaptor:     NewStatsAdaptor(usecases.StatsUseCase),
		PatientAdaptor:   NewPatientAdaptor(usecases.PatientUseCase),
	}
}
